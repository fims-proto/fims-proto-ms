package http

import (
	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/report/app/command"
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
