package http

type GenerateReportRequest struct {
	PeriodFiscalYear int `json:"periodFiscalYear"`
	PeriodNumber     int `json:"periodNumber"`
}
