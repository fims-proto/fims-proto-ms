package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestApp_HandleRecordVoucherHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		constructor func() *RecordVoucherCmd
	}{
		{
			name: "normal_success",
			constructor: func() *RecordVoucherCmd {
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
			handler := NewRecordVoucherHandler(repoMock, &accServiceMock, &cntServiceMock)
			newUUID, err := handler.Handle(context.Background(), *cmd)

			assertions.NoError(err)
			vouchers := repoMock.vouchers

			d100, _ := decimal.NewFromString("100")

			v := vouchers[newUUID]

			assertions.Len(vouchers, 1)
			assertions.Len(v.LineItems(), 2)
			assertions.Equal(d100, v.Credit())
			assertions.Equal(d100, v.Debit())
			assertions.Equal("0000", v.Creator())
			assertions.True(accServiceMock.invoked)
		})
	}
}

func createVoucherCmd() *RecordVoucherCmd {
	lineItems := []LineItemCmd{
		{
			Summary:       "test_item1",
			AccountNumber: "1000",
			Debit:         "100",
			Credit:        "",
		},
		{
			Summary:       "test_item2",
			AccountNumber: "2000",
			Debit:         "",
			Credit:        "100",
		},
	}
	return &RecordVoucherCmd{
		Sob:                "test_sob",
		VoucherType:        "GENERAL_VOUCHER",
		AttachmentQuantity: 0,
		LineItems:          lineItems,
		Creator:            "0000",
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

type voucherRepoMock struct {
	vouchers map[uuid.UUID]domain.Voucher
}

func (r voucherRepoMock) AddVoucher(ctx context.Context, v *domain.Voucher) (uuid.UUID, error) {
	r.vouchers[v.Id()] = *v
	return v.Id(), nil
}

func (r voucherRepoMock) UpdateVoucher(ctx context.Context, voucherUUID uuid.UUID, updateFn func(v *domain.Voucher) (*domain.Voucher, error)) error {
	v, ok := r.vouchers[voucherUUID]
	if !ok {
		return errors.Errorf("voucher %s not exists", voucherUUID)
	}

	updatedVoucher, err := updateFn(&v)
	if err != nil {
		return errors.Wrapf(err, "voucher %s updated failed", voucherUUID)
	}
	r.vouchers[voucherUUID] = *updatedVoucher
	return nil
}

type accountServiceMock struct {
	invoked bool
}

func (s *accountServiceMock) ValidateExistence(ctx context.Context, sob string, accNumbers []string) error {
	s.invoked = true
	return nil
}

type counterServiceMock struct {
	invoked bool
}

func (c *counterServiceMock) GetNextIdentifier(ctx context.Context, bo ...string) (string, error) {
	c.invoked = true
	return "1", nil
}

func prepareBalancedItems() []*domain.LineItem {
	item1, _ := domain.NewLineItem(uuid.New(), "test", "1000", "100", "")
	item2, _ := domain.NewLineItem(uuid.New(), "test", "1001", "100", "")
	item3, _ := domain.NewLineItem(uuid.New(), "test", "2000", "", "150")
	item4, _ := domain.NewLineItem(uuid.New(), "test", "2001", "", "50")
	return []*domain.LineItem{
		item1,
		item2,
		item3,
		item4,
	}
}
