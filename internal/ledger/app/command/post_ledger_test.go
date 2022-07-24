package command

import (
	"context"
	"testing"
	"time"

	"github/fims-proto/fims-proto-ms/internal/ledger/app/service"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	accountQuery "github/fims-proto/fims-proto-ms/internal/account/app/query"
	"github/fims-proto/fims-proto-ms/internal/ledger/domain"
)

var (
	accountId1         = uuid.New()
	superiorAccountId1 = uuid.New()
	accountId2         = uuid.New()
	resultLedgers      []*domain.Ledger
)

func TestPostLedgersHandler_Handle(t *testing.T) {
	type fields struct {
		repo           domain.Repository
		accountService service.AccountService
	}
	type args struct {
		cmd PostLedgersCmd
	}

	mockRepoInstance := mockRepo{}
	mockAccountServiceInstance := mockAccountService{}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		verify  func(t *testing.T)
	}{
		{
			name: "normal_success",
			fields: fields{
				repo:           mockRepoInstance,
				accountService: mockAccountServiceInstance,
			},
			args: args{cmd: PostLedgersCmd{
				Accounts: map[uuid.UUID]PostLedgersAmountCmd{
					accountId1: {
						Debit:  decimal.RequireFromString("11"),
						Credit: decimal.RequireFromString("11.11"),
					},
					accountId2: {
						Debit:  decimal.RequireFromString("22"),
						Credit: decimal.RequireFromString("22.22"),
					},
				},
				PeriodId: uuid.New(),
			}},
			wantErr: false,
			verify: func(t *testing.T) {
				assert.Equal(t, 3, len(resultLedgers))

				assert.Equal(t, accountId1, resultLedgers[0].AccountId())
				assert.Equal(t, "11.11", resultLedgers[0].Credit().String())
				assert.Equal(t, "11", resultLedgers[0].Debit().String())
				assert.Equal(t, "0.11", resultLedgers[0].EndingBalance().String())

				assert.Equal(t, superiorAccountId1, resultLedgers[1].AccountId())
				assert.Equal(t, "11.11", resultLedgers[1].Credit().String())
				assert.Equal(t, "11", resultLedgers[1].Debit().String())
				assert.Equal(t, "0.11", resultLedgers[1].EndingBalance().String())

				assert.Equal(t, accountId2, resultLedgers[2].AccountId())
				assert.Equal(t, "22.22", resultLedgers[2].Credit().String())
				assert.Equal(t, "22", resultLedgers[2].Debit().String())
				assert.Equal(t, "-0.22", resultLedgers[2].EndingBalance().String())
			},
		},
	}
	for _, tt := range tests {
		resultLedgers = resultLedgers[:0]
		t.Run(tt.name, func(t *testing.T) {
			h := PostLedgersHandler{
				repo:           tt.fields.repo,
				accountService: tt.fields.accountService,
			}
			if err := h.Handle(context.Background(), tt.args.cmd); (err != nil) != tt.wantErr {
				t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				tt.verify(t)
			}
		})
	}
}

type mockRepo struct{}

func (m mockRepo) CreatePeriod(context.Context, *domain.Period) error {
	panic("implement me")
}

func (m mockRepo) CreateLedgers(context.Context, []*domain.Ledger) error {
	panic("implement me")
}

func (m mockRepo) UpdateLedgersByPeriodAndAccounts(_ context.Context, periodId uuid.UUID, accountIds []uuid.UUID, updateFn func(ledgers []*domain.Ledger) ([]*domain.Ledger, error)) error {
	var ledgers []*domain.Ledger
	for _, accountId := range accountIds {
		ledger, err := domain.NewLedger(uuid.New(), periodId, accountId, "dummy", decimal.Zero, decimal.Zero, decimal.Zero, decimal.Zero)
		if err != nil {
			return nil
		}
		ledgers = append(ledgers, ledger)
	}
	updatedLedgers, err := updateFn(ledgers)
	if err != nil {
		return err
	}

	resultLedgers = updatedLedgers
	return nil
}

func (m mockRepo) Migrate(context.Context) error {
	panic("implement me")
}

type mockAccountService struct{}

func (m mockAccountService) ReadAccountsBySobId(context.Context, uuid.UUID) ([]accountQuery.Account, error) {
	panic("implement me")
}

func (m mockAccountService) ReadAccountsByIds(context.Context, []uuid.UUID) (map[uuid.UUID]accountQuery.Account, error) {
	panic("implement me")
}

func (m mockAccountService) ReadAccountsWithSuperiorsByIds(_ context.Context, accountIds []uuid.UUID) ([]accountQuery.Account, error) {
	sobId := uuid.New()

	var accounts []accountQuery.Account
	for _, accountId := range accountIds {
		account := accountQuery.Account{
			Id:               accountId,
			SobId:            sobId,
			Title:            "dummy",
			AccountNumber:    "dummy",
			NumberHierarchy:  nil,
			AccountType:      1,
			BalanceDirection: 1,
			Level:            1,
			SuperiorAccount:  nil,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		if accountId == accountId1 {
			account.BalanceDirection = 2
			account.Level = 2
			account.SuperiorAccount = &accountQuery.Account{
				Id:               superiorAccountId1,
				SobId:            sobId,
				Title:            "dummy",
				AccountNumber:    "dummy",
				NumberHierarchy:  nil,
				AccountType:      1,
				BalanceDirection: 2,
				Level:            1,
				SuperiorAccount:  nil,
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			}
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}
