package command

import (
	"context"
	"testing"
	"time"

	"github/fims-proto/fims-proto-ms/internal/account/app/query"

	ledgerQuery "github/fims-proto/fims-proto-ms/internal/ledger/app/query"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestApp_HandleCreateVoucherHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		constructor func() *CreateVoucherCmd
	}{
		{
			name: "normal_success",
			constructor: func() *CreateVoucherCmd {
				return createVoucherCmd()
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			assertions := assert.New(t)

			cmd := test.constructor()

			repoMock := newVoucherRepoMock()
			accServiceMock := newAccountService()
			cntServiceMock := newCounterService()
			ledServiceMock := newLedgerService()
			handler := NewCreateVoucherHandler(repoMock, &accServiceMock, &cntServiceMock, &ledServiceMock)
			newUUID, err := handler.Handle(context.Background(), *cmd)

			assertions.NoError(err)
			vouchers := repoMock.vouchers

			d100, _ := decimal.NewFromString("100")

			v := vouchers[newUUID]

			assertions.Len(vouchers, 1)
			assertions.Len(v.LineItems(), 2)
			assertions.Equal(d100, v.Credit())
			assertions.Equal(d100, v.Debit())
			assertions.Equal(userA, v.Creator())
			assertions.True(accServiceMock.invoked)
		})
	}
}

func createVoucherCmd() *CreateVoucherCmd {
	lineItems := []LineItemCmd{
		{
			Summary:       "test_item1",
			AccountNumber: "1000",
			Debit:         decimal.RequireFromString("100"),
			Credit:        decimal.Zero,
		},
		{
			Summary:       "test_item2",
			AccountNumber: "2000",
			Debit:         decimal.Zero,
			Credit:        decimal.RequireFromString("100"),
		},
	}
	return &CreateVoucherCmd{
		SobId:              uuid.New(),
		VoucherType:        "GENERAL_VOUCHER",
		AttachmentQuantity: 0,
		LineItems:          lineItems,
		Creator:            userA,
		TransactionTime:    time.Now(),
	}
}

func prepareBalancedItems() []*domain.LineItem {
	item1, _ := domain.NewLineItem(uuid.New(), uuid.New(), "test", decimal.RequireFromString("100"), decimal.Zero)
	item2, _ := domain.NewLineItem(uuid.New(), uuid.New(), "test", decimal.RequireFromString("100"), decimal.Zero)
	item3, _ := domain.NewLineItem(uuid.New(), uuid.New(), "test", decimal.Zero, decimal.RequireFromString("150"))
	item4, _ := domain.NewLineItem(uuid.New(), uuid.New(), "test", decimal.Zero, decimal.RequireFromString("50"))
	return []*domain.LineItem{
		item1,
		item2,
		item3,
		item4,
	}
}

func newVoucherRepoMock() voucherRepoMock {
	return voucherRepoMock{vouchers: make(map[uuid.UUID]domain.Voucher)}
}

func newAccountService() accountServiceMock {
	return accountServiceMock{invoked: false}
}

func newCounterService() counterServiceMock {
	return counterServiceMock{invoked: false}
}

func newLedgerService() ledgerServiceMock {
	return ledgerServiceMock{}
}

type voucherRepoMock struct {
	vouchers map[uuid.UUID]domain.Voucher
}

func (r voucherRepoMock) CreateVoucher(_ context.Context, v *domain.Voucher) (uuid.UUID, error) {
	r.vouchers[v.Id()] = *v
	return v.Id(), nil
}

func (r voucherRepoMock) UpdateVoucher(_ context.Context, id uuid.UUID, updateFn func(voucher *domain.Voucher) (*domain.Voucher, error)) error {
	v, ok := r.vouchers[id]
	if !ok {
		return errors.Errorf("voucher %s not exists", id)
	}

	updatedVoucher, err := updateFn(&v)
	if err != nil {
		return errors.Wrapf(err, "voucher %s updated failed", id)
	}
	r.vouchers[id] = *updatedVoucher
	return nil
}

func (r voucherRepoMock) Migrate(context.Context) error {
	panic("not implemented")
}

type accountServiceMock struct {
	invoked bool
}

func (s *accountServiceMock) ValidateExistenceAndGetId(_ context.Context, _ uuid.UUID, accountNumbers []string) (map[string]uuid.UUID, error) {
	s.invoked = true
	accountIds := make(map[string]uuid.UUID)
	for _, accountNumber := range accountNumbers {
		accountIds[accountNumber] = uuid.New()
	}
	return accountIds, nil
}

func (s *accountServiceMock) ReadAccountsByIds(context.Context, []uuid.UUID) (map[uuid.UUID]query.Account, error) {
	panic("implement me")
}

type counterServiceMock struct {
	invoked bool
}

func (c *counterServiceMock) GetNextIdentifier(context.Context, ...string) (string, error) {
	c.invoked = true
	return "1", nil
}

type ledgerServiceMock struct{}

func (l ledgerServiceMock) ReadPeriodByTime(context.Context, uuid.UUID, time.Time) (ledgerQuery.Period, error) {
	return ledgerQuery.Period{
		Id:       uuid.New(),
		IsClosed: false,
	}, nil
}

func (l ledgerServiceMock) ReadPeriodsByIds(context.Context, []uuid.UUID) (map[uuid.UUID]ledgerQuery.Period, error) {
	panic("implement me")
}

func (l ledgerServiceMock) PostVoucher(context.Context, domain.Voucher) error {
	panic("implement me")
}
