package command

import (
	"context"
	"time"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/query"

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
	readModel        query.GeneralLedgerReadModel
	numberingService service.NumberingService
}

func NewCreateVoucherHandler(
	repo domain.Repository,
	readModel query.GeneralLedgerReadModel,
	numberingService service.NumberingService,
) CreateVoucherHandler {
	if repo == nil {
		panic("nil repo")
	}

	if readModel == nil {
		panic("nil read model")
	}

	if numberingService == nil {
		panic("nil numbering service")
	}

	return CreateVoucherHandler{
		repo:             repo,
		readModel:        readModel,
		numberingService: numberingService,
	}
}

func (h CreateVoucherHandler) Handle(ctx context.Context, cmd CreateVoucherCmd) error {
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		periodId, err := readOrCreatePeriodForVoucher(txCtx, h.repo, h.readModel, h.numberingService, cmd.SobId, cmd.TransactionTime)
		if err != nil {
			return errors.Wrap(err, "failed to read or create period")
		}

		return h.createVoucher(txCtx, cmd, periodId)
	})
}

func (h CreateVoucherHandler) createVoucher(ctx context.Context, cmd CreateVoucherCmd, periodId uuid.UUID) error {
	// prepare line items
	lineItems, err := prepareLineItems(ctx, h.readModel, cmd.SobId, cmd.LineItems)
	if err != nil {
		return errors.Wrap(err, "failed to prepare line items")
	}

	// get document number
	identifier, err := h.numberingService.GenerateIdentifier(ctx, periodId, cmd.VoucherType)
	if err != nil {
		return errors.Wrap(err, "unable to generate next number")
	}

	newVoucher, err := voucher.New(
		cmd.SobId,
		cmd.VoucherId,
		periodId,
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
