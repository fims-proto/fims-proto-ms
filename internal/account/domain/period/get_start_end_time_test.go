package period

import (
	"reflect"
	"testing"
	"time"
)

func Test_getOpeningTime(t *testing.T) {
	type args struct {
		fiscalYear   int
		periodNumber int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal_success",
			args: args{
				fiscalYear:   2023,
				periodNumber: 1,
			},
			want: "2023-01-01T00:00:00+08:00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getOpeningTime(tt.args.fiscalYear, tt.args.periodNumber).Format(time.RFC3339); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getOpeningTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getEndingTime(t *testing.T) {
	type args struct {
		fiscalYear   int
		periodNumber int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal_success",
			args: args{
				fiscalYear:   2023,
				periodNumber: 1,
			},
			want: "2023-02-01T00:00:00+08:00",
		},
		{
			name: "nextYear_success",
			args: args{
				fiscalYear:   2023,
				periodNumber: 12,
			},
			want: "2024-01-01T00:00:00+08:00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getEndingTime(tt.args.fiscalYear, tt.args.periodNumber).Format(time.RFC3339); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getEndingTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
