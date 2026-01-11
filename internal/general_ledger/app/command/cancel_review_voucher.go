package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/voucher"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"

	"github.com/google/uuid"
)

type CancelReviewVoucherCmd struct {
	VoucherId uuid.UUID
	Reviewer  uuid.UUID
}

type CancelReviewVoucherHandler struct {
	repo domain.Repository
}

func NewCancelReviewVoucherHandler(repo domain.Repository) CancelReviewVoucherHandler {
	if repo == nil {
		panic("nil repo")
	}
	return CancelReviewVoucherHandler{repo: repo}
}

func (h CancelReviewVoucherHandler) Handle(ctx context.Context, cmd CancelReviewVoucherCmd) error {
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		return h.cancelReview(txCtx, cmd)
	})
}

func (h CancelReviewVoucherHandler) cancelReview(ctx context.Context, cmd CancelReviewVoucherCmd) error {
	return h.repo.UpdateVoucher(
		ctx,
		cmd.VoucherId,
		func(v *voucher.Voucher) (*voucher.Voucher, error) {
			err := v.CancelReview(cmd.Reviewer)
			return v, err
		},
	)
}
