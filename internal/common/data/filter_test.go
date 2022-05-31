package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newFiltersFromQuery(t *testing.T) {
	type args struct {
		filter string
	}
	tests := []struct {
		name    string
		args    args
		want    func([]Filter, *testing.T)
		wantErr bool
	}{
		{
			name: "field eq value",
			args: args{
				filter: "field eq  value",
			},
			want: func(filters []Filter, t *testing.T) {
				assertions := assert.New(t)
				assertions.Equal(1, len(filters))
				assertions.Equal("field", filters[0].Field())
				assertions.Equal(OptEq, filters[0].Operator())
				assertions.Equal("value", filters[0].Values()[0])
			},
			wantErr: false,
		},
		{
			name: "field eq 'some thing' and fieldA eq value1",
			args: args{
				filter: "field  eq 'some thing'  and fieldA  eq value1",
			},
			want: func(filters []Filter, t *testing.T) {
				assertions := assert.New(t)
				assertions.Equal(2, len(filters))
				assertions.Equal("field", filters[0].Field())
				assertions.Equal(OptEq, filters[0].Operator())
				assertions.Equal("some thing", filters[0].Values()[0])

				assertions.Equal("field_a", filters[1].Field())
				assertions.Equal(OptEq, filters[1].Operator())
				assertions.Equal("value1", filters[1].Values()[0])
			},
			wantErr: false,
		},
		{
			name: "complex filter",
			args: args{
				filter: "field  eq value  and fieldA  bt value1, value2",
			},
			want: func(filters []Filter, t *testing.T) {
				assertions := assert.New(t)
				assertions.Equal(2, len(filters))
				assertions.Equal("field", filters[0].Field())
				assertions.Equal(OptEq, filters[0].Operator())
				assertions.Equal("value", filters[0].Values()[0])

				assertions.Equal("field_a", filters[1].Field())
				assertions.Equal(OptBt, filters[1].Operator())
				assertions.Equal(2, len(filters[1].Values()))
				assertions.Equal("value1", filters[1].Values()[0])
				assertions.Equal("value2", filters[1].Values()[1])
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newFiltersFromQuery(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("newFiltersFromQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.want(got, t)
		})
	}
}
