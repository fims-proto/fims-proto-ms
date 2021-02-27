package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/voucher"

	"github.com/google/uuid"
)

type AuditVoucherCmd struct {
	VoucherUUID uuid.UUID
	Auditor     string
}

type AuditVoucherHandler struct {
	repo voucher.Repository
}

func NewAuditVoucherHandler(repo voucher.Repository) AuditVoucherHandler {
	if repo == nil {
		panic("nil repo")
	}
	return AuditVoucherHandler{repo: repo}
}

func (h AuditVoucherHandler) Handle(ctx context.Context, cmd AuditVoucherCmd) error {
	return h.repo.UpdateVoucher(
		ctx,
		cmd.VoucherUUID,
		func(v *voucher.Voucher) (*voucher.Voucher, error) {
			err := v.Audit(cmd.Auditor)
			return v, err
		},
	)
}

func (h AuditVoucherHandler) HandleCancel(ctx context.Context, cmd AuditVoucherCmd) error {
	return h.repo.UpdateVoucher(
		ctx,
		cmd.VoucherUUID,
		func(v *voucher.Voucher) (*voucher.Voucher, error) {
			err := v.CancelAudit(cmd.Auditor)
			return v, err
		},
	)
}
