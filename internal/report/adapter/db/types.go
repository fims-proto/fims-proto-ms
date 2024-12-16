package db

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/shopspring/decimal"
	"github/fims-proto/fims-proto-ms/internal/common/data/converter"
	"github/fims-proto/fims-proto-ms/internal/report/app/query"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report"
)

type reportPO struct {
	Id          uuid.UUID  `gorm:"type:uuid;primaryKey"`
	SobId       uuid.UUID  `gorm:"type:uuid;uniqueIndex:UQ_Reports_SobId_PeriodId_Title"`
	PeriodId    *uuid.UUID `gorm:"type:uuid;uniqueIndex:UQ_Reports_SobId_PeriodId_Title"`
	Title       string     `gorm:"uniqueIndex:UQ_Reports_SobId_PeriodId_Title"`
	Template    bool
	Class       string
	AmountTypes pgtype.TextArray `gorm:"type:text[]"`
	Sections    []*sectionPO     `gorm:"foreignKey:ReportId"`

	Period *periodPO `gorm:"foreignKey:PeriodId"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

type sectionPO struct {
	Id        uuid.UUID  `gorm:"type:uuid;primaryKey"`
	ReportId  uuid.UUID  `gorm:"type:uuid"`
	SectionId *uuid.UUID `gorm:"type:uuid"`
	Title     string
	Sequence  int              // sequence within the whole report
	Amounts   pgtype.TextArray `gorm:"type:text[]"`
	Sections  []*sectionPO     `gorm:"foreignKey:SectionId"`
	Items     []*itemPO        `gorm:"foreignKey:SectionId"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

type itemPO struct {
	Id               uuid.UUID `gorm:"type:uuid;primaryKey"`
	SectionId        uuid.UUID `gorm:"type:uuid"`
	Text             string
	Level            int
	Sequence         int // sequence within the parent
	SumFactor        int
	DisplaySumFactor bool
	DataSource       string
	Formulas         []*formulaPO     `gorm:"foreignKey:ItemId"`
	Amounts          pgtype.TextArray `gorm:"type:text[]"`
	IsBreakdownItem  bool
	IsDeletable      bool
	IsTextModifiable bool
	IsDraggable      bool
	IsAbleToAddChild bool
	IsAbleToAddLeaf  bool

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

type formulaPO struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey"`
	ItemId    uuid.UUID `gorm:"type:uuid"`
	AccountId uuid.UUID `gorm:"type:uuid"`
	Sequence  int       // sequence within the parent
	SumFactor int
	Rule      string
	Amounts   pgtype.TextArray `gorm:"type:text[]"`

	Account accountPO `gorm:"foreignKey:AccountId"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

// table names

func (r reportPO) TableName() string {
	return "a_reports"
}

func (s sectionPO) TableName() string {
	return "a_sections"
}

func (i itemPO) TableName() string {
	return "a_items"
}

func (f formulaPO) TableName() string {
	return "a_formulas"
}

// schema

func (r reportPO) ResolveAssociation(entity string) (string, error) {
	if entity == "" {
		return r.TableName(), nil
	}
	if strings.EqualFold(entity, "sections") {
		return "Sections", nil
	}
	if strings.EqualFold(entity, "period") {
		return "Period", nil
	}
	return "", fmt.Errorf("reportPO doesn't have association named %s", entity)
}

func (s sectionPO) ResolveAssociation(entity string) (string, error) {
	if entity == "" {
		return s.TableName(), nil
	}
	if strings.EqualFold(entity, "sections") {
		return "Sections", nil
	}
	if strings.EqualFold(entity, "items") {
		return "Items", nil
	}
	return "", fmt.Errorf("sectionPO doesn't have association named %s", entity)
}

func (i itemPO) ResolveAssociation(entity string) (string, error) {
	if entity == "" {
		return i.TableName(), nil
	}
	if strings.EqualFold(entity, "formulas") {
		return "Formulas", nil
	}
	return "", fmt.Errorf("itemPO doesn't have association named %s", entity)
}

func (f formulaPO) ResolveAssociation(entity string) (string, error) {
	if entity == "" {
		return f.TableName(), nil
	}
	if strings.EqualFold(entity, "account") {
		return "Account", nil
	}
	return "", fmt.Errorf("formulaPO doesn't have association named %s", entity)
}

// mappers

func reportBOToPO(bo *report.Report) *reportPO {
	var sectionPOs []*sectionPO
	for _, section := range bo.Sections() {
		sectionPOs = append(sectionPOs, sectionBOToPO(section, bo.Id(), uuid.Nil))
	}

	var amountTypes []string
	for _, amountType := range bo.AmountTypes() {
		amountTypes = append(amountTypes, amountType.String())
	}

	var textArray pgtype.TextArray
	if err := textArray.Set(amountTypes); err != nil {
		panic(fmt.Errorf("failed to convert []string to TextArray: %w", err))
	}

	po := &reportPO{
		Id:          bo.Id(),
		SobId:       bo.SobId(),
		PeriodId:    converter.UUIDToPtr(bo.PeriodId()),
		Title:       bo.Title(),
		Template:    bo.Template(),
		Class:       bo.Class().String(),
		AmountTypes: textArray,
		Sections:    sectionPOs,
	}
	assignSequence(po.Sections)
	return po
}

func sectionBOToPO(bo *report.Section, reportId uuid.UUID, sectionId uuid.UUID) *sectionPO {
	var subSections []*sectionPO
	for _, subSection := range bo.Sections() {
		subSections = append(subSections, sectionBOToPO(subSection, reportId, bo.Id()))
	}

	var items []*itemPO
	for _, item := range bo.Items() {
		items = append(items, itemBOToPO(item, bo.Id()))
	}

	amounts, err := decimalArrayToTextArray(bo.Amounts())
	if err != nil {
		panic(fmt.Errorf("failed to convert []decimal.Decimal to TextArray: %w", err))
	}

	return &sectionPO{
		Id:        bo.Id(),
		ReportId:  reportId,
		SectionId: converter.UUIDToPtr(sectionId),
		Title:     bo.Title(),
		Amounts:   amounts,
		Sections:  subSections,
		Items:     items,
	}
}

func itemBOToPO(bo *report.Item, sectionId uuid.UUID) *itemPO {
	var formulas []*formulaPO
	for _, formula := range bo.Formulas() {
		formulas = append(formulas, formulaBOToPO(formula, bo.Id()))
	}

	amounts, err := decimalArrayToTextArray(bo.Amounts())
	if err != nil {
		panic(fmt.Errorf("failed to convert []decimal.Decimal to TextArray: %w", err))
	}

	return &itemPO{
		Id:               bo.Id(),
		SectionId:        sectionId,
		Text:             bo.Text(),
		Level:            bo.Level(),
		SumFactor:        bo.SumFactor(),
		DisplaySumFactor: bo.DisplaySumFactor(),
		DataSource:       bo.DataSource().String(),
		Formulas:         formulas,
		Amounts:          amounts,
		IsBreakdownItem:  bo.IsBreakdownItem(),
		IsDeletable:      bo.IsDeletable(),
		IsTextModifiable: bo.IsTextModifiable(),
		IsDraggable:      bo.IsDraggable(),
		IsAbleToAddChild: bo.IsAbleToAddChild(),
		IsAbleToAddLeaf:  bo.IsAbleToAddLeaf(),
	}
}

func formulaBOToPO(bo *report.Formula, itemId uuid.UUID) *formulaPO {
	amounts, err := decimalArrayToTextArray(bo.Amounts())
	if err != nil {
		panic(fmt.Errorf("failed to convert []decimal.Decimal to TextArray: %w", err))
	}

	return &formulaPO{
		Id:        bo.Id(),
		ItemId:    itemId,
		AccountId: bo.AccountId(),
		SumFactor: bo.SumFactor(),
		Rule:      bo.Rule().String(),
		Amounts:   amounts,
	}
}

func reportPOToBO(po *reportPO) (*report.Report, error) {
	// restore the section structure
	restructureAndSort(po)

	var amountTypes []string
	if err := po.AmountTypes.AssignTo(&amountTypes); err != nil {
		return nil, fmt.Errorf("failed to assign TextArray to []string: %w", err)
	}

	sections, err := converter.POsToBOs(po.Sections, sectionPOToBO)
	if err != nil {
		return nil, err
	}

	return report.New(
		po.Id,
		po.SobId,
		converter.UUIDFromPtr(po.PeriodId),
		po.Title,
		po.Template,
		po.Class,
		amountTypes,
		sections,
	)
}

func sectionPOToBO(po *sectionPO) (*report.Section, error) {
	subSections, err := converter.POsToBOs(po.Sections, sectionPOToBO)
	if err != nil {
		return nil, err
	}

	items, err := converter.POsToBOs(po.Items, itemPOToBO)
	if err != nil {
		return nil, err
	}

	amounts, err := textArrayToDecimalArray(po.Amounts)
	if err != nil {
		return nil, err
	}

	return report.NewSection(
		po.Id,
		po.Title,
		amounts,
		subSections,
		items,
	)
}

func itemPOToBO(po *itemPO) (*report.Item, error) {
	formulas, err := converter.POsToBOs(po.Formulas, formulaPOToBO)
	if err != nil {
		return nil, err
	}

	amounts, err := textArrayToDecimalArray(po.Amounts)
	if err != nil {
		return nil, err
	}

	return report.NewItem(
		po.Id,
		po.Text,
		po.Level,
		po.SumFactor,
		po.DisplaySumFactor,
		po.DataSource,
		formulas,
		amounts,
		po.IsBreakdownItem,
		po.IsDeletable,
		po.IsTextModifiable,
		po.IsDraggable,
		po.IsAbleToAddChild,
		po.IsAbleToAddLeaf,
	)
}

func formulaPOToBO(po *formulaPO) (*report.Formula, error) {
	amounts, err := textArrayToDecimalArray(po.Amounts)
	if err != nil {
		return nil, err
	}

	return report.NewFormula(
		po.Id,
		po.AccountId,
		po.SumFactor,
		po.Rule,
		amounts,
	)
}

func reportPOToDTO(po reportPO) query.Report {
	// restore the section structure
	restructureAndSort(&po)

	var amountTypes []string
	if err := po.AmountTypes.AssignTo(&amountTypes); err != nil {
		panic(fmt.Errorf("failed to assign TextArray to []string: %w", err))
	}

	return query.Report{
		Id:          po.Id,
		SobId:       po.SobId,
		Period:      periodPOToDTO(po.Period),
		Title:       po.Title,
		Template:    po.Template,
		Class:       po.Class,
		AmountTypes: amountTypes,
		Sections:    converter.POsToDTOs(po.Sections, sectionPOToDTO),
		CreatedAt:   po.CreatedAt,
		UpdatedAt:   po.UpdatedAt,
	}
}

func sectionPOToDTO(po *sectionPO) query.Section {
	amounts, err := textArrayToDecimalArray(po.Amounts)
	if err != nil {
		panic(fmt.Errorf("failed to assign TextArray to []decimal.Decimal: %w", err))
	}

	return query.Section{
		Id:       po.Id,
		Title:    po.Title,
		Amounts:  amounts,
		Sections: converter.POsToDTOs(po.Sections, sectionPOToDTO),
		Items:    converter.POsToDTOs(po.Items, itemPOToDTO),
	}
}

func itemPOToDTO(po *itemPO) query.Item {
	amounts, err := textArrayToDecimalArray(po.Amounts)
	if err != nil {
		panic(fmt.Errorf("failed to assign TextArray to []decimal.Decimal: %w", err))
	}

	return query.Item{
		Id:               po.Id,
		Text:             po.Text,
		Level:            po.Level,
		SumFactor:        po.SumFactor,
		DisplaySumFactor: po.DisplaySumFactor,
		DataSource:       po.DataSource,
		Formulas:         converter.POsToDTOs(po.Formulas, formulaPOToDTO),
		Amounts:          amounts,
		IsBreakdownItem:  po.IsBreakdownItem,
		IsDeletable:      po.IsDeletable,
		IsTextModifiable: po.IsTextModifiable,
		IsDraggable:      po.IsDraggable,
		IsAbleToAddChild: po.IsAbleToAddChild,
		IsAbleToAddLeaf:  po.IsAbleToAddLeaf,
	}
}

func formulaPOToDTO(po *formulaPO) query.Formula {
	amounts, err := textArrayToDecimalArray(po.Amounts)
	if err != nil {
		panic(fmt.Errorf("failed to assign TextArray to []decimal.Decimal: %w", err))
	}

	return query.Formula{
		Id:        po.Id,
		Account:   accountPOToDTO(po.Account),
		SumFactor: po.SumFactor,
		Rule:      po.Rule,
		Amounts:   amounts,
	}
}

func periodPOToDTO(po *periodPO) *query.Period {
	if po == nil {
		return nil
	}
	return &query.Period{
		FiscalYear:   po.FiscalYear,
		PeriodNumber: po.PeriodNumber,
	}
}

func accountPOToDTO(po accountPO) query.Account {
	var superiorAccountId *uuid.UUID
	if po.SuperiorAccountId != uuid.Nil {
		superiorAccountId = &po.SuperiorAccountId
	}

	return query.Account{
		Id:                po.Id,
		SobId:             po.SobId,
		SuperiorAccountId: superiorAccountId,
		Title:             po.Title,
		AccountNumber:     po.AccountNumber,
		Level:             po.Level,
		IsLeaf:            po.IsLeaf,
		Class:             po.Class,
		Group:             po.Group,
		BalanceDirection:  po.BalanceDirection,
	}
}

// assignSequence assigns correct Sequence field to section, item and formula
// in during runtime, the sequence is guaranteed by slice, however after it's saved into the database, Sequence field is how we can know the sequence
func assignSequence(sections []*sectionPO) {
	for i, section := range sections {
		section.Sequence = i

		// item and formula
		for j, item := range section.Items {
			item.Sequence = j
			for k, formula := range item.Formulas {
				formula.Sequence = k
			}
		}

		assignSequence(section.Sections)
	}
}

// restructureAndSort restores the correct sections level, and sort sections, items and formulas based on Sequence fields
// since sectionPO has both reportId and sectionId field, a subsection can have both of the fields.
// this could cause the reportPO having a flat section list (subsection is not nested in the higher level section), since sections is retrieved by foreign key reportId
func restructureAndSort(r *reportPO) {
	// restructure sections first
	sectionMap := make(map[uuid.UUID]*sectionPO)
	for _, section := range r.Sections {
		sectionMap[section.Id] = section
	}

	// assign subsections
	for _, section := range r.Sections {
		if section.SectionId != nil {
			if parentSection, ok := sectionMap[*section.SectionId]; ok {
				parentSection.Sections = append(parentSection.Sections, section)
			}
		}
	}

	// overwrite sections in report
	r.Sections = nil
	for section := range maps.Values(sectionMap) {
		if section.SectionId == nil {
			r.Sections = append(r.Sections, section)
		}
	}

	// sort
	for _, section := range r.Sections {
		sortRecursive(section)
	}
	slices.SortFunc(r.Sections, func(a, b *sectionPO) int { return a.Sequence - b.Sequence })
}

func sortRecursive(s *sectionPO) {
	// formulas
	for _, item := range s.Items {
		slices.SortFunc(item.Formulas, func(a, b *formulaPO) int { return a.Sequence - b.Sequence })
	}

	// items
	slices.SortFunc(s.Items, func(a, b *itemPO) int { return a.Sequence - b.Sequence })

	// sub sections
	slices.SortFunc(s.Sections, func(a, b *sectionPO) int { return a.Sequence - b.Sequence })
	for _, subSection := range s.Sections {
		sortRecursive(subSection)
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
		// null or undefined value
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
