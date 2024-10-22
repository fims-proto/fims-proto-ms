package http

type GenerateReportRequest struct {
	Title            string   `json:"title"`
	AmountTypes      []string `json:"amountTypes"`
	PeriodFiscalYear int      `json:"periodFiscalYear"`
	PeriodNumber     int      `json:"periodNumber"`
}
