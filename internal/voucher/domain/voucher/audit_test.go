package voucher

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/lineitem"
	"testing"
	"time"
)

func TestVoucher_Audit_Success(t *testing.T) {
	t.Parallel()
	voucher, err := NewVoucher("test_uuid", 1, time.Now(), 0, []lineitem.LineItem{}, "", "", false, "", false)
	require.NoError(t, err)

	err = voucher.Audit("aud1_uuid")
	require.NoError(t, err)

	assert.Equal(t, "aud1_uuid", voucher.AuditorUUID())
}

func TestVoucher_Audit_RepeatAudit_Error(t *testing.T) {
	t.Parallel()
	voucher, err := NewVoucher("test_uuid", 1, time.Now(), 0, []lineitem.LineItem{}, "", "", false, "aud1_uuid", true)
	require.NoError(t, err)

	err = voucher.Audit("aud1_uuid")
	assert.Equal(t, ErrVoucherAlreadyAudited, err)
}

func TestVoucher_CancelAudit_Success(t *testing.T) {
	t.Parallel()
	voucher, err := NewVoucher("test_uuid", 1, time.Now(), 0, []lineitem.LineItem{}, "", "", false, "aud1_uuid", true)
	require.NoError(t, err)

	err = voucher.CancelAudit("aud1_uuid")
	require.NoError(t, err)

	assert.Equal(t, "", voucher.AuditorUUID())
	assert.False(t, voucher.IsAudited())
}

func TestVoucher_CancelAudit_NotAudited_Error(t *testing.T) {
	t.Parallel()
	voucher, err := NewVoucher("test_uuid", 1, time.Now(), 0, []lineitem.LineItem{}, "", "", false, "", false)
	require.NoError(t, err)

	err = voucher.CancelAudit("aud1_uuid")
	assert.Equal(t, ErrVoucherNotAudited, err)
}

func TestVoucher_CancelAudit_DifferentAuditor_Error(t *testing.T) {
	t.Parallel()
	voucher, err := NewVoucher("test_uuid", 1, time.Now(), 0, []lineitem.LineItem{}, "", "", false, "aud1_uuid", true)
	require.NoError(t, err)

	err = voucher.CancelAudit("aud2_uuid")
	assert.Equal(t, ErrDifferentAuditorCancel, err)
}
