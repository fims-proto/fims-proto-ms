package voucher

import (
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/lineitem"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDomain_VoucherReview(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		doReview     bool // true - review, false - cancel review
		constructor  func(t *testing.T) *Voucher
		reviewerUUID string
		verify       func(t *testing.T, v Voucher, err error)
	}{
		{
			"review_success",
			true,
			func(t *testing.T) *Voucher {
				return createVoucherForReviewTest(t, "")
			},
			"aud1_uuid",
			func(t *testing.T, v Voucher, err error) {
				require.NoError(t, err)
				assert.Equal(t, "aud1_uuid", v.ReviewerUUID())
			},
		},
		{
			"review_repeat_review_error",
			true,
			func(t *testing.T) *Voucher {
				return createVoucherForReviewTest(t, "aud1_uuid")
			},
			"aud1_uuid",
			func(t *testing.T, v Voucher, err error) {
				assert.Equal(t, ErrVoucherAlreadyReviewed, err)
			},
		},
		{
			"review_cancel_review_success",
			false,
			func(t *testing.T) *Voucher {
				return createVoucherForReviewTest(t, "aud1_uuid")
			},
			"aud1_uuid",
			func(t *testing.T, v Voucher, err error) {
				require.NoError(t, err)
				assert.Equal(t, "", v.ReviewerUUID())
				assert.False(t, v.IsReviewed())
			},
		},
		{
			"review_cancel_review_not_reviewed_error",
			false,
			func(t *testing.T) *Voucher {
				return createVoucherForReviewTest(t, "")
			},
			"aud1_uuid",
			func(t *testing.T, v Voucher, err error) {
				assert.Equal(t, ErrVoucherNotReviewed, err)
			},
		},
		{
			"review_cancel_review_different_reviewer_error",
			false,
			func(t *testing.T) *Voucher {
				return createVoucherForReviewTest(t, "aud1_uuid")
			},
			"aud2_uuid",
			func(t *testing.T, v Voucher, err error) {
				assert.Equal(t, ErrDifferentReviewerCancel, err)
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			voucher := test.constructor(t)
			var err error
			if test.doReview {
				err = voucher.Review(test.reviewerUUID)
			} else {
				err = voucher.CancelReview(test.reviewerUUID)
			}
			test.verify(t, *voucher, err)
		})
	}
}

func createVoucherForReviewTest(t *testing.T, reviewerUUID string) *Voucher {
	voucher, err := NewVoucher("test_uuid", "1", time.Now(), 0, []lineitem.LineItem{}, "")
	require.NoError(t, err)
	if reviewerUUID != "" {
		err := voucher.Review(reviewerUUID)
		require.NoError(t, err)
	}
	return voucher
}
