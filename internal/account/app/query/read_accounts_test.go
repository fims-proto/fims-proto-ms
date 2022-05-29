package query

import (
	"context"
	"reflect"
	"testing"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github/fims-proto/fims-proto-ms/internal/sob/app/query"

	"github.com/google/uuid"
)

var (
	retrievedAccount = Account{
		Id:                uuid.UUID{},
		SobId:             uuid.UUID{},
		SuperiorAccountId: uuid.UUID{},
		NumberHierarchy:   []int{1, 1},
		AccountNumber:     "",
		Title:             "",
		AccountType:       1,
		BalanceDirection:  1,
		SuperiorAccount:   nil,
	}
	wantAccount = Account{
		Id:                uuid.UUID{},
		SobId:             uuid.UUID{},
		SuperiorAccountId: uuid.UUID{},
		NumberHierarchy:   []int{1, 1},
		AccountNumber:     "0001001",
		Title:             "",
		AccountType:       1,
		BalanceDirection:  1,
		SuperiorAccount:   nil,
	}
)

type mockSobService struct{}

func (m mockSobService) ReadById(context.Context, uuid.UUID) (query.Sob, error) {
	return query.Sob{
		Id:                  uuid.UUID{},
		Name:                "",
		Description:         "",
		BaseCurrency:        "",
		StartingPeriodYear:  0,
		StartingPeriodMonth: 0,
		AccountsCodeLength:  []int{4, 3, 3, 3},
	}, nil
}

type mockReadModel struct{}

func (r mockReadModel) ReadAllAccounts(_ context.Context, _ uuid.UUID, pageable data.Pageable) (data.Page[Account], error) {
	page, _ := data.NewPage([]Account{retrievedAccount}, pageable.Page(), pageable.Size(), 1)
	return page, nil
}

func (r mockReadModel) ReadById(context.Context, uuid.UUID) (Account, error) {
	return retrievedAccount, nil
}

func (r mockReadModel) ReadByIds(context.Context, []uuid.UUID) (map[uuid.UUID]*Account, error) {
	return map[uuid.UUID]*Account{
		retrievedAccount.Id: &retrievedAccount,
	}, nil
}

func (r mockReadModel) ReadByAccountNumber(context.Context, uuid.UUID, []int) (Account, error) {
	panic("implement me")
}

func TestReadAccountsHandler_cutAccountNumber(t *testing.T) {
	t.Parallel()
	type args struct {
		accountNumber      string
		accountCodeLengths []int
	}
	tests := []struct {
		name                string
		args                args
		wantNumberHierarchy []int
		wantErr             bool
	}{
		{
			name: "firstLevelNumber_success",
			args: args{
				accountNumber:      "1000",
				accountCodeLengths: []int{4, 3, 3},
			},
			wantNumberHierarchy: []int{1000},
			wantErr:             false,
		},
		{
			name: "secondLevelNumber_success",
			args: args{
				accountNumber:      "1000001",
				accountCodeLengths: []int{4, 3, 3},
			},
			wantNumberHierarchy: []int{1000, 1},
			wantErr:             false,
		},
		{
			name: "thirdLevelNumber_success",
			args: args{
				accountNumber:      "1000001002",
				accountCodeLengths: []int{4, 3, 3},
			},
			wantNumberHierarchy: []int{1000, 1, 2},
			wantErr:             false,
		},
		{
			name: "tooShort_error",
			args: args{
				accountNumber:      "100",
				accountCodeLengths: []int{4, 3, 3},
			},
			wantErr: true,
		},
		{
			name: "tooLong_error",
			args: args{
				accountNumber:      "10000001001",
				accountCodeLengths: []int{4, 3, 3},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			numberHierarchy, err := cutAccountNumber(tt.args.accountNumber, tt.args.accountCodeLengths)
			if (err != nil) != tt.wantErr {
				t.Errorf("cutAccountNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(numberHierarchy, tt.wantNumberHierarchy) {
				t.Errorf("cutAccountNumber() wantNumberHierarchy = %v, wantLevelNumber %v", numberHierarchy, tt.wantNumberHierarchy)
			}
		})
	}
}

func TestReadAccountsHandler_concatenateAccountNumber(t *testing.T) {
	type args struct {
		numberHierarchy    []int
		accountCodeLengths []int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "normal_success",
			args: args{
				numberHierarchy:    []int{1001, 1, 1},
				accountCodeLengths: []int{4, 3, 2, 1},
			},
			want:    "100100101",
			wantErr: false,
		},
		{
			name: "noSuperior_success",
			args: args{
				numberHierarchy:    []int{1},
				accountCodeLengths: []int{4, 3, 3, 3, 3},
			},
			want:    "0001",
			wantErr: false,
		},
		{
			name: "accountCodeLengthsTooShort_error",
			args: args{
				numberHierarchy:    []int{1001, 1, 1},
				accountCodeLengths: []int{4, 3},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := concatenateAccountNumber(tt.args.numberHierarchy, tt.args.accountCodeLengths)
			if (err != nil) != tt.wantErr {
				t.Errorf("concatenateAccountNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("concatenateAccountNumber() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadAccountsHandler_HandleReadAll(t *testing.T) {
	pageRequest, _ := data.NewPageRequest(1, 1, nil, nil)
	type fields struct {
		readModel  AccountsReadModel
		sobService SobService
	}
	type args struct {
		ctx      context.Context
		sobId    uuid.UUID
		pageable data.Pageable
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Account
		wantErr bool
	}{
		{
			name: "normal_success",
			fields: fields{
				readModel:  mockReadModel{},
				sobService: mockSobService{},
			},
			args: args{
				ctx:      context.Background(),
				sobId:    uuid.UUID{},
				pageable: pageRequest,
			},
			want:    []Account{wantAccount},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := ReadAccountsHandler{
				readModel:  tt.fields.readModel,
				sobService: tt.fields.sobService,
			}
			got, err := h.HandleReadAll(tt.args.ctx, tt.args.sobId, tt.args.pageable)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleReadAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Content, tt.want) {
				t.Errorf("HandleReadAll() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadAccountsHandler_HandleReadById(t *testing.T) {
	type fields struct {
		readModel  AccountsReadModel
		sobService SobService
	}
	type args struct {
		ctx       context.Context
		accountId uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Account
		wantErr bool
	}{
		{
			name: "normal_success",
			fields: fields{
				readModel:  mockReadModel{},
				sobService: mockSobService{},
			},
			args: args{
				ctx:       context.Background(),
				accountId: uuid.UUID{},
			},
			want:    wantAccount,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := ReadAccountsHandler{
				readModel:  tt.fields.readModel,
				sobService: tt.fields.sobService,
			}
			got, err := h.HandleReadById(tt.args.ctx, tt.args.accountId)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleReadById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleReadById() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadAccountsHandler_HandleReadByIds(t *testing.T) {
	type fields struct {
		readModel  AccountsReadModel
		sobService SobService
	}
	type args struct {
		ctx        context.Context
		accountIds []uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[uuid.UUID]Account
		wantErr bool
	}{
		{
			name: "normal_success",
			fields: fields{
				readModel:  mockReadModel{},
				sobService: mockSobService{},
			},
			args: args{
				ctx:        context.Background(),
				accountIds: []uuid.UUID{{}},
			},
			want: map[uuid.UUID]Account{
				{}: wantAccount,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := ReadAccountsHandler{
				readModel:  tt.fields.readModel,
				sobService: tt.fields.sobService,
			}
			got, err := h.HandleReadByIds(tt.args.ctx, tt.args.accountIds)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleReadByIds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleReadByIds() got = %v, want %v", got, tt.want)
			}
		})
	}
}
