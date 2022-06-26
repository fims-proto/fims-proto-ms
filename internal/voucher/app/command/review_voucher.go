package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/log"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain"

	"github.com/google/uuid"
)

type ReviewVoucherCmd struct {
	VoucherUUID uuid.UUID
	Reviewer    uuid.UUID
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

func (h ReviewVoucherHandler) Handle(ctx context.Context, cmd ReviewVoucherCmd) (err error) {
	log.Info(ctx, "handle reviewing voucher")
	log.Debug(ctx, "handle reviewing voucher, cmd: %+v", cmd)
	defer func() {
		if err != nil {
			log.Err(ctx, err, "handle reviewing failed")
		}
	}()

	return h.repo.UpdateVoucher(
		ctx,
		cmd.VoucherUUID,
		func(v *domain.Voucher) (*domain.Voucher, error) {
			log.Info(ctx, "reviewing voucher")
			err := v.Review(cmd.Reviewer)
			return v, err
		},
	)
}

func (h ReviewVoucherHandler) HandleCancel(ctx context.Context, cmd ReviewVoucherCmd) (err error) {
	log.Info(ctx, "handle cancelling review voucher")
	log.Debug(ctx, "handle cancelling review voucher, cmd: %+v", cmd)
	defer func() {
		if err != nil {
			log.Err(ctx, err, "handle cancelling review failed")
		}
	}()

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
