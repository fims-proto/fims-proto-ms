package db

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	"github/fims-proto/fims-proto-ms/internal/common/data/converter"
	"github/fims-proto/fims-proto-ms/internal/report/app/query"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/shopspring/decimal"
)

type reportPO struct {
	Id       uuid.UUID  `gorm:"type:uuid;primaryKey"`
	SobId    uuid.UUID  `gorm:"type:uuid;uniqueIndex:UQ_Reports_SobId_Class_PeriodId,where:template = false"`
	PeriodId *uuid.UUID `gorm:"type:uuid;uniqueIndex:UQ_Reports_SobId_Class_PeriodId,where:template = false"`
	Title    string
	Template bool
	Class    string `gorm:"uniqueIndex:UQ_Reports_SobId_Class_PeriodId,where:template = false"`

	Columns []*reportColumnPO `gorm:"foreignKey:ReportId"`
	Rows    []*reportRowPO    `gorm:"foreignKey:ReportId"`
	Period  *periodPO         `gorm:"foreignKey:PeriodId"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

type reportColumnPO struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey"`
	ReportId  uuid.UUID `gorm:"type:uuid"`
	Label     string
	ValueType string
	Sequence  int

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

type reportRowPO struct {
	Id          uuid.UUID  `gorm:"type:uuid;primaryKey"`
	ReportId    uuid.UUID  `gorm:"type:uuid"`
	ParentRowId *uuid.UUID `gorm:"type:uuid"`
	RowCode     string
	Text        string
	Sequence    int
	LineNo      *int
	ShowLineNo  bool
	SumFactor   int
	CanEdit     bool
	CanMove     bool
	CanAddChild bool
	Amounts     pgtype.TextArray    `gorm:"type:text[]"`
	Expression  *reportExpressionPO `gorm:"foreignKey:RowId"`
	Rows        []*reportRowPO      `gorm:"foreignKey:ParentRowId"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

type reportExpressionPO struct {
	Id                 uuid.UUID `gorm:"type:uuid;primaryKey"`
	RowId              uuid.UUID `gorm:"type:uuid;uniqueIndex"`
	Kind               string
	LedgerAccountsJSON string `gorm:"type:jsonb"`
	CashFlowItemsJSON  string `gorm:"type:jsonb"`
	RowReferencesJSON  string `gorm:"type:jsonb"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

type periodPO struct {
	Id           uuid.UUID `gorm:"type:uuid;primaryKey"`
	SobId        uuid.UUID `gorm:"type:uuid"`
	FiscalYear   int
	PeriodNumber int
}

func (r reportPO) ResolveAssociation(entity string) (string, error) {
	if entity == "" {
		return "reports", nil
	}
	if strings.EqualFold(entity, "columns") {
		return "Columns", nil
	}
	if strings.EqualFold(entity, "rows") {
		return "Rows", nil
	}
	if strings.EqualFold(entity, "period") {
		return "Period", nil
	}
	return "", fmt.Errorf("reportPO does not have association named %s", entity)
}

func reportBOToPO(bo *report.Report) *reportPO {
	var columns []*reportColumnPO
	for _, column := range bo.Columns() {
		columns = append(columns, columnBOToPO(column, bo.Id()))
	}

	var rows []*reportRowPO
	for _, row := range bo.Rows() {
		rows = append(rows, rowBOToPO(row, bo.Id(), uuid.Nil))
	}

	return &reportPO{
		Id:       bo.Id(),
		SobId:    bo.SobId(),
		PeriodId: converter.UUIDToPtr(bo.PeriodId()),
		Title:    bo.Title(),
		Template: bo.Template(),
		Class:    bo.Class(),
		Columns:  columns,
		Rows:     rows,
	}
}

func columnBOToPO(bo *report.Column, reportId uuid.UUID) *reportColumnPO {
	return &reportColumnPO{
		Id:        bo.Id(),
		ReportId:  reportId,
		Label:     bo.Label(),
		ValueType: bo.ValueType(),
		Sequence:  bo.Sequence(),
	}
}

func rowBOToPO(bo *report.Row, reportId uuid.UUID, parentRowId uuid.UUID) *reportRowPO {
	var rows []*reportRowPO
	for _, row := range bo.Rows() {
		rows = append(rows, rowBOToPO(row, reportId, bo.Id()))
	}
	amounts, err := decimalArrayToTextArray(bo.Amounts())
	if err != nil {
		panic(fmt.Errorf("failed to convert row amounts: %w", err))
	}
	return &reportRowPO{
		Id:          bo.Id(),
		ReportId:    reportId,
		ParentRowId: converter.UUIDToPtr(parentRowId),
		RowCode:     bo.RowCode(),
		Text:        bo.Text(),
		Sequence:    bo.Sequence(),
		LineNo:      bo.LineNo(),
		ShowLineNo:  bo.ShowLineNo(),
		SumFactor:   bo.SumFactor(),
		CanEdit:     bo.CanEdit(),
		CanMove:     bo.CanMove(),
		CanAddChild: bo.CanAddChild(),
		Amounts:     amounts,
		Expression:  expressionBOToPO(bo.Expression(), bo.Id()),
		Rows:        rows,
	}
}

func expressionBOToPO(bo *report.Expression, rowId uuid.UUID) *reportExpressionPO {
	return &reportExpressionPO{
		Id:                 bo.Id(),
		RowId:              rowId,
		Kind:               bo.Kind(),
		LedgerAccountsJSON: mustMarshalJSON(bo.LedgerAccounts()),
		CashFlowItemsJSON:  mustMarshalJSON(bo.CashFlowItems()),
		RowReferencesJSON:  mustMarshalJSON(bo.RowReferences()),
	}
}

func reportPOToBO(po *reportPO) (*report.Report, error) {
	restructureAndSort(po)

	columns, err := converter.POsToBOs(po.Columns, columnPOToBO)
	if err != nil {
		return nil, err
	}
	rows, err := converter.POsToBOs(po.Rows, rowPOToBO)
	if err != nil {
		return nil, err
	}
	return report.New(po.Id, po.SobId, converter.UUIDFromPtr(po.PeriodId), po.Title, po.Template, po.Class, columns, rows)
}

func columnPOToBO(po *reportColumnPO) (*report.Column, error) {
	return report.NewColumn(po.Id, po.Label, po.ValueType, po.Sequence)
}

func rowPOToBO(po *reportRowPO) (*report.Row, error) {
	expr, err := expressionPOToBO(po.Expression)
	if err != nil {
		return nil, err
	}
	rows, err := converter.POsToBOs(po.Rows, rowPOToBO)
	if err != nil {
		return nil, err
	}
	amounts, err := textArrayToDecimalArray(po.Amounts)
	if err != nil {
		return nil, err
	}
	return report.NewRow(po.Id, po.RowCode, po.Text, po.Sequence, po.LineNo, po.ShowLineNo, po.SumFactor, po.CanEdit, po.CanMove, po.CanAddChild, expr, rows, amounts)
}

func expressionPOToBO(po *reportExpressionPO) (*report.Expression, error) {
	if po == nil {
		return report.NewExpression(uuid.New(), report.ExpressionNone, nil, nil, nil)
	}
	var ledgerAccounts []report.LedgerAccountReference
	var cashFlowItems []report.CashFlowItemReference
	var rowReferences []report.RowReference
	if err := json.Unmarshal([]byte(defaultJSON(po.LedgerAccountsJSON)), &ledgerAccounts); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(defaultJSON(po.CashFlowItemsJSON)), &cashFlowItems); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(defaultJSON(po.RowReferencesJSON)), &rowReferences); err != nil {
		return nil, err
	}
	return report.NewExpression(po.Id, po.Kind, ledgerAccounts, cashFlowItems, rowReferences)
}

func reportPOToDTO(po reportPO) query.Report {
	restructureAndSort(&po)
	return query.Report{
		Id:        po.Id,
		SobId:     po.SobId,
		Period:    periodPOToDTO(po.Period),
		Title:     po.Title,
		Template:  po.Template,
		Class:     po.Class,
		Columns:   converter.POsToDTOs(po.Columns, columnPOToDTO),
		Rows:      converter.POsToDTOs(po.Rows, rowPOToDTO),
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,
	}
}

func columnPOToDTO(po *reportColumnPO) query.Column {
	return query.Column{Id: po.Id, Label: po.Label, ValueType: po.ValueType, Sequence: po.Sequence}
}

func rowPOToDTO(po *reportRowPO) query.Row {
	amounts, err := textArrayToDecimalArray(po.Amounts)
	if err != nil {
		panic(fmt.Errorf("failed to convert row amounts: %w", err))
	}
	return query.Row{
		Id:          po.Id,
		RowCode:     po.RowCode,
		Text:        po.Text,
		Sequence:    po.Sequence,
		LineNo:      po.LineNo,
		ShowLineNo:  po.ShowLineNo,
		SumFactor:   po.SumFactor,
		CanEdit:     po.CanEdit,
		CanMove:     po.CanMove,
		CanAddChild: po.CanAddChild,
		Expression:  expressionPOToDTO(po.Expression),
		Rows:        converter.POsToDTOs(po.Rows, rowPOToDTO),
		Amounts:     amounts,
	}
}

func expressionPOToDTO(po *reportExpressionPO) query.Expression {
	if po == nil {
		return query.Expression{Kind: report.ExpressionNone}
	}
	var ledgerAccounts []query.LedgerAccountReference
	var cashFlowItems []query.CashFlowItemReference
	var rowReferences []query.RowReference
	if err := json.Unmarshal([]byte(defaultJSON(po.LedgerAccountsJSON)), &ledgerAccounts); err != nil {
		panic(err)
	}
	if err := json.Unmarshal([]byte(defaultJSON(po.CashFlowItemsJSON)), &cashFlowItems); err != nil {
		panic(err)
	}
	if err := json.Unmarshal([]byte(defaultJSON(po.RowReferencesJSON)), &rowReferences); err != nil {
		panic(err)
	}
	return query.Expression{
		Id:             po.Id,
		Kind:           po.Kind,
		LedgerAccounts: ledgerAccounts,
		CashFlowItems:  cashFlowItems,
		RowReferences:  rowReferences,
	}
}

func periodPOToDTO(po *periodPO) *query.Period {
	if po == nil {
		return nil
	}
	return &query.Period{FiscalYear: po.FiscalYear, PeriodNumber: po.PeriodNumber}
}

func restructureAndSort(r *reportPO) {
	rowMap := make(map[uuid.UUID]*reportRowPO)
	for _, row := range r.Rows {
		row.Rows = nil
		rowMap[row.Id] = row
	}
	for _, row := range r.Rows {
		if row.ParentRowId != nil {
			if parent, ok := rowMap[*row.ParentRowId]; ok {
				parent.Rows = append(parent.Rows, row)
			}
		}
	}
	r.Rows = nil
	for row := range maps.Values(rowMap) {
		if row.ParentRowId == nil {
			r.Rows = append(r.Rows, row)
		}
	}
	for _, row := range r.Rows {
		sortRowsRecursive(row)
	}
	slices.SortFunc(r.Columns, func(a, b *reportColumnPO) int { return a.Sequence - b.Sequence })
	slices.SortFunc(r.Rows, func(a, b *reportRowPO) int { return a.Sequence - b.Sequence })
}

func sortRowsRecursive(row *reportRowPO) {
	slices.SortFunc(row.Rows, func(a, b *reportRowPO) int { return a.Sequence - b.Sequence })
	for _, child := range row.Rows {
		sortRowsRecursive(child)
	}
}

func decimalArrayToTextArray(decimalArray []decimal.Decimal) (pgtype.TextArray, error) {
	var tempStrs []string
	for _, d := range decimalArray {
		tempStrs = append(tempStrs, d.String())
	}
	var textArray pgtype.TextArray
	if err := textArray.Set(tempStrs); err != nil {
		return pgtype.TextArray{}, err
	}
	return textArray, nil
}

func textArrayToDecimalArray(textArray pgtype.TextArray) ([]decimal.Decimal, error) {
	if textArray.Status != pgtype.Present {
		return nil, nil
	}
	var tempStrs []string
	if err := textArray.AssignTo(&tempStrs); err != nil {
		return nil, err
	}
	var decimalArray []decimal.Decimal
	for _, str := range tempStrs {
		d, err := decimal.NewFromString(str)
		if err != nil {
			return nil, err
		}
		decimalArray = append(decimalArray, d)
	}
	return decimalArray, nil
}

func mustMarshalJSON(v any) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func defaultJSON(value string) string {
	if value == "" {
		return "[]"
	}
	return value
}
