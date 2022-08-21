package sob

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSob_UpdateName(t *testing.T) {
	t.Parallel()
	type fields struct {
		id                  uuid.UUID
		name                string
		description         string
		baseCurrency        string
		startingPeriodYear  int
		startingPeriodMonth int
		accountsCodeLength  []int
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "normal_success",
			fields: fields{
				id:                  uuid.New(),
				name:                "name1",
				description:         "desc1",
				baseCurrency:        "RMB",
				startingPeriodYear:  2000,
				startingPeriodMonth: 1,
				accountsCodeLength:  []int{4, 2, 2},
			},
			args: args{
				name: "name2",
			},
			wantErr: false,
		},
		{
			name: "failed",
			fields: fields{
				id:                  uuid.New(),
				name:                "name1",
				description:         "desc1",
				baseCurrency:        "RMB",
				startingPeriodYear:  2000,
				startingPeriodMonth: 1,
				accountsCodeLength:  []int{4, 2, 2},
			},
			args: args{
				name: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Sob{
				id:                  tt.fields.id,
				name:                tt.fields.name,
				description:         tt.fields.description,
				baseCurrency:        tt.fields.baseCurrency,
				startingPeriodYear:  tt.fields.startingPeriodYear,
				startingPeriodMonth: tt.fields.startingPeriodMonth,
				accountsCodeLength:  tt.fields.accountsCodeLength,
			}
			if err := s.UpdateName(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("SobId.UpdateName() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				assert.Equal(t, tt.args.name, s.Name())
			}
		})
	}
}

func TestSob_UpdateAccountsCodeLength(t *testing.T) {
	t.Parallel()
	type fields struct {
		id                  uuid.UUID
		name                string
		description         string
		baseCurrency        string
		startingPeriodYear  int
		startingPeriodMonth int
		accountsCodeLength  []int
	}
	type args struct {
		accountsCodeLength []int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "normal_success",
			fields: fields{
				id:                  uuid.New(),
				name:                "name1",
				description:         "desc1",
				baseCurrency:        "RMB",
				startingPeriodYear:  2000,
				startingPeriodMonth: 1,
				accountsCodeLength:  []int{4, 2, 2},
			},
			args: args{
				accountsCodeLength: []int{4, 2, 2, 2},
			},
			wantErr: false,
		},
		{
			name: "shorten_level_error",
			fields: fields{
				id:                  uuid.New(),
				name:                "name1",
				description:         "desc1",
				baseCurrency:        "RMB",
				startingPeriodYear:  2000,
				startingPeriodMonth: 1,
				accountsCodeLength:  []int{4, 2, 2},
			},
			args: args{
				accountsCodeLength: []int{4, 2},
			},
			wantErr: true,
		},
		{
			name: "shorten_length_error",
			fields: fields{
				id:                  uuid.New(),
				name:                "name1",
				description:         "desc1",
				baseCurrency:        "RMB",
				startingPeriodYear:  2000,
				startingPeriodMonth: 1,
				accountsCodeLength:  []int{4, 2, 2},
			},
			args: args{
				accountsCodeLength: []int{4, 2, 1},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Sob{
				id:                  tt.fields.id,
				name:                tt.fields.name,
				description:         tt.fields.description,
				baseCurrency:        tt.fields.baseCurrency,
				startingPeriodYear:  tt.fields.startingPeriodYear,
				startingPeriodMonth: tt.fields.startingPeriodMonth,
				accountsCodeLength:  tt.fields.accountsCodeLength,
			}
			if err := s.UpdateAccountsCodeLength(tt.args.accountsCodeLength); (err != nil) != tt.wantErr {
				t.Errorf("SobId.UpdateAccountsCodeLength() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
