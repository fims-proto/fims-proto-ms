package period

import "time"

var fistDayOfMonth = 1

func getOpeningTime(fiscalYear, periodNumber int) time.Time {
	// first day of current month
	return time.Date(fiscalYear, time.Month(periodNumber), fistDayOfMonth, 0, 0, 0, 0, time.UTC)
}

func getEndingTime(fiscalYear, periodNumber int) time.Time {
	// first day of next month
	return getOpeningTime(fiscalYear, periodNumber).AddDate(0, 1, 0)
}

func (p *Period) NextNumber() (int, int) {
	firstDay := time.Date(p.fiscalYear, time.Month(p.periodNumber), fistDayOfMonth, 0, 0, 0, 0, time.UTC)
	nextFirstDay := firstDay.AddDate(0, 1, 0)

	return nextFirstDay.Year(), int(nextFirstDay.Month())
}

func (p *Period) PreviousNumber() (int, int) {
	firstDay := time.Date(p.fiscalYear, time.Month(p.periodNumber), fistDayOfMonth, 0, 0, 0, 0, time.UTC)
	previousFirstDay := firstDay.AddDate(0, -1, 0)

	return previousFirstDay.Year(), int(previousFirstDay.Month())
}
