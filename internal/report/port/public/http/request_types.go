package http

import (
	"github/fims-proto/fims-proto-ms/internal/report/app/command"

	"github.com/google/uuid"
)

type GenerateReportRequest struct {
	Title            string   `json:"title"`
	AmountTypes      []string `json:"amountTypes"`
	PeriodFiscalYear int      `json:"periodFiscalYear"`
	PeriodNumber     int      `json:"periodNumber"`
}

type UpdateItemRequest struct {
	Text       string                     `json:"text,omitempty"`
	SumFactor  *int                       `json:"sumFactor,omitempty"`
	DataSource string                     `json:"dataSource,omitempty"`
	Formulas   []UpdateItemFormulaRequest `json:"formulas,omitempty"`
}

type UpdateItemFormulaRequest struct {
	SumFactor     int    `json:"sumFactor"`
	AccountNumber string `json:"accountNumber"`
	Rule          string `json:"rule"`
}

type AddItemRequest struct {
	InsertAfterSequence *int                       `json:"insertAfterSequence,omitempty"`
	Text                string                     `json:"text"`
	Level               int                        `json:"level"`
	SumFactor           int                        `json:"sumFactor"`
	DataSource          string                     `json:"dataSource"`
	Formulas            []UpdateItemFormulaRequest `json:"formulas,omitempty"`
	IsBreakdownItem     bool                       `json:"isBreakdownItem,omitempty"`
	IsAbleToAddChild    bool                       `json:"isAbleToAddChild,omitempty"`
}

// mappers

func (r UpdateItemRequest) mapToCommand(sobId, itemId uuid.UUID) command.UpdateItemCmd {
	var formulaCmds []command.FormulaCmd
	for _, formulaReq := range r.Formulas {
		formulaCmds = append(formulaCmds, command.FormulaCmd{
			SumFactor:     formulaReq.SumFactor,
			AccountNumber: formulaReq.AccountNumber,
			Rule:          formulaReq.Rule,
		})
	}

	return command.UpdateItemCmd{
		SobId:      sobId,
		Id:         itemId,
		Text:       r.Text,
		SumFactor:  r.SumFactor,
		DataSource: r.DataSource,
		Formulas:   formulaCmds,
	}
}

func (r AddItemRequest) mapToCommand(sobId, reportId, sectionId uuid.UUID) command.AddItemCmd {
	var formulaCmds []command.FormulaCmd
	for _, formulaReq := range r.Formulas {
		formulaCmds = append(formulaCmds, command.FormulaCmd{
			SumFactor:     formulaReq.SumFactor,
			AccountNumber: formulaReq.AccountNumber,
			Rule:          formulaReq.Rule,
		})
	}

	insertAfterSequence := 0
	if r.InsertAfterSequence != nil {
		insertAfterSequence = *r.InsertAfterSequence
	}

	return command.AddItemCmd{
		SobId:               sobId,
		ReportId:            reportId,
		SectionId:           sectionId,
		InsertAfterSequence: insertAfterSequence,
		Text:                r.Text,
		Level:               r.Level,
		SumFactor:           r.SumFactor,
		DataSource:          r.DataSource,
		Formulas:            formulaCmds,
		IsBreakdownItem:     r.IsBreakdownItem,
		IsAbleToAddChild:    r.IsAbleToAddChild,
	}
}
