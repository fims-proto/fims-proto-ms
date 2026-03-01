package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/voucher"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"

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
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		return h.cancelAudit(txCtx, cmd)
	})
}

func (h CancelAuditVoucherHandler) cancelAudit(ctx context.Context, cmd CancelAuditVoucherCmd) error {
	return h.repo.UpdateVoucherHeader(
		ctx,
		cmd.VoucherId,
		func(v *voucher.Voucher) (*voucher.Voucher, error) {
			err := v.CancelAudit(cmd.Auditor)
			return v, err
		},
	)
}
