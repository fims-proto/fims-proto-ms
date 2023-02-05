package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/voucher/domain/voucher"

	"github/fims-proto/fims-proto-ms/internal/voucher/domain"

	"github.com/google/uuid"
)

type ReviewVoucherCmd struct {
	VoucherId uuid.UUID
	Reviewer  uuid.UUID
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
		cmd.VoucherId,
		func(j *voucher.Voucher) (*voucher.Voucher, error) {
			err := j.Review(cmd.Reviewer)
			return j, err
		},
	)
}
