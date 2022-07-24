package command

import (
	"context"
	"testing"
	"time"

	"github/fims-proto/fims-proto-ms/internal/voucher/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApp_HandleReviewVoucher(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		cancel      bool
		constructor func(t *testing.T) *domain.Voucher
		reviewer    uuid.UUID
	}{
		{
			name:   "review_success",
			cancel: false,
			constructor: func(t *testing.T) *domain.Voucher {
				return createVoucherForReviewTest(t, uuid.Nil)
			},
			reviewer: userA,
		},
		{
			name:   "cancel_review_success",
			cancel: true,
			constructor: func(t *testing.T) *domain.Voucher {
				return createVoucherForReviewTest(t, userA)
			},
			reviewer: userA,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			v := test.constructor(t)
			repoMock := newVoucherRepoMock()
			repoMock.vouchers = map[uuid.UUID]domain.Voucher{
				v.Id(): *v,
			}

			handler := NewReviewVoucherHandler(repoMock)
			cancelHandler := NewCancelReviewVoucherHandler(repoMock)

			if test.cancel {
				err := cancelHandler.Handle(context.Background(), CancelReviewVoucherCmd{
					VoucherUUID: v.Id(),
					Reviewer:    test.reviewer,
				})
				assert.NoError(t, err)
				assert.False(t, repoMock.vouchers[v.Id()].IsReviewed())
			} else {
				err := handler.Handle(context.Background(), ReviewVoucherCmd{
					VoucherUUID: v.Id(),
					Reviewer:    test.reviewer,
				})
				assert.NoError(t, err)
				assert.True(t, repoMock.vouchers[v.Id()].IsReviewed())
			}
		})
	}
}

func createVoucherForReviewTest(t *testing.T, reviewer uuid.UUID) *domain.Voucher {
	v, err := domain.NewVoucher(uuid.New(), uuid.New(), uuid.New(), "general_voucher", "1", 0, prepareBalancedItems(), uuid.New(), uuid.Nil, uuid.Nil, false, false, false, time.Now())
	require.NoError(t, err)
	if reviewer != uuid.Nil {
		err := v.Review(reviewer)
		require.NoError(t, err)
	}
	return v
}
