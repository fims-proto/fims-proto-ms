package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/lineitem"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/voucher"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApp_HandleReviewVoucher(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		constructor  func(t *testing.T) *voucher.Voucher
		reviewerUUID string
	}{
		{
			name: "normal_success",
			constructor: func(t *testing.T) *voucher.Voucher {
				return createVoucherForReviewTest(t, "")
			},
			reviewerUUID: "aud1_uuid",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			v := test.constructor(t)
			repoMock := newVoucherRepoMock()
			repoMock.vouchers = map[string]voucher.Voucher{
				v.UUID(): *v,
			}

			handler := NewReviewVoucherHandler(repoMock)

			err := handler.Handle(context.Background(), ReviewVoucherCmd{
				VoucherUUID:  v.UUID(),
				ReviewerUUID: test.reviewerUUID,
			})
			assert.NoError(t, err)

			assert.True(t, repoMock.vouchers[v.UUID()].IsReviewed())
		})
	}
}

func createVoucherForReviewTest(t *testing.T, reviewerUUID string) *voucher.Voucher {
	v, err := voucher.NewVoucher("test_uuid", "1", time.Now(), 0, []lineitem.LineItem{}, "")
	require.NoError(t, err)
	if reviewerUUID != "" {
		err := v.Review(reviewerUUID)
		require.NoError(t, err)
	}
	return v
}
