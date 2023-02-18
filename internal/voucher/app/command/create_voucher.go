package command

import (
	"context"
	"time"

	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"

	"github/fims-proto/fims-proto-ms/internal/voucher/domain/line_item"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/voucher"

	"github/fims-proto/fims-proto-ms/internal/voucher/app/service"

	"github/fims-proto/fims-proto-ms/internal/voucher/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type CreateVoucherCmd struct {
	VoucherId          uuid.UUID
	SobId              uuid.UUID
	HeaderText         string
	VoucherType        string
	AttachmentQuantity int
	LineItems          []LineItemCmd
	Creator            uuid.UUID
	TransactionTime    time.Time
}

type CreateVoucherHandler struct {
	repo             domain.Repository
	accountService   service.AccountService
	numberingService service.NumberingService
}

func NewCreateVoucherHandler(
	repo domain.Repository,
	accountService service.AccountService,
	numberingService service.NumberingService,
) CreateVoucherHandler {
	if repo == nil {
		panic("nil repo")
	}

	if accountService == nil {
		panic("nil account service")
	}

	if numberingService == nil {
		panic("nil numbering service")
	}

	return CreateVoucherHandler{
		repo:             repo,
		accountService:   accountService,
		numberingService: numberingService,
	}
}

func (h CreateVoucherHandler) Handle(ctx context.Context, cmd CreateVoucherCmd) error {
	// read period by transaction time
	period, err := h.accountService.ReadOrCreatePeriodByTime(ctx, cmd.SobId, cmd.TransactionTime)
	if err != nil {
		return errors.Wrap(err, "failed to read period by transaction time")
	}

	if period.IsClosed {
		return commonErrors.NewSlugError("voucher-periodClosed")
	}

	// validate account numbers
	var accountNumbers []string
	for _, item := range cmd.LineItems {
		accountNumbers = append(accountNumbers, item.AccountNumber)
	}
	accountIds, err := h.accountService.ValidateExistenceAndGetId(ctx, cmd.SobId, accountNumbers)
	if err != nil {
		return errors.Wrap(err, "unable to validate account numbers")
	}

	// validate line items
	var lineItems []line_item.LineItem
	for _, item := range cmd.LineItems {
		accountId, ok := accountIds[item.AccountNumber]
		if !ok {
			return errors.Errorf("should not happen, unable to find account id by number %s", item.AccountNumber)
		}
		lineItem, err := line_item.New(
			uuid.New(),
			accountId,
			item.Text,
			item.Debit,
			item.Credit,
		)
		if err != nil {
			return err
		}
		lineItems = append(lineItems, *lineItem)
	}

	// get document number
	identifier, err := h.numberingService.GenerateIdentifier(ctx, period.Id, cmd.VoucherType)
	if err != nil {
		return errors.Wrap(err, "unable to generate next number")
	}

	newVoucher, err := voucher.New(
		cmd.SobId,
		cmd.VoucherId,
		period.Id,
		cmd.HeaderText,
		cmd.VoucherType,
		identifier,
		cmd.AttachmentQuantity,
		cmd.Creator,
		uuid.Nil,
		uuid.Nil,
		uuid.Nil,
		false,
		false,
		false,
		cmd.TransactionTime,
		lineItems,
	)
	if err != nil {
		return err
	}

	return h.repo.CreateVoucher(ctx, newVoucher)
}
