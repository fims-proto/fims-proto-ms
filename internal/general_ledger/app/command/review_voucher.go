package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/voucher"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"

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
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		return h.review(txCtx, cmd)
	})
}

func (h ReviewVoucherHandler) review(ctx context.Context, cmd ReviewVoucherCmd) error {
	return h.repo.UpdateVoucherHeader(
		ctx,
		cmd.VoucherId,
		func(v *voucher.Voucher) (*voucher.Voucher, error) {
			err := v.Review(cmd.Reviewer)
			return v, err
		},
	)
}
