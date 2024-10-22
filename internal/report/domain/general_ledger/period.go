package general_ledger

type Period struct {
	fiscalYear   int
	periodNumber int
}

func NewPeriod(fiscalYear int, periodNumber int) *Period {
	return &Period{
		fiscalYear:   fiscalYear,
		periodNumber: periodNumber,
	}
}

func (p Period) FiscalYear() int {
	return p.fiscalYear
}

func (p Period) PeriodNumber() int {
	return p.periodNumber
}
