package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/voucher"

	"github.com/google/uuid"
)

type AuditVoucherCmd struct {
	VoucherId uuid.UUID
	Auditor   uuid.UUID
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
		cmd.VoucherId,
		func(v *voucher.Voucher) (*voucher.Voucher, error) {
			err := v.Audit(cmd.Auditor)
			return v, err
		},
	)
}
