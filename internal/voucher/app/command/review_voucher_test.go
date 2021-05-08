package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain"
	"testing"
	"time"

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
		reviewer    string
	}{
		{
			name:   "review_success",
			cancel: false,
			constructor: func(t *testing.T) *domain.Voucher {
				return createVoucherForReviewTest(t, "")
			},
			reviewer: "reviewer1",
		},
		{
			name:   "cancel_review_success",
			cancel: true,
			constructor: func(t *testing.T) *domain.Voucher {
				return createVoucherForReviewTest(t, "reviewer1")
			},
			reviewer: "reviewer1",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			v := test.constructor(t)
			repoMock := newVoucherRepoMock()
			repoMock.vouchers = map[uuid.UUID]domain.Voucher{
				v.UUID(): *v,
			}

			handler := NewReviewVoucherHandler(repoMock)

			if test.cancel {
				err := handler.HandleCancel(context.Background(), ReviewVoucherCmd{
					VoucherUUID: v.UUID(),
					Reviewer:    test.reviewer,
				})
				assert.NoError(t, err)
				assert.False(t, repoMock.vouchers[v.UUID()].IsReviewed())
			} else {
				err := handler.Handle(context.Background(), ReviewVoucherCmd{
					VoucherUUID: v.UUID(),
					Reviewer:    test.reviewer,
				})
				assert.NoError(t, err)
				assert.True(t, repoMock.vouchers[v.UUID()].IsReviewed())
			}
		})
	}
}

func createVoucherForReviewTest(t *testing.T, reviewer string) *domain.Voucher {
	v, err := domain.NewVoucher(uuid.New(), domain.GeneralVoucher, "1", time.Now(), 0, prepareBalancedItems(), "")
	require.NoError(t, err)
	if reviewer != "" {
		err := v.Review(reviewer)
		require.NoError(t, err)
	}
	return v
}
