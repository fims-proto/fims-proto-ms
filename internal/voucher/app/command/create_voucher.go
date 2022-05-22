package command

import (
	"context"
	"time"

	"github/fims-proto/fims-proto-ms/internal/common/log"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type CreateVoucherCmd struct {
	SobId              uuid.UUID
	VoucherType        string
	AttachmentQuantity uint
	LineItems          []LineItemCmd
	Creator            string
	TransactionTime    time.Time
}

type CreateVoucherHandler struct {
	repo           domain.Repository
	accountService AccountService
	counterService CounterService
}

func NewCreateVoucherHandler(repo domain.Repository, accountService AccountService, counterService CounterService) CreateVoucherHandler {
	if repo == nil {
		panic("nil repo")
	}
	if accountService == nil {
		panic("nil account service")
	}
	if counterService == nil {
		panic("nil counter service")
	}
	return CreateVoucherHandler{
		repo:           repo,
		accountService: accountService,
		counterService: counterService,
	}
}

func (h CreateVoucherHandler) Handle(ctx context.Context, cmd CreateVoucherCmd) (createdId uuid.UUID, err error) {
	log.Info(ctx, "handle creating voucher")
	log.Debug(ctx, "handle creating voucher, cmd: %+v", cmd)
	defer func() {
		if err != nil {
			log.Err(ctx, err, "handle creating failed")
		}
	}()

	// validate account numbers
	var accountNumbers []string
	for _, item := range cmd.LineItems {
		accountNumbers = append(accountNumbers, item.AccountNumber)
	}
	accountIds, err := h.accountService.ValidateExistenceAndGetId(ctx, cmd.SobId, accountNumbers)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "unable to validate account numbers")
	}

	// validate line items
	var lineItems []*domain.LineItem
	for _, item := range cmd.LineItems {
		accountId, ok := accountIds[item.AccountNumber]
		if !ok {
			return uuid.Nil, errors.Wrapf(err, "unable to find account id by number %s", item.AccountNumber)
		}
		lineItem, err := domain.NewLineItem(
			uuid.New(),
			accountId,
			item.Summary,
			item.Debit,
			item.Credit,
		)
		if err != nil {
			return uuid.Nil, err
		}
		lineItems = append(lineItems, lineItem)
	}

	// get voucher number
	identifier, err := h.counterService.GetNextIdentifier(ctx, cmd.SobId.String(), cmd.VoucherType)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "unable to generate next number")
	}

	newVoucher, err := domain.NewVoucher(
		uuid.New(),
		cmd.SobId,
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
		cmd.TransactionTime,
	)
	if err != nil {
		return uuid.Nil, err
	}

	return h.repo.CreateVoucher(ctx, newVoucher)
}
