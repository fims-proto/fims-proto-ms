package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/voucher/domain/voucher"

	"github/fims-proto/fims-proto-ms/internal/voucher/domain"

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
	return h.repo.UpdateVoucher(
		ctx,
		cmd.VoucherId,
		func(j *voucher.Voucher) (*voucher.Voucher, error) {
			err := j.CancelReview(cmd.Reviewer)
			return j, err
		},
	)
}
