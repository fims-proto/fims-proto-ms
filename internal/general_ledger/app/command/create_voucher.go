package command

import (
	"context"
	"time"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/voucher"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"

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
	numberingService service.NumberingService
}

func NewCreateVoucherHandler(repo domain.Repository, numberingService service.NumberingService) CreateVoucherHandler {
	if repo == nil {
		panic("nil repo")
	}

	if numberingService == nil {
		panic("nil numbering service")
	}

	return CreateVoucherHandler{
		repo:             repo,
		numberingService: numberingService,
	}
}

func (h CreateVoucherHandler) Handle(ctx context.Context, cmd CreateVoucherCmd) error {
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		periodId, err := readPeriodIdAndCheck(txCtx, h.repo, h.numberingService, cmd.SobId, cmd.TransactionTime)
		if err != nil {
			return errors.Wrap(err, "failed to read or create period")
		}

		return h.createVoucher(txCtx, cmd, periodId)
	})
}

func (h CreateVoucherHandler) createVoucher(ctx context.Context, cmd CreateVoucherCmd, periodId uuid.UUID) error {
	// prepare line items
	lineItems, err := prepareLineItems(ctx, h.repo, cmd.SobId, cmd.LineItems)
	if err != nil {
		return errors.Wrap(err, "failed to prepare line items")
	}

	// get document number
	identifier, err := h.numberingService.GenerateIdentifier(ctx, periodId, cmd.VoucherType)
	if err != nil {
		return errors.Wrap(err, "unable to generate next number")
	}

	newVoucher, err := voucher.New(
		cmd.VoucherId,
		cmd.SobId,
		periodId,
		cmd.VoucherType,
		cmd.HeaderText,
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
