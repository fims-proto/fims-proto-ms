package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/voucher"
)

type ReviewVoucherCmd struct {
	VoucherUUID  string
	ReviewerUUID string
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
			err := v.Review(cmd.ReviewerUUID)
			return v, err
		},
	)
}
