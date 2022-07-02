package query

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github/fims-proto/fims-proto-ms/internal/account/app/query"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	retrievedLedger = Ledger{
		Id:       uuid.UUID{},
		PeriodId: uuid.UUID{},
		Account: Account{
			Id:                uuid.UUID{},
			SuperiorAccountId: uuid.UUID{},
			AccountNumber:     "",
			Title:             "",
			AccountType:       0,
			BalanceDirection:  0,
		},
		OpeningBalance: decimal.Decimal{},
		EndingBalance:  decimal.Decimal{},
		Debit:          decimal.Decimal{},
		Credit:         decimal.Decimal{},
		CreatedAt:      time.Time{},
		UpdatedAt:      time.Time{},
	}

	wantLedger = Ledger{
		Id:       uuid.UUID{},
		PeriodId: uuid.UUID{},
		Account: Account{
			Id:                uuid.UUID{},
			SuperiorAccountId: uuid.UUID{},
			AccountNumber:     "12345",
			Title:             "",
			AccountType:       0,
			BalanceDirection:  0,
		},
		OpeningBalance: decimal.Decimal{},
		EndingBalance:  decimal.Decimal{},
		Debit:          decimal.Decimal{},
		Credit:         decimal.Decimal{},
		CreatedAt:      time.Time{},
		UpdatedAt:      time.Time{},
	}
)

type mockReadService struct{}

func (m mockReadService) ReadLedgerById(context.Context, uuid.UUID) (Ledger, error) {
	panic("implement me")
}

func (m mockReadService) ReadAllLedgersByAccountingPeriod(context.Context, uuid.UUID, data.Pageable) (data.Page[Ledger], error) {
	return data.NewPage([]Ledger{retrievedLedger}, 1, 1, 1)
}

func (m mockReadService) ReadAllAccountingPeriods(context.Context, uuid.UUID, data.Pageable) (data.Page[AccountingPeriod], error) {
	panic("implement me")
}

func (m mockReadService) ReadAccountingPeriodById(context.Context, uuid.UUID) (AccountingPeriod, error) {
	panic("implement me")
}

func (m mockReadService) ReadOpenAccountingPeriod(context.Context, uuid.UUID) (AccountingPeriod, error) {
	panic("implement me")
}

func (m mockReadService) ReadLedgerLogsByAccountIdsAndTimes(context.Context, []uuid.UUID, time.Time, time.Time) ([]LedgerLog, error) {
	panic("implement me")
}

type mockAccountService struct{}

func (m mockAccountService) ReadSuperiorAccountIds(context.Context, uuid.UUID) ([]uuid.UUID, error) {
	panic("implement me")
}

func (m mockAccountService) ReadAccountsByIds(context.Context, []uuid.UUID) (map[uuid.UUID]query.Account, error) {
	return map[uuid.UUID]query.Account{
		{}: {
			Id:                uuid.UUID{},
			SobId:             uuid.UUID{},
			SuperiorAccountId: uuid.UUID{},
			NumberHierarchy:   nil,
			AccountNumber:     "12345",
			Title:             "",
			AccountType:       0,
			BalanceDirection:  0,
			SuperiorAccount:   nil,
		},
	}, nil
}

func (m mockAccountService) ReadAllAccountIdsBySobId(context.Context, uuid.UUID) ([]uuid.UUID, error) {
	panic("implement me")
}

func TestReadLedgerHandler_HandleReadAllLedgersByAccountingPeriod(t *testing.T) {
	type fields struct {
		readModel      LedgerReadModel
		accountService AccountService
	}
	type args struct {
		ctx      context.Context
		periodId uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Ledger
		wantErr bool
	}{
		{
			name: "normal_success",
			fields: fields{
				readModel:      mockReadService{},
				accountService: mockAccountService{},
			},
			args: args{
				ctx:      context.Background(),
				periodId: uuid.UUID{},
			},
			want:    []Ledger{wantLedger},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := ReadLedgerHandler{
				readModel:      tt.fields.readModel,
				accountService: tt.fields.accountService,
			}
			got, err := h.HandleReadAllLedgersByAccountingPeriod(tt.args.ctx, tt.args.periodId, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleReadAllLedgersByAccountingPeriod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Content, tt.want) {
				t.Errorf("HandleReadAllLedgersByAccountingPeriod() got = %v, want %v", got, tt.want)
			}
		})
	}
}
