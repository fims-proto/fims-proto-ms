package voucher

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDomain_VoucherUpdate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		constructor func(t *testing.T) *Voucher
		verify      func(t *testing.T, v Voucher, err error)
	}{
		{
			"update_success",
			func(t *testing.T) *Voucher {
				return createVoucherForUpdateTest(t, false, false)
			},
			func(t *testing.T, v Voucher, err error) {
				require.NoError(t, err)
			},
		},
		{
			"update_reviewed",
			func(t *testing.T) *Voucher {
				return createVoucherForUpdateTest(t, true, false)
			},
			func(t *testing.T, v Voucher, err error) {
				assert.EqualError(t, err, ErrVoucherAlreadyReviewed.Error())
			},
		},
		{
			"update_audited",
			func(t *testing.T) *Voucher {
				return createVoucherForUpdateTest(t, false, true)
			},
			func(t *testing.T, v Voucher, err error) {
				assert.EqualError(t, err, ErrVoucherAlreadyAudited.Error())
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			voucher := test.constructor(t)

			err := voucher.Update(prepareBalancedItems())

			test.verify(t, *voucher, err)
		})
	}
}

func createVoucherForUpdateTest(t *testing.T, reviewed bool, audited bool) *Voucher {
	v, err := NewVoucher("test_uuid", "1", time.Now(), 0, prepareBalancedItems(), "")
	require.NoError(t, err)
	if reviewed {
		require.NoError(t, v.Review("r"))
	}
	if audited {
		require.NoError(t, v.Audit("a"))
	}
	return v
}
