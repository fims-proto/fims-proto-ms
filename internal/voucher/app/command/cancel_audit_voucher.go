package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/log"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain"

	"github.com/google/uuid"
)

type CancelAuditVoucherCmd struct {
	VoucherUUID uuid.UUID
	Auditor     uuid.UUID
}

type CancelAuditVoucherHandler struct {
	repo domain.Repository
}

func NewCancelAuditVoucherHandler(repo domain.Repository) CancelAuditVoucherHandler {
	if repo == nil {
		panic("nil repo")
	}
	return CancelAuditVoucherHandler{repo: repo}
}

func (h CancelAuditVoucherHandler) Handle(ctx context.Context, cmd CancelAuditVoucherCmd) error {
	return h.repo.UpdateVoucher(
		ctx,
		cmd.VoucherUUID,
		func(v *domain.Voucher) (*domain.Voucher, error) {
			log.Info(ctx, "cancelling audit voucher")
			err := v.CancelAudit(cmd.Auditor)
			return v, err
		},
	)
}
