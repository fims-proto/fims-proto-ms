package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/voucher"

	"github.com/google/uuid"
)

type ReviewVoucherCmd struct {
	VoucherUUID uuid.UUID
	Reviewer    string
}

type ReviewVoucherHandler struct {
	repo voucher.Repository
}

func NewReviewVoucherHandler(repo voucher.Repository) ReviewVoucherHandler {
	if repo == nil {
		panic("nil repo")
	}
	return ReviewVoucherHandler{repo: repo}
}

func (h ReviewVoucherHandler) Handle(ctx context.Context, cmd ReviewVoucherCmd) error {
	return h.repo.UpdateVoucher(
		ctx,
		cmd.VoucherUUID,
		func(v *voucher.Voucher) (*voucher.Voucher, error) {
			err := v.Review(cmd.Reviewer)
			return v, err
		},
	)
}

func (h ReviewVoucherHandler) HandleCancel(ctx context.Context, cmd ReviewVoucherCmd) error {
	return h.repo.UpdateVoucher(
		ctx,
		cmd.VoucherUUID,
		func(v *voucher.Voucher) (*voucher.Voucher, error) {
			err := v.CancelReview(cmd.Reviewer)
			return v, err
		},
	)
}
