package report

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const (
	ClassBalanceSheet      = "balance_sheet"
	ClassIncomeStatement   = "income_statement"
	ClassCashFlowStatement = "cash_flow_statement"

	ColumnPeriodEndingBalance = "period_ending_balance"
	ColumnYearOpeningBalance  = "year_opening_balance"
	ColumnPeriodAmount        = "period_amount"
	ColumnYearToDateAmount    = "year_to_date_amount"
	ColumnLastYearAmount      = "last_year_amount"

	ExpressionLedgerAccounts = "ledger_accounts"
	ExpressionCashFlowItems  = "cash_flow_items"
	ExpressionRowsExplicit   = "rows_explicit"
	ExpressionChildrenSum    = "children_sum"
	ExpressionNone           = "none"

	LedgerMeasureNet            = "net"
	LedgerMeasureDebit          = "debit"
	LedgerMeasureCredit         = "credit"
	LedgerMeasureTransaction    = "transaction"
	LedgerMeasureOpeningBalance = "opening_balance"
)

var validClasses = []string{ClassBalanceSheet, ClassIncomeStatement, ClassCashFlowStatement}

var validColumnTypes = []string{
	ColumnPeriodEndingBalance,
	ColumnYearOpeningBalance,
	ColumnPeriodAmount,
	ColumnYearToDateAmount,
	ColumnLastYearAmount,
}

var validExpressionKinds = []string{
	ExpressionLedgerAccounts,
	ExpressionCashFlowItems,
	ExpressionRowsExplicit,
	ExpressionChildrenSum,
	ExpressionNone,
}

var validLedgerMeasures = []string{
	LedgerMeasureNet,
	LedgerMeasureDebit,
	LedgerMeasureCredit,
	LedgerMeasureTransaction,
	LedgerMeasureOpeningBalance,
}

type Report struct {
	id       uuid.UUID
	sobId    uuid.UUID
	periodId uuid.UUID
	title    string
	template bool
	class    string
	columns  []*Column
	rows     []*Row
}

type Column struct {
	id        uuid.UUID
	label     string
	valueType string
	sequence  int
}

type Row struct {
	id               uuid.UUID
	rowCode          string
	text             string
	sequence         int
	lineNo           *int
	showLineNo       bool
	sumFactor        int
	displaySumFactor bool
	indent           int
	canEdit          bool
	canMove          bool
	canAddChild      bool
	expression       *Expression
	rows             []*Row
	amounts          []decimal.Decimal
}

type Expression struct {
	id             uuid.UUID
	kind           string
	ledgerAccounts []LedgerAccountReference
	cashFlowItems  []CashFlowItemReference
	rowReferences  []RowReference
}

type LedgerAccountReference struct {
	RawAccountNumber string
	AccountId        uuid.UUID
	SumFactor        int
	Measure          string
}

type CashFlowItemReference struct {
	Code      string
	ItemId    uuid.UUID
	SumFactor int
}

type RowReference struct {
	RowCode   string
	SumFactor int
}

func New(
	id uuid.UUID,
	sobId uuid.UUID,
	periodId uuid.UUID,
	title string,
	template bool,
	class string,
	columns []*Column,
	rows []*Row,
) (*Report, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("report id is required")
	}
	if sobId == uuid.Nil {
		return nil, fmt.Errorf("sob id is required")
	}
	if template && periodId != uuid.Nil {
		return nil, fmt.Errorf("report template cannot have period id")
	}
	if !template && periodId == uuid.Nil {
		return nil, fmt.Errorf("report instance must have period id")
	}
	if title == "" {
		return nil, fmt.Errorf("report title is required")
	}
	if !slices.Contains(validClasses, class) {
		return nil, fmt.Errorf("unknown report class %s", class)
	}
	if len(columns) == 0 {
		return nil, fmt.Errorf("report columns are required")
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("report rows are required")
	}

	return &Report{
		id:       id,
		sobId:    sobId,
		periodId: periodId,
		title:    title,
		template: template,
		class:    class,
		columns:  columns,
		rows:     rows,
	}, nil
}

func NewColumn(id uuid.UUID, label string, valueType string, sequence int) (*Column, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("column id is required")
	}
	if label == "" {
		return nil, fmt.Errorf("column label is required")
	}
	if sequence == 0 {
		return nil, fmt.Errorf("column sequence is required")
	}
	if !slices.Contains(validColumnTypes, valueType) {
		return nil, fmt.Errorf("unknown column value type %s", valueType)
	}
	return &Column{id: id, label: label, valueType: valueType, sequence: sequence}, nil
}

func NewRow(
	id uuid.UUID,
	rowCode string,
	text string,
	sequence int,
	lineNo *int,
	showLineNo bool,
	sumFactor int,
	displaySumFactor bool,
	indent int,
	canEdit bool,
	canMove bool,
	canAddChild bool,
	expression *Expression,
	rows []*Row,
	amounts []decimal.Decimal,
) (*Row, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("row id is required")
	}
	if rowCode == "" {
		return nil, fmt.Errorf("row code is required")
	}
	if text == "" {
		return nil, fmt.Errorf("row text is required")
	}
	if sequence == 0 {
		return nil, fmt.Errorf("row sequence is required")
	}
	if sumFactor != -1 && sumFactor != 0 && sumFactor != 1 {
		return nil, fmt.Errorf("row sum factor must be -1, 0, or 1")
	}
	if expression == nil {
		return nil, fmt.Errorf("row expression is required")
	}
	return &Row{
		id:               id,
		rowCode:          rowCode,
		text:             text,
		sequence:         sequence,
		lineNo:           lineNo,
		showLineNo:       showLineNo,
		sumFactor:        sumFactor,
		displaySumFactor: displaySumFactor,
		indent:           indent,
		canEdit:          canEdit,
		canMove:          canMove,
		canAddChild:      canAddChild,
		expression:       expression,
		rows:             rows,
		amounts:          amounts,
	}, nil
}

func NewExpression(
	id uuid.UUID,
	kind string,
	ledgerAccounts []LedgerAccountReference,
	cashFlowItems []CashFlowItemReference,
	rowReferences []RowReference,
) (*Expression, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("expression id is required")
	}
	if !slices.Contains(validExpressionKinds, kind) {
		return nil, fmt.Errorf("unknown expression kind %s", kind)
	}
	for _, ref := range ledgerAccounts {
		if ref.SumFactor != -1 && ref.SumFactor != 1 {
			return nil, fmt.Errorf("ledger account sum factor must be -1 or 1")
		}
		if !slices.Contains(validLedgerMeasures, ref.Measure) {
			return nil, fmt.Errorf("unknown ledger measure %s", ref.Measure)
		}
	}
	for _, ref := range cashFlowItems {
		if ref.SumFactor != -1 && ref.SumFactor != 1 {
			return nil, fmt.Errorf("cash flow item sum factor must be -1 or 1")
		}
	}
	for _, ref := range rowReferences {
		if ref.RowCode == "" {
			return nil, fmt.Errorf("row reference code is required")
		}
		if ref.SumFactor != -1 && ref.SumFactor != 1 {
			return nil, fmt.Errorf("row reference sum factor must be -1 or 1")
		}
	}
	return &Expression{
		id:             id,
		kind:           kind,
		ledgerAccounts: ledgerAccounts,
		cashFlowItems:  cashFlowItems,
		rowReferences:  rowReferences,
	}, nil
}

func (r *Report) Instantiate(reportId uuid.UUID, periodId uuid.UUID, title string) (*Report, error) {
	if !r.template {
		return nil, fmt.Errorf("only report templates can be instantiated")
	}
	if title == "" {
		title = r.title
	}
	columns := make([]*Column, 0, len(r.columns))
	for _, c := range r.columns {
		columns = append(columns, c.copy())
	}
	rows := make([]*Row, 0, len(r.rows))
	for _, row := range r.rows {
		rows = append(rows, row.copy())
	}
	return New(reportId, r.sobId, periodId, title, false, r.class, columns, rows)
}

func (r *Report) UpdateStructure(params UpdateParams) error {
	if params.Title != nil {
		r.title = *params.Title
	}
	if params.Rows != nil {
		rows, err := rebuildRows(r.rows, params.Rows)
		if err != nil {
			return err
		}
		r.rows = rows
	}
	return nil
}

func (r *Report) Id() uuid.UUID       { return r.id }
func (r *Report) SobId() uuid.UUID    { return r.sobId }
func (r *Report) PeriodId() uuid.UUID { return r.periodId }
func (r *Report) Title() string       { return r.title }
func (r *Report) Template() bool      { return r.template }
func (r *Report) Class() string       { return r.class }
func (r *Report) Columns() []*Column  { return r.columns }
func (r *Report) Rows() []*Row        { return r.rows }

func (c *Column) Id() uuid.UUID     { return c.id }
func (c *Column) Label() string     { return c.label }
func (c *Column) ValueType() string { return c.valueType }
func (c *Column) Sequence() int     { return c.sequence }

func (r *Row) Id() uuid.UUID                        { return r.id }
func (r *Row) RowCode() string                      { return r.rowCode }
func (r *Row) Text() string                         { return r.text }
func (r *Row) Sequence() int                        { return r.sequence }
func (r *Row) LineNo() *int                         { return r.lineNo }
func (r *Row) ShowLineNo() bool                     { return r.showLineNo }
func (r *Row) SumFactor() int                       { return r.sumFactor }
func (r *Row) DisplaySumFactor() bool               { return r.displaySumFactor }
func (r *Row) Indent() int                          { return r.indent }
func (r *Row) CanEdit() bool                        { return r.canEdit }
func (r *Row) CanMove() bool                        { return r.canMove }
func (r *Row) CanAddChild() bool                    { return r.canAddChild }
func (r *Row) Expression() *Expression              { return r.expression }
func (r *Row) Rows() []*Row                         { return r.rows }
func (r *Row) Amounts() []decimal.Decimal           { return r.amounts }
func (r *Row) SetAmounts(amounts []decimal.Decimal) { r.amounts = amounts }

func (e *Expression) Id() uuid.UUID                            { return e.id }
func (e *Expression) Kind() string                             { return e.kind }
func (e *Expression) LedgerAccounts() []LedgerAccountReference { return e.ledgerAccounts }
func (e *Expression) CashFlowItems() []CashFlowItemReference   { return e.cashFlowItems }
func (e *Expression) RowReferences() []RowReference            { return e.rowReferences }

func (c *Column) copy() *Column {
	col, _ := NewColumn(uuid.New(), c.label, c.valueType, c.sequence)
	return col
}

func (r *Row) copy() *Row {
	rows := make([]*Row, 0, len(r.rows))
	for _, row := range r.rows {
		rows = append(rows, row.copy())
	}
	row, _ := NewRow(
		uuid.New(),
		r.rowCode,
		r.text,
		r.sequence,
		copyIntPtr(r.lineNo),
		r.showLineNo,
		r.sumFactor,
		r.displaySumFactor,
		r.indent,
		r.canEdit,
		r.canMove,
		r.canAddChild,
		r.expression.copy(),
		rows,
		nil,
	)
	return row
}

func (e *Expression) copy() *Expression {
	ledgerAccounts := append([]LedgerAccountReference(nil), e.ledgerAccounts...)
	cashFlowItems := append([]CashFlowItemReference(nil), e.cashFlowItems...)
	rowReferences := append([]RowReference(nil), e.rowReferences...)
	expr, _ := NewExpression(uuid.New(), e.kind, ledgerAccounts, cashFlowItems, rowReferences)
	return expr
}

func copyIntPtr(v *int) *int {
	if v == nil {
		return nil
	}
	return new(*v)
}

func expressionsEqual(a, b *Expression) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.kind == b.kind &&
		reflect.DeepEqual(a.ledgerAccounts, b.ledgerAccounts) &&
		reflect.DeepEqual(a.cashFlowItems, b.cashFlowItems) &&
		reflect.DeepEqual(a.rowReferences, b.rowReferences)
}
