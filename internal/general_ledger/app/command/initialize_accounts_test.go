package command

import (
	"context"
	"testing"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/journal"

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

var sampleAccountEntries = []accountEntry{
	{
		number:           "001001",
		level:            1,
		title:            "库存现金",
		superiorNumber:   "",
		class:            1,
		group:            101,
		balanceDirection: "debit",
	},
	{
		number:           "001002",
		level:            1,
		title:            "银行存款",
		superiorNumber:   "",
		class:            1,
		group:            101,
		balanceDirection: "debit",
	},
	{
		number:           "001002000001",
		level:            2,
		title:            "中国银行存款",
		superiorNumber:   "001002",
		class:            1,
		group:            101,
		balanceDirection: "debit",
	},
	{
		number:           "001002000002",
		level:            2,
		title:            "招商银行存款",
		superiorNumber:   "001002",
		class:            1,
		group:            101,
		balanceDirection: "debit",
	},
	{
		number:           "006602",
		level:            1,
		title:            "管理费用",
		superiorNumber:   "",
		class:            5,
		group:            501,
		balanceDirection: "not_defined",
	},
	{
		number:           "006602000001",
		level:            2,
		title:            "办公费",
		superiorNumber:   "006602",
		class:            5,
		group:            503,
		balanceDirection: "not_defined",
	},
	{
		number:           "006602000001000001",
		level:            3,
		title:            "办公室租金",
		superiorNumber:   "006602000001",
		class:            5,
		group:            503,
		balanceDirection: "not_defined",
	},
	{
		number:           "006602000001000002",
		level:            3,
		title:            "文具费用",
		superiorNumber:   "006602000001",
		class:            5,
		group:            503,
		balanceDirection: "not_defined",
	},
}

func TestAccountDataLoadHandler_prepareAccounts(t *testing.T) {
	t.Parallel()
	type fields struct {
		repo       domain.Repository
		sobService service.SobService
	}
	type args struct {
		sobId          uuid.UUID
		accountEntries []accountEntry
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantNumber map[string]string
		wantErr    bool
	}{
		{
			name: "success",
			fields: fields{
				repo:       mockRepo{},
				sobService: mockSobService{},
			},
			args: args{
				sobId:          sobId,
				accountEntries: sampleAccountEntries,
			},
			wantNumber: map[string]string{
				"库存现金":   "001001",
				"银行存款":   "001002",
				"中国银行存款": "001002000001",
				"招商银行存款": "001002000002",
				"管理费用":   "006602",
				"办公费":    "006602000001",
				"办公室租金":  "006602000001000001",
				"文具费用":   "006602000001000002",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := prepareAccounts(tt.args.sobId, tt.args.accountEntries)
			if (err != nil) != tt.wantErr {
				t.Errorf("prepareAccounts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, 8, len(got))
			for _, acc := range got {
				switch acc.Title() {
				case "库存现金":
					assert.Equal(t, tt.wantNumber["库存现金"], acc.RawAccountNumber())
					hierarchy, _ := account.HierarchyFromRaw(acc.RawAccountNumber())
					assert.EqualValues(t, []int{1001}, hierarchy)
					assert.Equal(t, 1, acc.Level())
					assert.True(t, acc.IsLeaf())
				case "银行存款":
					assert.Equal(t, tt.wantNumber["银行存款"], acc.RawAccountNumber())
					hierarchy, _ := account.HierarchyFromRaw(acc.RawAccountNumber())
					assert.EqualValues(t, []int{1002}, hierarchy)
					assert.Equal(t, 1, acc.Level())
					assert.False(t, acc.IsLeaf())
				case "中国银行存款":
					assert.Equal(t, tt.wantNumber["中国银行存款"], acc.RawAccountNumber())
					hierarchy, _ := account.HierarchyFromRaw(acc.RawAccountNumber())
					assert.EqualValues(t, []int{1002, 1}, hierarchy)
					assert.Equal(t, 2, acc.Level())
					assert.True(t, acc.IsLeaf())
				case "招商银行存款":
					assert.Equal(t, tt.wantNumber["招商银行存款"], acc.RawAccountNumber())
					hierarchy, _ := account.HierarchyFromRaw(acc.RawAccountNumber())
					assert.EqualValues(t, []int{1002, 2}, hierarchy)
					assert.Equal(t, 2, acc.Level())
					assert.True(t, acc.IsLeaf())
				case "管理费用":
					assert.Equal(t, tt.wantNumber["管理费用"], acc.RawAccountNumber())
					hierarchy, _ := account.HierarchyFromRaw(acc.RawAccountNumber())
					assert.EqualValues(t, []int{6602}, hierarchy)
					assert.Equal(t, 1, acc.Level())
					assert.False(t, acc.IsLeaf())
				case "办公费":
					assert.Equal(t, tt.wantNumber["办公费"], acc.RawAccountNumber())
					hierarchy, _ := account.HierarchyFromRaw(acc.RawAccountNumber())
					assert.EqualValues(t, []int{6602, 1}, hierarchy)
					assert.Equal(t, 2, acc.Level())
					assert.False(t, acc.IsLeaf())
				case "办公室租金":
					assert.Equal(t, tt.wantNumber["办公室租金"], acc.RawAccountNumber())
					hierarchy, _ := account.HierarchyFromRaw(acc.RawAccountNumber())
					assert.EqualValues(t, []int{6602, 1, 1}, hierarchy)
					assert.Equal(t, 3, acc.Level())
					assert.True(t, acc.IsLeaf())
				case "文具费用":
					assert.Equal(t, tt.wantNumber["文具费用"], acc.RawAccountNumber())
					hierarchy, _ := account.HierarchyFromRaw(acc.RawAccountNumber())
					assert.EqualValues(t, []int{6602, 1, 2}, hierarchy)
					assert.Equal(t, 3, acc.Level())
					assert.True(t, acc.IsLeaf())
				}
			}
		})
	}
}

type mockRepo struct{}

func (m mockRepo) Migrate(context.Context) error {
	panic("implement me")
}

func (m mockRepo) EnableTx(context.Context, func(txCtx context.Context) error) error {
	panic("implement me")
}

func (m mockRepo) InitialAccounts(context.Context, []*account.Account) error {
	panic("implement me")
}

func (m mockRepo) UpdateAccount(context.Context, uuid.UUID, func(a *account.Account) (*account.Account, error)) error {
	panic("implement me")
}

func (m mockRepo) ReadAllAccounts(context.Context, uuid.UUID) ([]*account.Account, error) {
	panic("implement me")
}

func (m mockRepo) ReadAccountsByNumbers(context.Context, uuid.UUID, []string) ([]*account.Account, error) {
	panic("implement me")
}

func (m mockRepo) ReadAccountByRawNumber(context.Context, uuid.UUID, string) (*account.Account, error) {
	panic("implement me")
}

func (m mockRepo) ReadAccountsByRawNumbers(context.Context, uuid.UUID, []string) ([]*account.Account, error) {
	panic("implement me")
}

func (m mockRepo) ReadSuperiorAccountsById(context.Context, uuid.UUID) ([]*account.Account, error) {
	panic("implement me")
}

func (m mockRepo) CreatePeriodIfNotExists(context.Context, *period.Period) (*period.Period, bool, error) {
	panic("implement me")
}

func (m mockRepo) UpdatePeriod(context.Context, uuid.UUID, func(p *period.Period) (*period.Period, error)) error {
	panic("implement me")
}

func (m mockRepo) ReadCurrentPeriod(context.Context, uuid.UUID) (*period.Period, error) {
	panic("implement me")
}

func (m mockRepo) ReadPeriodById(context.Context, uuid.UUID, uuid.UUID) (*period.Period, error) {
	panic("implement me")
}

func (m mockRepo) ReadPreviousPeriod(context.Context, uuid.UUID) (*period.Period, error) {
	panic("implement me")
}

func (m mockRepo) CreateLedgers(context.Context, []*ledger.Ledger) error {
	panic("implement me")
}

func (m mockRepo) UpdateLedgersByPeriodAndAccountIds(context.Context, uuid.UUID, []uuid.UUID, func(accounts []*ledger.Ledger) ([]*ledger.Ledger, error)) error {
	panic("implement me")
}

func (m mockRepo) ReadLedgersByPeriod(context.Context, uuid.UUID) ([]*ledger.Ledger, error) {
	panic("implement me")
}

func (m mockRepo) ReadFirstLevelLedgersInPeriod(context.Context, uuid.UUID, uuid.UUID) ([]*ledger.Ledger, error) {
	panic("implement me")
}

func (m mockRepo) ExistsProfitAndLossLedgersHavingBalanceInPeriod(context.Context, uuid.UUID, uuid.UUID) (bool, error) {
	panic("implement me")
}

func (m mockRepo) ExistsLedgerHavingBalanceByRawAccountNumberInPeriod(context.Context, uuid.UUID, string, uuid.UUID) (bool, error) {
	panic("implement me")
}

func (m mockRepo) CreateJournal(context.Context, *journal.Journal) error {
	panic("implement me")
}

func (m mockRepo) UpdateJournalHeader(context.Context, uuid.UUID, func(j *journal.Journal) (*journal.Journal, error)) error {
	panic("implement me")
}

func (m mockRepo) UpdateEntireJournal(context.Context, uuid.UUID, func(j *journal.Journal) (*journal.Journal, error)) error {
	panic("implement me")
}

func (m mockRepo) ExistsJournalsNotPostedInPeriod(context.Context, uuid.UUID, uuid.UUID) (bool, error) {
	panic("implement me")
}

func (m mockRepo) ExistsJournalById(context.Context, uuid.UUID, uuid.UUID) (bool, error) {
	panic("implement me")
}

func (m mockRepo) ReadAccountsWithSuperiorsByIds(context.Context, uuid.UUID, []uuid.UUID) ([]*account.Account, error) {
	panic("implement me")
}

func (m mockRepo) ReadAllSubAccountsWithSuperiors(context.Context, uuid.UUID) ([]*account.Account, error) {
	panic("implement me")
}

func (m mockRepo) ReadFirstPeriod(context.Context, uuid.UUID) (*period.Period, error) {
	panic("implement me")
}

func (m mockRepo) CreateAccount(context.Context, *account.Account) error {
	panic("implement me")
}

func (m mockRepo) DeleteAccount(context.Context, uuid.UUID) error {
	panic("implement me")
}

func (m mockRepo) DeleteLedgersByAccountId(context.Context, uuid.UUID) error {
	panic("implement me")
}

func (m mockRepo) ReadAccountByNumber(context.Context, uuid.UUID, string) (*account.Account, error) {
	panic("implement me")
}

func (m mockRepo) ReadAccountById(context.Context, uuid.UUID) (*account.Account, error) {
	panic("implement me")
}

func (m mockRepo) ExistsChildAccountsByAccountId(context.Context, uuid.UUID) (bool, error) {
	panic("implement me")
}

func (m mockRepo) ExistsJournalLinesByAccountId(context.Context, uuid.UUID) (bool, error) {
	panic("implement me")
}

func (m mockRepo) ExistsLedgerWithOpeningBalanceByAccountId(context.Context, uuid.UUID) (bool, error) {
	panic("implement me")
}

func (m mockRepo) ReadProfitAndLossLedgersHavingBalanceInPeriod(context.Context, uuid.UUID, uuid.UUID) ([]*ledger.Ledger, error) {
	panic("implement me")
}

func (m mockRepo) ReadLedgerByRawAccountNumberInPeriod(context.Context, uuid.UUID, string, uuid.UUID) (*ledger.Ledger, error) {
	panic("implement me")
}

func (m mockRepo) ExistsClosingJournalInPeriod(context.Context, uuid.UUID, uuid.UUID, journal.JournalType) (bool, error) {
	panic("implement me")
}

type mockSobService struct{}

func (m mockSobService) ReadById(context.Context, uuid.UUID) (query.Sob, error) {
	panic("implement me")
}
