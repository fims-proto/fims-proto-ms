package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDomain_VoucherAudit(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		doAudit     bool // true - audit, false - cancel audit
		constructor func(t *testing.T) *Voucher
		auditor     string
		verify      func(t *testing.T, v Voucher, err error)
	}{
		{
			"audit_success",
			true,
			func(t *testing.T) *Voucher {
				return createVoucherForAuditTest(t, "")
			},
			"aud1_uuid",
			func(t *testing.T, v Voucher, err error) {
				require.NoError(t, err)
				assert.Equal(t, "aud1_uuid", v.Auditor())
			},
		},
		{
			"audit_repeat_audit_error",
			true,
			func(t *testing.T) *Voucher {
				return createVoucherForAuditTest(t, "aud1_uuid")
			},
			"aud1_uuid",
			func(t *testing.T, v Voucher, err error) {
				assert.Equal(t, ErrVoucherAlreadyAudited, err)
			},
		},
		{
			"audit_cancel_audit_success",
			false,
			func(t *testing.T) *Voucher {
				return createVoucherForAuditTest(t, "aud1_uuid")
			},
			"aud1_uuid",
			func(t *testing.T, v Voucher, err error) {
				require.NoError(t, err)
				assert.Equal(t, "", v.Auditor())
				assert.False(t, v.IsAudited())
			},
		},
		{
			"audit_cancel_audit_not_audited_error",
			false,
			func(t *testing.T) *Voucher {
				return createVoucherForAuditTest(t, "")
			},
			"aud1_uuid",
			func(t *testing.T, v Voucher, err error) {
				assert.Equal(t, ErrVoucherNotAudited, err)
			},
		},
		{
			"audit_cancel_audit_different_auditor_error",
			false,
			func(t *testing.T) *Voucher {
				return createVoucherForAuditTest(t, "aud1_uuid")
			},
			"aud2_uuid",
			func(t *testing.T, v Voucher, err error) {
				assert.Equal(t, ErrDifferentAuditorCancel, err)
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			voucher := test.constructor(t)
			var err error
			if test.doAudit {
				err = voucher.Audit(test.auditor)
			} else {
				err = voucher.CancelAudit(test.auditor)
			}
			test.verify(t, *voucher, err)
		})
	}
}

func createVoucherForAuditTest(t *testing.T, auditor string) *Voucher {
	voucher, err := NewVoucher("test_sob", uuid.New(), GeneralVoucher, "1", time.Now(), 0, prepareBalancedItems(), "")
	require.NoError(t, err)
	if auditor != "" {
		err := voucher.Audit(auditor)
		require.NoError(t, err)
	}
	return voucher
}
