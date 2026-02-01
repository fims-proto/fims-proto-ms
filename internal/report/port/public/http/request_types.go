package http

import (
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/report/app/command"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/amount_type"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/data_source"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/formula_rule"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/item_type"

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
	Title       *string                  `json:"title,omitempty"`       // Optional: update report title
	AmountTypes []string                 `json:"amountTypes,omitempty"` // Optional: update amount types
	Sections    []UpdateSectionRequest   `json:"sections"`              // Required: complete section structure
}

type UpdateSectionRequest struct {
	Id    string                       `json:"id"`                // Section ID
	Title *string                      `json:"title,omitempty"`   // Optional: update section title
	Items []UpdateReportItemRequest    `json:"items"`             // Complete item list for this section
}

type UpdateReportItemRequest struct {
	// Identity
	Id *string `json:"id,omitempty"` // Existing item ID, or null/omit for new item

	// Position
	Sequence int `json:"sequence"` // Required: position within section (1-based)

	// Content (required for new items, optional for updates to existing items)
	Text             *string                      `json:"text,omitempty"`
	Level            *int                         `json:"level,omitempty"`
	SumFactor        *int                         `json:"sumFactor,omitempty"`
	DisplaySumFactor *bool                        `json:"displaySumFactor,omitempty"`
	ItemType         *string                      `json:"itemType,omitempty"`
	DataSource       *string                      `json:"dataSource,omitempty"`
	Formulas         []UpdateReportFormulaRequest `json:"formulas,omitempty"`
	IsBreakdownItem  *bool                        `json:"isBreakdownItem,omitempty"`
	IsAbleToAddChild *bool                        `json:"isAbleToAddChild,omitempty"`
}

type UpdateReportFormulaRequest struct {
	SumFactor     int    `json:"sumFactor" binding:"required"`
	AccountNumber string `json:"accountNumber" binding:"required"`
	Rule          string `json:"rule" binding:"required"`
}

type UpdateReportResponse struct {
	CreatedItemIds map[string]string `json:"createdItemIds,omitempty"` // Maps client temp ID -> actual UUID
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

	// Convert sections
	var sections []command.UpdateSectionData
	for _, sectionReq := range r.Sections {
		sectionId, err := uuid.Parse(sectionReq.Id)
		if err != nil {
			return command.UpdateReportCmd{}, fmt.Errorf("invalid sectionId: %s", sectionReq.Id)
		}

		// Convert items
		var items []command.UpdateItemData
		for _, itemReq := range sectionReq.Items {
			itemData, err := itemReq.toUpdateItemData(sobId)
			if err != nil {
				return command.UpdateReportCmd{}, err
			}
			items = append(items, itemData)
		}

		sections = append(sections, command.UpdateSectionData{
			SectionId: sectionId,
			Title:     sectionReq.Title,
			Items:     items,
		})
	}

	return command.UpdateReportCmd{
		ReportId:    reportId,
		SobId:       sobId,
		Title:       r.Title,
		AmountTypes: amountTypes,
		Sections:    sections,
	}, nil
}

func (r UpdateReportItemRequest) toUpdateItemData(sobId uuid.UUID) (command.UpdateItemData, error) {
	var itemId *uuid.UUID
	if r.Id != nil {
		parsed, err := uuid.Parse(*r.Id)
		if err != nil {
			return command.UpdateItemData{}, fmt.Errorf("invalid itemId: %s", *r.Id)
		}
		itemId = &parsed
	}

	// Convert formulas
	var formulas []command.FormulaData
	for _, f := range r.Formulas {
		rule, err := formula_rule.FromString(f.Rule)
		if err != nil {
			return command.UpdateItemData{}, err
		}
		formulas = append(formulas, command.FormulaData{
			SumFactor:     f.SumFactor,
			AccountNumber: f.AccountNumber,
			Rule:          rule,
		})
	}

	var itemType *item_type.ItemType
	if r.ItemType != nil {
		t, err := item_type.FromString(*r.ItemType)
		if err != nil {
			return command.UpdateItemData{}, err
		}
		itemType = &t
	}

	var dataSource *data_source.DataSource
	if r.DataSource != nil {
		ds, err := data_source.FromString(*r.DataSource)
		if err != nil {
			return command.UpdateItemData{}, err
		}
		dataSource = &ds
	}

	return command.UpdateItemData{
		ItemId:           itemId,
		Sequence:         r.Sequence,
		Text:             r.Text,
		Level:            r.Level,
		SumFactor:        r.SumFactor,
		DisplaySumFactor: r.DisplaySumFactor,
		ItemType:         itemType,
		DataSource:       dataSource,
		Formulas:         formulas,
		IsBreakdownItem:  r.IsBreakdownItem,
		IsAbleToAddChild: r.IsAbleToAddChild,
	}, nil
}
