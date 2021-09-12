package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/common/log"
	"github/fims-proto/fims-proto-ms/internal/user/lib/authorization"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type RecordVoucherCmd struct {
	Sob                string
	VoucherType        string
	AttachmentQuantity uint
	LineItems          []LineItemCmd
	Creator            string
}

type RecordVoucherHandler struct {
	repo           domain.Repository
	accountService AccountService
	counterService CounterService
}

func NewRecordVoucherHandler(repo domain.Repository, accountService AccountService, counterService CounterService) RecordVoucherHandler {
	if repo == nil {
		panic("nil repo")
	}
	if accountService == nil {
		panic("nil account service")
	}
	if counterService == nil {
		panic("nil counter service")
	}
	return RecordVoucherHandler{
		repo:           repo,
		accountService: accountService,
		counterService: counterService,
	}
}

func (h RecordVoucherHandler) Handle(ctx context.Context, cmd RecordVoucherCmd) (createdId uuid.UUID, err error) {
	log.Info(ctx, "handle recording voucher")
	log.Debug(ctx, "handle recording voucher, cmd: %+v", cmd)
	defer func() {
		if err != nil {
			log.Err(ctx, err, "handle recording failed")
		}
	}()

	if err := authorization.VerifyAuth("current-user", cmd.Sob, "voucher", "create"); err != nil {
		return uuid.Nil, errors.Wrap(err, "failed verifing permission")
	}

	var accNumbers []string
	var lineItems []*domain.LineItem
	for _, item := range cmd.LineItems {
		lineItem, err := domain.NewLineItem(
			uuid.New(),
			item.Summary,
			item.AccountNumber,
			item.Debit,
			item.Credit,
		)
		if err != nil {
			return uuid.Nil, err
		}
		lineItems = append(lineItems, lineItem)
		accNumbers = append(accNumbers, item.AccountNumber)
	}

	if err := h.accountService.ValidateExistence(ctx, cmd.Sob, accNumbers); err != nil {
		return uuid.Nil, errors.Wrap(err, "unable to validate account numbers")
	}

	identifier, err := h.counterService.GetNextIdentifier(ctx, cmd.Sob, cmd.VoucherType)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "unable to generate next number")
	}

	newVoucher, err := domain.NewVoucher(
		uuid.New(),
		cmd.Sob,
		cmd.VoucherType,
		identifier,
		cmd.AttachmentQuantity,
		lineItems,
		cmd.Creator,
		"",    // no reviewer
		"",    // no auditor
		false, // not reviewed yet
		false, // not audited yet
		false, // not posted yet
	)
	if err != nil {
		return uuid.Nil, err
	}

	return h.repo.AddVoucher(ctx, newVoucher)
}
