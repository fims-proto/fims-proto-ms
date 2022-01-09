package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/account/domain"
	"github/fims-proto/fims-proto-ms/internal/sob/app/query"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var sobId = uuid.New()

func TestAccountDataLoadHandler_prepareAccounts(t *testing.T) {
	t.Parallel()
	type fields struct {
		repo       domain.Repository
		sobService SobService
	}
	type args struct {
		sobId            uuid.UUID
		accountEntries   []accountDataLoadEntry
		codeLengthLimits []int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*domain.Account
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
				accountEntries: []accountDataLoadEntry{
					{
						number:           "1001",
						level:            1,
						title:            "库存现金",
						superiorNumber:   "",
						accountType:      "ASSETS",
						balanceDirection: "DEBIT",
					},
					{
						number:           "1002",
						level:            1,
						title:            "银行存款",
						superiorNumber:   "",
						accountType:      "ASSETS",
						balanceDirection: "DEBIT",
					},
					{
						number:           "1002001",
						level:            2,
						title:            "中国银行存款",
						superiorNumber:   "1002",
						accountType:      "ASSETS",
						balanceDirection: "DEBIT",
					},
					{
						number:           "1002002",
						level:            2,
						title:            "招商银行存款",
						superiorNumber:   "1002",
						accountType:      "ASSETS",
						balanceDirection: "DEBIT",
					},
					{
						number:           "6602",
						level:            1,
						title:            "管理费用",
						superiorNumber:   "",
						accountType:      "PROFIT_AND_LOSS",
						balanceDirection: "NOT_DEFINED",
					},
					{
						number:           "6602001",
						level:            2,
						title:            "办公费",
						superiorNumber:   "6602",
						accountType:      "PROFIT_AND_LOSS",
						balanceDirection: "NOT_DEFINED",
					},
					{
						number:           "6602001001",
						level:            3,
						title:            "办公室租金",
						superiorNumber:   "6602001",
						accountType:      "PROFIT_AND_LOSS",
						balanceDirection: "NOT_DEFINED",
					},
					{
						number:           "6602001002",
						level:            3,
						title:            "文具费用",
						superiorNumber:   "6602001",
						accountType:      "PROFIT_AND_LOSS",
						balanceDirection: "NOT_DEFINED",
					},
				},
				codeLengthLimits: []int{4, 3, 3},
			},
			want:    []*domain.Account{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := AccountDataLoadHandler{
				repo:       tt.fields.repo,
				sobService: tt.fields.sobService,
			}
			got, err := h.prepareAccounts(tt.args.sobId, tt.args.accountEntries, tt.args.codeLengthLimits)
			if (err != nil) != tt.wantErr {
				t.Errorf("prepareAccounts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, 8, len(got))
			for _, account := range got {
				switch account.Title() {
				case "库存现金":
					assert.Equal(t, 1001, account.LevelNumber())
					assert.Equal(t, 1, account.Level())
					assert.EqualValues(t, []int{}, account.SuperiorNumbers())
				case "银行存款":
					assert.Equal(t, 1002, account.LevelNumber())
					assert.Equal(t, 1, account.Level())
					assert.EqualValues(t, []int{}, account.SuperiorNumbers())
				case "中国银行存款":
					assert.Equal(t, 1, account.LevelNumber())
					assert.Equal(t, 2, account.Level())
					assert.EqualValues(t, []int{1002}, account.SuperiorNumbers())
				case "招商银行存款":
					assert.Equal(t, 2, account.LevelNumber())
					assert.Equal(t, 2, account.Level())
					assert.EqualValues(t, []int{1002}, account.SuperiorNumbers())
				case "管理费用":
					assert.Equal(t, 6602, account.LevelNumber())
					assert.Equal(t, 1, account.Level())
					assert.EqualValues(t, []int{}, account.SuperiorNumbers())
				case "办公费":
					assert.Equal(t, 1, account.LevelNumber())
					assert.Equal(t, 2, account.Level())
					assert.EqualValues(t, []int{6602}, account.SuperiorNumbers())
				case "办公室租金":
					assert.Equal(t, 1, account.LevelNumber())
					assert.Equal(t, 3, account.Level())
					assert.EqualValues(t, []int{6602, 1}, account.SuperiorNumbers())
				case "文具费用":
					assert.Equal(t, 2, account.LevelNumber())
					assert.Equal(t, 3, account.Level())
					assert.EqualValues(t, []int{6602, 1}, account.SuperiorNumbers())
				}
			}
		})
	}
}

type mockRepo struct{}

func (m mockRepo) CreateAccount(context.Context, *domain.Account) error {
	panic("implement me")
}

func (m mockRepo) DataLoad(context.Context, []*domain.Account) error {
	panic("implement me")
}

func (m mockRepo) Migrate(context.Context) error {
	panic("implement me")
}

type mockSobService struct{}

func (m mockSobService) ReadById(context.Context, uuid.UUID) (query.Sob, error) {
	panic("implement me")
}
