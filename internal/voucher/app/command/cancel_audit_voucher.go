package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/voucher/domain/voucher"

	"github/fims-proto/fims-proto-ms/internal/voucher/domain"

	"github.com/google/uuid"
)

type CancelAuditVoucherCmd struct {
	VoucherId uuid.UUID
	Auditor   uuid.UUID
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
		cmd.VoucherId,
		func(j *voucher.Voucher) (*voucher.Voucher, error) {
			err := j.CancelAudit(cmd.Auditor)
			return j, err
		},
	)
}
