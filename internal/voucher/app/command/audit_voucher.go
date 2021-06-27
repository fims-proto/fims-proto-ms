package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain"

	"github.com/google/uuid"
)

type AuditVoucherCmd struct {
	Sob         string
	VoucherUUID uuid.UUID
	Auditor     string
}

type AuditVoucherHandler struct {
	repo domain.Repository
}

func NewAuditVoucherHandler(repo domain.Repository) AuditVoucherHandler {
	if repo == nil {
		panic("nil repo")
	}
	return AuditVoucherHandler{repo: repo}
}

func (h AuditVoucherHandler) Handle(ctx context.Context, cmd AuditVoucherCmd) error {
	return h.repo.UpdateVoucher(
		ctx,
		cmd.Sob,
		cmd.VoucherUUID,
		func(v *domain.Voucher) (*domain.Voucher, error) {
			err := v.Audit(cmd.Auditor)
			return v, err
		},
	)
}

func (h AuditVoucherHandler) HandleCancel(ctx context.Context, cmd AuditVoucherCmd) error {
	return h.repo.UpdateVoucher(
		ctx,
		cmd.Sob,
		cmd.VoucherUUID,
		func(v *domain.Voucher) (*domain.Voucher, error) {
			err := v.CancelAudit(cmd.Auditor)
			return v, err
		},
	)
}
