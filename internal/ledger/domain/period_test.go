package domain

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var periodId = uuid.New()

func TestNewPeriod(t *testing.T) {
	type args struct {
		id            uuid.UUID
		financialYear int
		number        int
		openingTime   time.Time
		endingTime    time.Time
		isClosed      bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "normal_success",
			args: args{
				id:            periodId,
				financialYear: 2021,
				number:        1,
				openingTime:   mustParseTime("2021-01-01T00:00:00Z"),
				endingTime:    mustParseTime("2021-02-01T00:00:00Z"),
				isClosed:      false,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.NoError(t, err)
				return false
			},
		},
		{
			name: "invalid_financial_year",
			args: args{
				id:            periodId,
				financialYear: 1900,
				number:        1,
				openingTime:   mustParseTime("2021-01-01T00:00:00Z"),
				endingTime:    mustParseTime("2021-02-01T00:00:00Z"),
				isClosed:      false,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Error(t, err)
				return false
			},
		},
		{
			name: "invalid_period_number",
			args: args{
				id:            periodId,
				financialYear: 2021,
				number:        18,
				openingTime:   mustParseTime("2021-01-01T00:00:00Z"),
				endingTime:    mustParseTime("2021-02-01T00:00:00Z"),
				isClosed:      false,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Error(t, err)
				return false
			},
		},
		{
			name: "invalid_time",
			args: args{
				id:            periodId,
				financialYear: 2021,
				number:        1,
				openingTime:   mustParseTime("2021-02-01T00:00:00Z"),
				endingTime:    mustParseTime("2021-01-01T00:00:00Z"),
				isClosed:      false,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Error(t, err)
				return false
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPeriod(tt.args.id, uuid.New(), uuid.Nil, tt.args.financialYear, tt.args.number, tt.args.openingTime, tt.args.endingTime, tt.args.isClosed)
			if !tt.wantErr(t, err, fmt.Sprintf("NewPeriod(%v, %v, %v, %v, %v, %v)", tt.args.id, tt.args.financialYear, tt.args.number, tt.args.openingTime, tt.args.endingTime, tt.args.isClosed)) {
				return
			}
			assert.NotNil(t, got)
		})
	}
}

func mustParseTime(timeStr string) time.Time {
	timeVal, _ := time.Parse(time.RFC3339, timeStr)
	return timeVal
}
