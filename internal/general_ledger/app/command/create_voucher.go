package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/transaction_date"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/period"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/voucher"

	"github.com/google/uuid"
)

type CreateVoucherCmd struct {
	VoucherId          uuid.UUID
	SobId              uuid.UUID
	HeaderText         string
	VoucherType        string
	AttachmentQuantity int
	LineItems          []LineItemCmd
	Creator            uuid.UUID
	TransactionDate    transaction_date.TransactionDate
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
		p, err := readPeriodIdAndCheck(txCtx, h.repo, h.numberingService, cmd.SobId, cmd.TransactionDate)
		if err != nil {
			return fmt.Errorf("failed to read or create period: %w", err)
		}

		return h.createVoucher(txCtx, cmd, p)
	})
}

func (h CreateVoucherHandler) createVoucher(ctx context.Context, cmd CreateVoucherCmd, p *period.Period) error {
	// prepare line items
	lineItems, err := prepareLineItems(ctx, h.repo, cmd.SobId, cmd.LineItems)
	if err != nil {
		return fmt.Errorf("failed to prepare line items: %w", err)
	}

	// get document number
	identifier, err := h.numberingService.GenerateIdentifier(ctx, p.Id(), cmd.VoucherType)
	if err != nil {
		return fmt.Errorf("failed to generate next number: %w", err)
	}

	newVoucher, err := voucher.New(
		cmd.VoucherId,
		cmd.SobId,
		p,
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
		cmd.TransactionDate,
		lineItems,
	)
	if err != nil {
		return err
	}

	return h.repo.CreateVoucher(ctx, newVoucher)
}
