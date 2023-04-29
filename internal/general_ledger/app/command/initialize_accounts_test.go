package command

import (
	"context"
	"testing"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/voucher"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/ledger"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/period"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
	"github/fims-proto/fims-proto-ms/internal/sob/app/query"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var sobId = uuid.New()

func TestAccountDataLoadHandler_prepareAccounts(t *testing.T) {
	t.Parallel()
	type fields struct {
		repo       domain.Repository
		sobService service.SobService
	}
	type args struct {
		sobId            uuid.UUID
		accountEntries   []accountEntry
		codeLengthLimits []int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*account.Account
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				repo:       mockRepo{},
				sobService: mockSobService{},
			},
			args: args{
				sobId: sobId,
				accountEntries: []accountEntry{
					{
						number:           "1001",
						level:            1,
						title:            "库存现金",
						superiorNumber:   "",
						accountType:      "assets",
						balanceDirection: "debit",
					},
					{
						number:           "1002",
						level:            1,
						title:            "银行存款",
						superiorNumber:   "",
						accountType:      "assets",
						balanceDirection: "debit",
					},
					{
						number:           "1002001",
						level:            2,
						title:            "中国银行存款",
						superiorNumber:   "1002",
						accountType:      "assets",
						balanceDirection: "debit",
					},
					{
						number:           "1002002",
						level:            2,
						title:            "招商银行存款",
						superiorNumber:   "1002",
						accountType:      "assets",
						balanceDirection: "debit",
					},
					{
						number:           "6602",
						level:            1,
						title:            "管理费用",
						superiorNumber:   "",
						accountType:      "profit_and_loss",
						balanceDirection: "not_defined",
					},
					{
						number:           "6602001",
						level:            2,
						title:            "办公费",
						superiorNumber:   "6602",
						accountType:      "profit_and_loss",
						balanceDirection: "not_defined",
					},
					{
						number:           "6602001001",
						level:            3,
						title:            "办公室租金",
						superiorNumber:   "6602001",
						accountType:      "profit_and_loss",
						balanceDirection: "not_defined",
					},
					{
						number:           "6602001002",
						level:            3,
						title:            "文具费用",
						superiorNumber:   "6602001",
						accountType:      "profit_and_loss",
						balanceDirection: "not_defined",
					},
				},
				codeLengthLimits: []int{4, 3, 3},
			},
			want:    []*account.Account{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := prepareAccounts(tt.args.sobId, tt.args.accountEntries, tt.args.codeLengthLimits)
			if (err != nil) != tt.wantErr {
				t.Errorf("prepareAccounts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, 8, len(got))
			for _, acc := range got {
				switch acc.Title() {
				case "库存现金":
					assert.Equal(t, "1001", acc.AccountNumber())
					assert.EqualValues(t, []int{1001}, acc.NumberHierarchy())
					assert.Equal(t, 1, acc.Level())
				case "银行存款":
					assert.Equal(t, "1002", acc.AccountNumber())
					assert.EqualValues(t, []int{1002}, acc.NumberHierarchy())
					assert.Equal(t, 1, acc.Level())
				case "中国银行存款":
					assert.Equal(t, "1002001", acc.AccountNumber())
					assert.EqualValues(t, []int{1002, 1}, acc.NumberHierarchy())
					assert.Equal(t, 2, acc.Level())
				case "招商银行存款":
					assert.Equal(t, "1002002", acc.AccountNumber())
					assert.EqualValues(t, []int{1002, 2}, acc.NumberHierarchy())
					assert.Equal(t, 2, acc.Level())
				case "管理费用":
					assert.Equal(t, "6602", acc.AccountNumber())
					assert.EqualValues(t, []int{6602}, acc.NumberHierarchy())
					assert.Equal(t, 1, acc.Level())
				case "办公费":
					assert.Equal(t, "6602001", acc.AccountNumber())
					assert.EqualValues(t, []int{6602, 1}, acc.NumberHierarchy())
					assert.Equal(t, 2, acc.Level())
				case "办公室租金":
					assert.Equal(t, "6602001001", acc.AccountNumber())
					assert.EqualValues(t, []int{6602, 1, 1}, acc.NumberHierarchy())
					assert.Equal(t, 3, acc.Level())
				case "文具费用":
					assert.Equal(t, "6602001002", acc.AccountNumber())
					assert.EqualValues(t, []int{6602, 1, 2}, acc.NumberHierarchy())
					assert.Equal(t, 3, acc.Level())
				}
			}
		})
	}
}

type mockRepo struct{}

func (m mockRepo) InitialAccounts(context.Context, []*account.Account) error {
	panic("implement me")
}

func (m mockRepo) CreatePeriod(context.Context, *period.Period) error {
	panic("implement me")
}

func (m mockRepo) UpdatePeriod(context.Context, uuid.UUID, func(p *period.Period) (*period.Period, error)) error {
	panic("implement me")
}

func (m mockRepo) CreateLedgers(context.Context, []*ledger.Ledger) error {
	panic("implement me")
}

func (m mockRepo) UpdateLedgersByPeriodAndAccountIds(context.Context, uuid.UUID, []uuid.UUID, func(accounts []*ledger.Ledger) ([]*ledger.Ledger, error)) error {
	panic("implement me")
}

func (m mockRepo) CreateVoucher(context.Context, *voucher.Voucher) error {
	panic("implement me")
}

func (m mockRepo) UpdateVoucher(context.Context, uuid.UUID, func(d *voucher.Voucher) (*voucher.Voucher, error)) error {
	panic("implement me")
}

func (m mockRepo) EnableTx(context.Context, func(txCtx context.Context) error) error {
	panic("implement me")
}

func (m mockRepo) Migrate(context.Context) error {
	panic("implement me")
}

type mockSobService struct{}

func (m mockSobService) ReadById(context.Context, uuid.UUID) (query.Sob, error) {
	panic("implement me")
}
