package command

import (
	"context"
	"testing"

	"github/fims-proto/fims-proto-ms/internal/account/domain/account"
	"github/fims-proto/fims-proto-ms/internal/account/domain/account_configuration"
	"github/fims-proto/fims-proto-ms/internal/account/domain/period"

	"github/fims-proto/fims-proto-ms/internal/account/app/service"

	"github/fims-proto/fims-proto-ms/internal/account/domain"
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
		accountEntries   []accountConfigurationEntry
		codeLengthLimits []int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*account_configuration.AccountConfiguration
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
				accountEntries: []accountConfigurationEntry{
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
			want:    []*account_configuration.AccountConfiguration{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := InitialAccountConfigurationHandler{
				repo:       tt.fields.repo,
				sobService: tt.fields.sobService,
			}
			got, err := h.prepareAccountConfigurations(tt.args.sobId, tt.args.accountEntries, tt.args.codeLengthLimits)
			if (err != nil) != tt.wantErr {
				t.Errorf("prepareAccounts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, 8, len(got))
			for _, accountConfiguration := range got {
				switch accountConfiguration.Title() {
				case "库存现金":
					assert.Equal(t, "1001", accountConfiguration.AccountNumber())
					assert.EqualValues(t, []int{1001}, accountConfiguration.NumberHierarchy())
					assert.Equal(t, 1, accountConfiguration.Level())
				case "银行存款":
					assert.Equal(t, "1002", accountConfiguration.AccountNumber())
					assert.EqualValues(t, []int{1002}, accountConfiguration.NumberHierarchy())
					assert.Equal(t, 1, accountConfiguration.Level())
				case "中国银行存款":
					assert.Equal(t, "1002001", accountConfiguration.AccountNumber())
					assert.EqualValues(t, []int{1002, 1}, accountConfiguration.NumberHierarchy())
					assert.Equal(t, 2, accountConfiguration.Level())
				case "招商银行存款":
					assert.Equal(t, "1002002", accountConfiguration.AccountNumber())
					assert.EqualValues(t, []int{1002, 2}, accountConfiguration.NumberHierarchy())
					assert.Equal(t, 2, accountConfiguration.Level())
				case "管理费用":
					assert.Equal(t, "6602", accountConfiguration.AccountNumber())
					assert.EqualValues(t, []int{6602}, accountConfiguration.NumberHierarchy())
					assert.Equal(t, 1, accountConfiguration.Level())
				case "办公费":
					assert.Equal(t, "6602001", accountConfiguration.AccountNumber())
					assert.EqualValues(t, []int{6602, 1}, accountConfiguration.NumberHierarchy())
					assert.Equal(t, 2, accountConfiguration.Level())
				case "办公室租金":
					assert.Equal(t, "6602001001", accountConfiguration.AccountNumber())
					assert.EqualValues(t, []int{6602, 1, 1}, accountConfiguration.NumberHierarchy())
					assert.Equal(t, 3, accountConfiguration.Level())
				case "文具费用":
					assert.Equal(t, "6602001002", accountConfiguration.AccountNumber())
					assert.EqualValues(t, []int{6602, 1, 2}, accountConfiguration.NumberHierarchy())
					assert.Equal(t, 3, accountConfiguration.Level())
				}
			}
		})
	}
}

type mockRepo struct{}

func (m mockRepo) InitialAccountConfiguration(context.Context, []*account_configuration.AccountConfiguration) error {
	panic("implement me")
}

func (m mockRepo) CreatePeriod(context.Context, *period.Period, func() error) error {
	panic("implement me")
}

func (m mockRepo) CreateAccounts(context.Context, []*account.Account) error {
	panic("implement me")
}

func (m mockRepo) UpdateAccountsByPeriodAndIds(context.Context, uuid.UUID, []uuid.UUID, func(accounts []*account.Account) ([]*account.Account, error)) error {
	panic("implement me")
}

func (m mockRepo) Migrate(context.Context) error {
	panic("implement me")
}

type mockSobService struct{}

func (m mockSobService) ReadById(context.Context, uuid.UUID) (query.Sob, error) {
	panic("implement me")
}
