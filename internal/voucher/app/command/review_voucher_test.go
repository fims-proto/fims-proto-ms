package command

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/lineitem"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/voucher"
	"testing"
	"time"
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
			deps := newReviewMockDeps()
			deps.repository.vouchers = map[string]voucher.Voucher{
				v.UUID(): *v,
			}

			err := deps.handler.Handle(context.Background(), ReviewVoucherCmd{
				VoucherUUID:  v.UUID(),
				ReviewerUUID: test.reviewerUUID,
			})
			assert.NoError(t, err)

			assert.True(t, deps.repository.vouchers[v.UUID()].IsReviewed())
		})
	}
}

func createVoucherForReviewTest(t *testing.T, reviewerUUID string) *voucher.Voucher {
	v, err := voucher.NewVoucher("test_uuid", 1, time.Now(), 0, []lineitem.LineItem{}, "")
	require.NoError(t, err)
	if reviewerUUID != "" {
		err := v.Review(reviewerUUID)
		require.NoError(t, err)
	}
	return v
}

type reviewMockDeps struct {
	repository *reviewRepoMock
	handler    ReviewVoucherHandler
}

func newReviewMockDeps() reviewMockDeps {
	repository := &reviewRepoMock{}
	return reviewMockDeps{
		repository: repository,
		handler:    ReviewVoucherHandler{repository},
	}
}

type reviewRepoMock struct {
	vouchers map[string]voucher.Voucher
}

func (r *reviewRepoMock) AddVoucher(ctx context.Context, v *voucher.Voucher) error {
	panic("implement me")
}

func (r *reviewRepoMock) UpdateVoucher(
	ctx context.Context,
	voucherUUID string,
	updateFn func(v *voucher.Voucher) (*voucher.Voucher, error),
) error {
	v, ok := r.vouchers[voucherUUID]
	if !ok {
		return voucher.NotFoundError{VoucherUUID: voucherUUID}
	}

	updatedVoucher, err := updateFn(&v)
	if err != nil {
		return err
	}

	r.vouchers[voucherUUID] = *updatedVoucher
	return nil
}
