package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/voucher"
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
		constructor func(t *testing.T) *voucher.Voucher
		reviewer    string
	}{
		{
			name: "normal_success",
			constructor: func(t *testing.T) *voucher.Voucher {
				return createVoucherForReviewTest(t, "")
			},
			reviewer: "aud1_uuid",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			v := test.constructor(t)
			repoMock := newVoucherRepoMock()
			repoMock.vouchers = map[uuid.UUID]voucher.Voucher{
				v.UUID(): *v,
			}

			handler := NewReviewVoucherHandler(repoMock)

			err := handler.Handle(context.Background(), ReviewVoucherCmd{
				VoucherUUID: v.UUID(),
				Reviewer:    test.reviewer,
			})
			assert.NoError(t, err)

			assert.True(t, repoMock.vouchers[v.UUID()].IsReviewed())
		})
	}
}

func createVoucherForReviewTest(t *testing.T, reviewer string) *voucher.Voucher {
	v, err := voucher.NewVoucher(uuid.New(), "1", time.Now(), 0, prepareBalancedItems(), "")
	require.NoError(t, err)
	if reviewer != "" {
		err := v.Review(reviewer)
		require.NoError(t, err)
	}
	return v
}
