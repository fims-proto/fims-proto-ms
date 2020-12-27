package voucher

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/lineitem"
	"testing"
	"time"
)

func TestVoucher_Review_Success(t *testing.T) {
	t.Parallel()
	voucher, err := NewVoucher("test_uuid", 1, time.Now(), 0, []lineitem.LineItem{}, "", "", false, "", false)
	require.NoError(t, err)

	err = voucher.Review("rev1_uuid")
	require.NoError(t, err)

	assert.Equal(t, "rev1_uuid", voucher.ReviewerUUID())
}

func TestVoucher_Review_RepeatReview_Error(t *testing.T) {
	t.Parallel()
	voucher, err := NewVoucher("test_uuid", 1, time.Now(), 0, []lineitem.LineItem{}, "", "rev1_uuid", true, "", false)
	require.NoError(t, err)

	err = voucher.Review("rev1_uuid")
	assert.Equal(t, ErrVoucherAlreadyReviewed, err)
}

func TestVoucher_CancelReview_Success(t *testing.T) {
	t.Parallel()
	voucher, err := NewVoucher("test_uuid", 1, time.Now(), 0, []lineitem.LineItem{}, "", "rev1_uuid", true, "", false)
	require.NoError(t, err)

	err = voucher.CancelReview("rev1_uuid")
	require.NoError(t, err)

	assert.Equal(t, "", voucher.ReviewerUUID())
	assert.False(t, voucher.IsReviewed())
}

func TestVoucher_CancelReview_NotReviewed_Error(t *testing.T) {
	t.Parallel()
	voucher, err := NewVoucher("test_uuid", 1, time.Now(), 0, []lineitem.LineItem{}, "", "", false, "", false)
	require.NoError(t, err)

	err = voucher.CancelReview("rev1_uuid")
	assert.Equal(t, ErrVoucherNotReviewed, err)
}

func TestVoucher_CancelReview_DifferentReviewer_Error(t *testing.T) {
	t.Parallel()
	voucher, err := NewVoucher("test_uuid", 1, time.Now(), 0, []lineitem.LineItem{}, "", "rev1_uuid", true, "", false)
	require.NoError(t, err)

	err = voucher.CancelReview("rev2_uuid")
	assert.Equal(t, ErrDifferentReviewerCancel, err)
}
