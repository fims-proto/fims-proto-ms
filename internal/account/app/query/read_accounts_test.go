package query

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/account/app/command"
	"github/fims-proto/fims-proto-ms/internal/sob/app/query"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestReadAccountsHandler_cutAccountNumber(t *testing.T) {
	t.Parallel()
	type fields struct {
		readModel  AccountsReadModel
		sobService command.SobService
	}
	type args struct {
		accountNumber      string
		accountCodeLengths []int
	}
	tests := []struct {
		name                string
		fields              fields
		args                args
		wantLevelNumber     int
		wantSuperiorNumbers []int
		wantErr             bool
	}{
		{
			name: "firstLevelNumber_success",
			fields: fields{
				mockReadModel{},
				mockSobService{},
			},
			args: args{
				accountNumber:      "1000",
				accountCodeLengths: []int{4, 3, 3},
			},
			wantLevelNumber:     1000,
			wantSuperiorNumbers: []int{},
			wantErr:             false,
		},
		{
			name: "secondLevelNumber_success",
			fields: fields{
				mockReadModel{},
				mockSobService{},
			},
			args: args{
				accountNumber:      "1000001",
				accountCodeLengths: []int{4, 3, 3},
			},
			wantLevelNumber:     1,
			wantSuperiorNumbers: []int{1000},
			wantErr:             false,
		},
		{
			name: "thirdLevelNumber_success",
			fields: fields{
				mockReadModel{},
				mockSobService{},
			},
			args: args{
				accountNumber:      "1000001002",
				accountCodeLengths: []int{4, 3, 3},
			},
			wantLevelNumber:     2,
			wantSuperiorNumbers: []int{1000, 1},
			wantErr:             false,
		},
		{
			name: "tooShort_error",
			fields: fields{
				mockReadModel{},
				mockSobService{},
			},
			args: args{
				accountNumber:      "100",
				accountCodeLengths: []int{4, 3, 3},
			},
			wantErr: true,
		},
		{
			name: "tooLong_error",
			fields: fields{
				mockReadModel{},
				mockSobService{},
			},
			args: args{
				accountNumber:      "10000001001",
				accountCodeLengths: []int{4, 3, 3},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := ReadAccountsHandler{
				readModel:  tt.fields.readModel,
				sobService: tt.fields.sobService,
			}
			levelNumber, superiorNumbers, err := h.cutAccountNumber(tt.args.accountNumber, tt.args.accountCodeLengths)
			if (err != nil) != tt.wantErr {
				t.Errorf("cutAccountNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if levelNumber != tt.wantLevelNumber {
				t.Errorf("cutAccountNumber() wantLevelNumber = %v, wantLevelNumber %v", levelNumber, tt.wantLevelNumber)
			}
			if !reflect.DeepEqual(superiorNumbers, tt.wantSuperiorNumbers) {
				t.Errorf("cutAccountNumber() wantSuperiorNumbers = %v, wantLevelNumber %v", superiorNumbers, tt.wantSuperiorNumbers)
			}
		})
	}
}

type mockReadModel struct{}

func (r mockReadModel) ReadAllAccounts(context.Context, uuid.UUID) ([]Account, error) {
	panic("implement me")
}

func (r mockReadModel) ReadById(context.Context, uuid.UUID) (Account, error) {
	panic("implement me")
}

func (r mockReadModel) ReadByIds(context.Context, []uuid.UUID) (map[uuid.UUID]Account, error) {
	panic("implement me")
}

func (r mockReadModel) ReadByAccountNumber(context.Context, uuid.UUID, int, []int) (Account, error) {
	panic("implement me")
}

type mockSobService struct{}

func (m mockSobService) ReadById(context.Context, uuid.UUID) (query.Sob, error) {
	panic("implement me")
}
