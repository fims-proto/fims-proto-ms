package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/common/log"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain"

	"github.com/google/uuid"
)

type AuditVoucherCmd struct {
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

func (h AuditVoucherHandler) Handle(ctx context.Context, cmd AuditVoucherCmd) (err error) {
	log.Info(ctx, "handle auditing voucher")
	log.Debug(ctx, "handle auditing voucher, cmd: %+v", cmd)
	defer func() {
		if err != nil {
			log.Err(ctx, err, "handle auditing failed")
		}
	}()

	return h.repo.UpdateVoucher(
		ctx,
		cmd.VoucherUUID,
		func(v *domain.Voucher) (*domain.Voucher, error) {
			log.Info(ctx, "auditing voucher")
			err := v.Audit(cmd.Auditor)
			return v, err
		},
	)
}

func (h AuditVoucherHandler) HandleCancel(ctx context.Context, cmd AuditVoucherCmd) (err error) {
	log.Info(ctx, "handle cancelling audit voucher")
	log.Debug(ctx, "handle cancelling audit voucher, cmd: %+v", cmd)
	defer func() {
		if err != nil {
			log.Err(ctx, err, "handle cancelling audit failed")
		}
	}()

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
