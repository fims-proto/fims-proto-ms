package http

import (
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/report/app/command"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/amount_type"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/data_source"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/formula_rule"

	"github.com/google/uuid"
)

type GenerateReportRequest struct {
	Title            string   `json:"title"`
	AmountTypes      []string `json:"amountTypes"`
	PeriodFiscalYear int      `json:"periodFiscalYear"`
	PeriodNumber     int      `json:"periodNumber"`
}

// UpdateReportRequest represents the comprehensive update request for a report
type UpdateReportRequest struct {
	Title       *string                `json:"title,omitempty"`       // Optional: update report title
	AmountTypes []string               `json:"amountTypes,omitempty"` // Optional: update amount types
	Sections    []UpdateSectionRequest `json:"sections"`              // Required: complete section structure
}

type UpdateSectionRequest struct {
	Id       string                    `json:"id"`                 // Section ID
	Title    *string                   `json:"title,omitempty"`    // Optional: update section title
	Items    []UpdateReportItemRequest `json:"items"`              // Complete item list for this section
	Sections []UpdateSectionRequest    `json:"sections,omitempty"` // Optional: nested sections
}

type UpdateReportItemRequest struct {
	// Identity
	Id *string `json:"id,omitempty"` // Existing item ID, or null/omit for new item

	// Content (required for new items, optional for updates to existing items)
	Text             *string                      `json:"text,omitempty"`
	Level            *int                         `json:"level,omitempty"`
	SumFactor        *int                         `json:"sumFactor,omitempty"`
	DisplaySumFactor *bool                        `json:"displaySumFactor,omitempty"`
	DataSource       *string                      `json:"dataSource,omitempty"`
	Formulas         []UpdateReportFormulaRequest `json:"formulas,omitempty"`
	IsBreakdownItem  *bool                        `json:"isBreakdownItem,omitempty"`
	IsAbleToAddChild *bool                        `json:"isAbleToAddChild,omitempty"`
}

type UpdateReportFormulaRequest struct {
	Id            *string `json:"id,omitempty"`
	SumFactor     int     `json:"sumFactor" binding:"required"`
	AccountNumber string  `json:"accountNumber" binding:"required"`
	Rule          string  `json:"rule" binding:"required"`
}

// mapToCommand converts UpdateReportRequest to UpdateReportCmd
func (r UpdateReportRequest) mapToCommand(reportId uuid.UUID, sobId uuid.UUID) (command.UpdateReportCmd, error) {
	// Convert amount types
	var amountTypes []amount_type.AmountType
	for _, at := range r.AmountTypes {
		amountType, err := amount_type.FromString(at)
		if err != nil {
			return command.UpdateReportCmd{}, err
		}
		amountTypes = append(amountTypes, amountType)
	}

	// Convert sections recursively
	sections, err := convertSections(r.Sections, sobId)
	if err != nil {
		return command.UpdateReportCmd{}, err
	}

	return command.UpdateReportCmd{
		ReportId:    reportId,
		SobId:       sobId,
		Title:       r.Title,
		AmountTypes: amountTypes,
		Sections:    sections,
	}, nil
}

// convertSections recursively converts sections and their nested sections
func convertSections(sectionsReq []UpdateSectionRequest, sobId uuid.UUID) ([]command.UpdateReportCmdSection, error) {
	var sections []command.UpdateReportCmdSection
	for _, sectionReq := range sectionsReq {
		sectionId, err := uuid.Parse(sectionReq.Id)
		if err != nil {
			return nil, fmt.Errorf("invalid sectionId: %s", sectionReq.Id)
		}

		// Convert items
		var items []command.UpdateReportCmdItem
		for _, itemReq := range sectionReq.Items {
			itemData, err := itemReq.toUpdateItemData()
			if err != nil {
				return nil, err
			}
			items = append(items, itemData)
		}

		// Recursively convert nested sections
		var nestedSections []command.UpdateReportCmdSection
		if len(sectionReq.Sections) > 0 {
			nestedSections, err = convertSections(sectionReq.Sections, sobId)
			if err != nil {
				return nil, err
			}
		}

		sections = append(sections, command.UpdateReportCmdSection{
			SectionId: sectionId,
			Title:     sectionReq.Title,
			Items:     items,
			Sections:  nestedSections,
		})
	}

	return sections, nil
}

func (r UpdateReportItemRequest) toUpdateItemData() (command.UpdateReportCmdItem, error) {
	var itemId *uuid.UUID
	if r.Id != nil {
		parsed, err := uuid.Parse(*r.Id)
		if err != nil {
			return command.UpdateReportCmdItem{}, fmt.Errorf("invalid itemId: %s", *r.Id)
		}
		itemId = &parsed
	}

	// Convert formulas
	var formulas []command.UpdateReportCmdFormula
	for _, f := range r.Formulas {
		var formulaId *uuid.UUID
		if f.Id != nil {
			parsed, err := uuid.Parse(*f.Id)
			if err != nil {
				return command.UpdateReportCmdItem{}, fmt.Errorf("invalid formulaId: %s", *f.Id)
			}
			formulaId = &parsed
		}
		rule, err := formula_rule.FromString(f.Rule)
		if err != nil {
			return command.UpdateReportCmdItem{}, err
		}
		formulas = append(formulas, command.UpdateReportCmdFormula{
			FormulaId:     formulaId,
			SumFactor:     f.SumFactor,
			AccountNumber: f.AccountNumber,
			Rule:          rule,
		})
	}

	var dataSource *data_source.DataSource
	if r.DataSource != nil {
		ds, err := data_source.FromString(*r.DataSource)
		if err != nil {
			return command.UpdateReportCmdItem{}, err
		}
		dataSource = &ds
	}

	return command.UpdateReportCmdItem{
		ItemId:           itemId,
		Text:             r.Text,
		Level:            r.Level,
		SumFactor:        r.SumFactor,
		DisplaySumFactor: r.DisplaySumFactor,
		DataSource:       dataSource,
		Formulas:         formulas,
		IsBreakdownItem:  r.IsBreakdownItem,
		IsAbleToAddChild: r.IsAbleToAddChild,
	}, nil
}
