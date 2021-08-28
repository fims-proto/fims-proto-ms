package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain"

	"github.com/google/uuid"
)

type ReviewVoucherCmd struct {
	VoucherUUID uuid.UUID
	Reviewer    string
}

type ReviewVoucherHandler struct {
	repo domain.Repository
}

func NewReviewVoucherHandler(repo domain.Repository) ReviewVoucherHandler {
	if repo == nil {
		panic("nil repo")
	}
	return ReviewVoucherHandler{repo: repo}
}

func (h ReviewVoucherHandler) Handle(ctx context.Context, cmd ReviewVoucherCmd) error {
	return h.repo.UpdateVoucher(
		ctx,
		cmd.VoucherUUID,
		func(v *domain.Voucher) (*domain.Voucher, error) {
			err := v.Review(cmd.Reviewer)
			return v, err
		},
	)
}

func (h ReviewVoucherHandler) HandleCancel(ctx context.Context, cmd ReviewVoucherCmd) error {
	return h.repo.UpdateVoucher(
		ctx,
		cmd.VoucherUUID,
		func(v *domain.Voucher) (*domain.Voucher, error) {
			err := v.CancelReview(cmd.Reviewer)
			return v, err
		},
	)
}
