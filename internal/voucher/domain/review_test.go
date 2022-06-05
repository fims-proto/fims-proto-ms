package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDomain_VoucherReview(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		doReview    bool // true - review, false - cancel review
		constructor func(t *testing.T) *Voucher
		reviewer    string
		verify      func(t *testing.T, v Voucher, err error)
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
				assert.Equal(t, "aud1_uuid", v.Reviewer())
			},
		},
		{
			"review_sameAsCreator_error",
			true,
			func(t *testing.T) *Voucher {
				v := createVoucherForReviewTest(t, "")
				v.creator = "aud1_uuid"
				return v
			},
			"aud1_uuid",
			func(t *testing.T, v Voucher, err error) {
				assert.Error(t, err)
			},
		},
		{
			"review_sameAsAuditor_error",
			true,
			func(t *testing.T) *Voucher {
				v := createVoucherForReviewTest(t, "")
				v.auditor = "aud1_uuid"
				return v
			},
			"aud1_uuid",
			func(t *testing.T, v Voucher, err error) {
				assert.Error(t, err)
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
				assert.Equal(t, "", v.Reviewer())
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
				err = voucher.Review(test.reviewer)
			} else {
				err = voucher.CancelReview(test.reviewer)
			}
			test.verify(t, *voucher, err)
		})
	}
}

func createVoucherForReviewTest(t *testing.T, reviewer string) *Voucher {
	voucher, err := NewVoucher(uuid.New(), uuid.New(), "GENERAL_VOUCHER", "1", 0, prepareBalancedItems(), "creator", "", "", false, false, false, time.Now())
	require.NoError(t, err)
	if reviewer != "" {
		err := voucher.Review(reviewer)
		require.NoError(t, err)
	}
	return voucher
}
