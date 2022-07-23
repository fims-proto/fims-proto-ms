package command

import (
	"context"
	"time"

	"github/fims-proto/fims-proto-ms/internal/voucher/app/service"

	"github/fims-proto/fims-proto-ms/internal/voucher/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type CreateVoucherCmd struct {
	SobId              uuid.UUID
	VoucherType        string
	AttachmentQuantity uint
	LineItems          []LineItemCmd
	Creator            uuid.UUID
	TransactionTime    time.Time
}

type CreateVoucherHandler struct {
	repo             domain.Repository
	accountService   service.AccountService
	numberingService service.NumberingService
	ledgerService    service.LedgerService
}

func NewCreateVoucherHandler(repo domain.Repository, accountService service.AccountService, numberingService service.NumberingService, ledgerService service.LedgerService) CreateVoucherHandler {
	if repo == nil {
		panic("nil repo")
	}
	if accountService == nil {
		panic("nil account service")
	}
	if numberingService == nil {
		panic("nil numbering service")
	}
	if ledgerService == nil {
		panic("nil ledger service")
	}
	return CreateVoucherHandler{
		repo:             repo,
		accountService:   accountService,
		numberingService: numberingService,
		ledgerService:    ledgerService,
	}
}

func (h CreateVoucherHandler) Handle(ctx context.Context, cmd CreateVoucherCmd) (uuid.UUID, error) {
	// read period by transaction time
	period, err := h.ledgerService.ReadPeriodByTime(ctx, cmd.SobId, cmd.TransactionTime)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "failed to read period by transaction time")
	}

	if period.IsClosed {
		return uuid.Nil, errors.New("period is closed")
	}

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
	identifier, err := h.numberingService.GenerateIdentifier(ctx, period.Id, cmd.VoucherType)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "unable to generate next number")
	}

	newVoucher, err := domain.NewVoucher(
		uuid.New(),
		cmd.SobId,
		period.Id,
		cmd.VoucherType,
		identifier,
		cmd.AttachmentQuantity,
		lineItems,
		cmd.Creator,
		uuid.Nil,
		uuid.Nil,
		false,
		false,
		false,
		cmd.TransactionTime,
	)
	if err != nil {
		return uuid.Nil, err
	}

	return h.repo.CreateVoucher(ctx, newVoucher)
}
