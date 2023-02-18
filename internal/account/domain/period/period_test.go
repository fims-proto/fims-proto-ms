package period

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewByTime(t *testing.T) {
	type args struct {
		id        uuid.UUID
		sobId     uuid.UUID
		timePoint time.Time
	}
	tests := []struct {
		name    string
		args    args
		verify  func(t *testing.T, period *Period, err error)
		wantErr bool
	}{
		{
			name: "normal_success",
			args: args{
				id:        uuid.New(),
				sobId:     uuid.New(),
				timePoint: time.Date(time.Now().Year()+1, time.May, 1, 1, 1, 1, 1, time.UTC),
			},
			verify: func(t *testing.T, period *Period, _ error) {
				assert.Equal(t, time.Now().Year()+1, period.FiscalYear())
				assert.Equal(t, 5, period.PeriodNumber())
				assert.Equal(t, time.Date(period.FiscalYear(), time.May, 1, 0, 0, 0, 0, time.UTC), period.openingTime)
				assert.Equal(t, time.Date(period.FiscalYear(), time.June, 1, 0, 0, 0, 0, time.UTC), period.endingTime)
			},
			wantErr: false,
		},
		{
			name: "timeInPast_error",
			args: args{
				id:        uuid.New(),
				sobId:     uuid.New(),
				timePoint: time.Date(time.Now().Year()-1, time.May, 1, 1, 1, 1, 1, time.UTC),
			},
			verify: func(t *testing.T, _ *Period, err error) {
				assert.Equal(t, "period-timeInPast", err.Error())
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewByTime(tt.args.id, tt.args.sobId, tt.args.timePoint)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewByTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.verify(t, got, err)
		})
	}
}
