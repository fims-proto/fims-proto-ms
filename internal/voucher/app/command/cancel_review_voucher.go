package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/log"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain"

	"github.com/google/uuid"
)

type CancelReviewVoucherCmd struct {
	VoucherUUID uuid.UUID
	Reviewer    uuid.UUID
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
		cmd.VoucherUUID,
		func(v *domain.Voucher) (*domain.Voucher, error) {
			log.Info(ctx, "cancelling review voucher")
			err := v.CancelReview(cmd.Reviewer)
			return v, err
		},
	)
}
