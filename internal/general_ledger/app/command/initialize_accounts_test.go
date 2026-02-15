package command

import (
	"context"
	"testing"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_account"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_category"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_ledger"

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

var sampleAccountEntries = []accountEntry{
	{
		number:           "1001",
		level:            1,
		title:            "库存现金",
		superiorNumber:   "",
		class:            1,
		group:            101,
		balanceDirection: "debit",
	},
	{
		number:           "1002",
		level:            1,
		title:            "银行存款",
		superiorNumber:   "",
		class:            1,
		group:            101,
		balanceDirection: "debit",
	},
	{
		number:           "1002001",
		level:            2,
		title:            "中国银行存款",
		superiorNumber:   "1002",
		class:            1,
		group:            101,
		balanceDirection: "debit",
	},
	{
		number:           "1002002",
		level:            2,
		title:            "招商银行存款",
		superiorNumber:   "1002",
		class:            1,
		group:            101,
		balanceDirection: "debit",
	},
	{
		number:           "6602",
		level:            1,
		title:            "管理费用",
		superiorNumber:   "",
		class:            5,
		group:            501,
		balanceDirection: "not_defined",
	},
	{
		number:           "6602001",
		level:            2,
		title:            "办公费",
		superiorNumber:   "6602",
		class:            5,
		group:            503,
		balanceDirection: "not_defined",
	},
	{
		number:           "6602001001",
		level:            3,
		title:            "办公室租金",
		superiorNumber:   "6602001",
		class:            5,
		group:            503,
		balanceDirection: "not_defined",
	},
	{
		number:           "6602001002",
		level:            3,
		title:            "文具费用",
		superiorNumber:   "6602001",
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
		sobId            uuid.UUID
		accountEntries   []accountEntry
		codeLengthLimits []int
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
				sobId:            sobId,
				accountEntries:   sampleAccountEntries,
				codeLengthLimits: []int{4, 3, 3},
			},
			wantNumber: map[string]string{
				"库存现金":   "1001",
				"银行存款":   "1002",
				"中国银行存款": "1002001",
				"招商银行存款": "1002002",
				"管理费用":   "6602",
				"办公费":    "6602001",
				"办公室租金":  "6602001001",
				"文具费用":   "6602001002",
			},
			wantErr: false,
		},
		{
			name: "shorter_code_length_success",
			fields: fields{
				repo:       mockRepo{},
				sobService: mockSobService{},
			},
			args: args{
				sobId:            sobId,
				accountEntries:   sampleAccountEntries,
				codeLengthLimits: []int{4, 2, 2},
			},
			wantNumber: map[string]string{
				"库存现金":   "1001",
				"银行存款":   "1002",
				"中国银行存款": "100201",
				"招商银行存款": "100202",
				"管理费用":   "6602",
				"办公费":    "660201",
				"办公室租金":  "66020101",
				"文具费用":   "66020102",
			},
			wantErr: false,
		},
		{
			name: "longer_code_length_success",
			fields: fields{
				repo:       mockRepo{},
				sobService: mockSobService{},
			},
			args: args{
				sobId:            sobId,
				accountEntries:   sampleAccountEntries,
				codeLengthLimits: []int{4, 4, 4},
			},
			wantNumber: map[string]string{
				"库存现金":   "1001",
				"银行存款":   "1002",
				"中国银行存款": "10020001",
				"招商银行存款": "10020002",
				"管理费用":   "6602",
				"办公费":    "66020001",
				"办公室租金":  "660200010001",
				"文具费用":   "660200010002",
			},
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
					assert.Equal(t, tt.wantNumber["库存现金"], acc.AccountNumber())
					assert.EqualValues(t, []int{1001}, acc.NumberHierarchy())
					assert.Equal(t, 1, acc.Level())
					assert.True(t, acc.IsLeaf())
				case "银行存款":
					assert.Equal(t, tt.wantNumber["银行存款"], acc.AccountNumber())
					assert.EqualValues(t, []int{1002}, acc.NumberHierarchy())
					assert.Equal(t, 1, acc.Level())
					assert.False(t, acc.IsLeaf())
				case "中国银行存款":
					assert.Equal(t, tt.wantNumber["中国银行存款"], acc.AccountNumber())
					assert.EqualValues(t, []int{1002, 1}, acc.NumberHierarchy())
					assert.Equal(t, 2, acc.Level())
					assert.True(t, acc.IsLeaf())
				case "招商银行存款":
					assert.Equal(t, tt.wantNumber["招商银行存款"], acc.AccountNumber())
					assert.EqualValues(t, []int{1002, 2}, acc.NumberHierarchy())
					assert.Equal(t, 2, acc.Level())
					assert.True(t, acc.IsLeaf())
				case "管理费用":
					assert.Equal(t, tt.wantNumber["管理费用"], acc.AccountNumber())
					assert.EqualValues(t, []int{6602}, acc.NumberHierarchy())
					assert.Equal(t, 1, acc.Level())
					assert.False(t, acc.IsLeaf())
				case "办公费":
					assert.Equal(t, tt.wantNumber["办公费"], acc.AccountNumber())
					assert.EqualValues(t, []int{6602, 1}, acc.NumberHierarchy())
					assert.Equal(t, 2, acc.Level())
					assert.False(t, acc.IsLeaf())
				case "办公室租金":
					assert.Equal(t, tt.wantNumber["办公室租金"], acc.AccountNumber())
					assert.EqualValues(t, []int{6602, 1, 1}, acc.NumberHierarchy())
					assert.Equal(t, 3, acc.Level())
					assert.True(t, acc.IsLeaf())
				case "文具费用":
					assert.Equal(t, tt.wantNumber["文具费用"], acc.AccountNumber())
					assert.EqualValues(t, []int{6602, 1, 2}, acc.NumberHierarchy())
					assert.Equal(t, 3, acc.Level())
					assert.True(t, acc.IsLeaf())
				}
			}
		})
	}
}

type mockRepo struct{}

func (m mockRepo) ReadAuxiliaryAccountsByPairs(context.Context, uuid.UUID, []auxiliary_account.AuxiliaryPair) ([]*auxiliary_account.AuxiliaryAccount, error) {
	panic("implement me")
}

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

func (m mockRepo) CreateVoucher(context.Context, *voucher.Voucher) error {
	panic("implement me")
}

func (m mockRepo) UpdateVoucher(context.Context, uuid.UUID, func(v *voucher.Voucher) (*voucher.Voucher, error)) error {
	panic("implement me")
}

func (m mockRepo) ExistsVouchersNotPostedInPeriod(context.Context, uuid.UUID, uuid.UUID) (bool, error) {
	panic("implement me")
}

func (m mockRepo) CreateAuxiliaryCategories(context.Context, []*auxiliary_category.AuxiliaryCategory) error {
	panic("implement me")
}

func (m mockRepo) CreateAuxiliaryAccounts(context.Context, []*auxiliary_account.AuxiliaryAccount) error {
	panic("implement me")
}

func (m mockRepo) ReadAllAuxiliaryAccounts(context.Context, uuid.UUID) ([]*auxiliary_account.AuxiliaryAccount, error) {
	panic("implement me")
}

func (m mockRepo) CreateAuxiliaryLedgers(context.Context, []*auxiliary_ledger.AuxiliaryLedger) error {
	panic("implement me")
}

func (m mockRepo) UpsertAuxiliaryLedgersByPeriodAndAccounts(context.Context, uuid.UUID, uuid.UUID, []domain.AuxiliaryLedgerKey, func(auxiliaryLedgers []*auxiliary_ledger.AuxiliaryLedger) ([]*auxiliary_ledger.AuxiliaryLedger, error)) error {
	panic("implement me")
}

func (m mockRepo) ReadAuxiliaryLedgersByPeriod(context.Context, uuid.UUID) ([]*auxiliary_ledger.AuxiliaryLedger, error) {
	panic("implement me")
}

func (m mockRepo) ReadAuxiliaryLedgersByAccountAndPeriod(context.Context, uuid.UUID, uuid.UUID) ([]*auxiliary_ledger.AuxiliaryLedger, error) {
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

func (m mockRepo) ReadAccountByNumber(context.Context, uuid.UUID, string) (*account.Account, error) {
	panic("implement me")
}

func (m mockRepo) ReadAuxiliaryCategoryByKey(context.Context, uuid.UUID, string) (*auxiliary_category.AuxiliaryCategory, error) {
	panic("implement me")
}

func (m mockRepo) ReadAuxiliaryCategoriesByKeys(context.Context, uuid.UUID, []string) ([]*auxiliary_category.AuxiliaryCategory, error) {
	panic("implement me")
}

type mockSobService struct{}

func (m mockSobService) ReadById(context.Context, uuid.UUID) (query.Sob, error) {
	panic("implement me")
}
