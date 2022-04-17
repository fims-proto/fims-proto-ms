package query

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/account/app/query"
	"github/fims-proto/fims-proto-ms/internal/ledger/app/service"
	"reflect"
	"testing"
	"time"

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
			Level:             0,
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
			Level:             0,
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

func (m mockReadService) ReadLedgerById(ctx context.Context, id uuid.UUID) (Ledger, error) {
	panic("implement me")
}

func (m mockReadService) ReadAllLedgersByAccountingPeriod(ctx context.Context, periodId uuid.UUID) ([]Ledger, error) {
	return []Ledger{retrievedLedger}, nil
}

func (m mockReadService) ReadAllAccountingPeriods(ctx context.Context, sobId uuid.UUID) ([]AccountingPeriod, error) {
	panic("implement me")
}

func (m mockReadService) ReadAccountingPeriodById(ctx context.Context, id uuid.UUID) (AccountingPeriod, error) {
	panic("implement me")
}

func (m mockReadService) ReadOpenAccountingPeriod(ctx context.Context, sobId uuid.UUID) (AccountingPeriod, error) {
	panic("implement me")
}

func (m mockReadService) ReadLedgerLogsByAccountIdsAndTimes(ctx context.Context, accountId []uuid.UUID, openingTime, endingTime time.Time) ([]LedgerLog, error) {
	panic("implement me")
}

type mockAccountService struct{}

func (m mockAccountService) ReadSuperiorAccountIds(ctx context.Context, accountId uuid.UUID) ([]uuid.UUID, error) {
	panic("implement me")
}

func (m mockAccountService) ReadAccountsByIds(ctx context.Context, accountIds []uuid.UUID) (map[uuid.UUID]query.Account, error) {
	return map[uuid.UUID]query.Account{
		{}: {
			Id:                uuid.UUID{},
			SobId:             uuid.UUID{},
			SuperiorAccountId: uuid.UUID{},
			SuperiorNumbers:   nil,
			LevelNumber:       0,
			AccountNumber:     "12345",
			Title:             "",
			Level:             0,
			AccountType:       0,
			BalanceDirection:  0,
			SuperiorAccount:   nil,
		},
	}, nil
}

func (m mockAccountService) ReadAllAccountIdsBySobId(ctx context.Context, sobId uuid.UUID) ([]uuid.UUID, error) {
	panic("implement me")
}

func TestReadLedgerHandler_HandleReadAllLedgersByAccountingPeriod(t *testing.T) {
	type fields struct {
		readModel      LedgerReadModel
		accountService service.AccountService
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
			got, err := h.HandleReadAllLedgersByAccountingPeriod(tt.args.ctx, tt.args.periodId)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleReadAllLedgersByAccountingPeriod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleReadAllLedgersByAccountingPeriod() got = %v, want %v", got, tt.want)
			}
		})
	}
}
