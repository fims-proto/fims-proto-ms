package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/voucher"
)

type AuditVoucherCmd struct {
	VoucherUUID string
	AuditorUUID string
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
			err := v.Audit(cmd.AuditorUUID)
			return v, err
		},
	)
}
