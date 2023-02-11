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
