package command

import (
	"context"
	"time"

	"github/fims-proto/fims-proto-ms/internal/journal/domain/journal_entry"
	"github/fims-proto/fims-proto-ms/internal/journal/domain/line_item"

	"github/fims-proto/fims-proto-ms/internal/journal/app/service"

	"github/fims-proto/fims-proto-ms/internal/journal/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type CreateJournalEntryCmd struct {
	EntryId            uuid.UUID
	SobId              uuid.UUID
	HeaderText         string
	JournalType        string
	AttachmentQuantity int
	LineItems          []LineItemCmd
	Creator            uuid.UUID
	TransactionTime    time.Time
}

type CreateJournalEntryHandler struct {
	repo             domain.Repository
	accountService   service.AccountService
	numberingService service.NumberingService
}

func NewCreateJournalEntryHandler(
	repo domain.Repository,
	accountService service.AccountService,
	numberingService service.NumberingService,
) CreateJournalEntryHandler {
	if repo == nil {
		panic("nil repo")
	}

	if accountService == nil {
		panic("nil account service")
	}

	if numberingService == nil {
		panic("nil numbering service")
	}

	return CreateJournalEntryHandler{
		repo:             repo,
		accountService:   accountService,
		numberingService: numberingService,
	}
}

func (h CreateJournalEntryHandler) Handle(ctx context.Context, cmd CreateJournalEntryCmd) error {
	// read period by transaction time
	period, err := h.accountService.ReadPeriodByTime(ctx, cmd.SobId, cmd.TransactionTime)
	if err != nil {
		return errors.Wrap(err, "failed to read period by transaction time")
	}

	if period.IsClosed {
		return errors.New("period is closed")
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
	identifier, err := h.numberingService.GenerateIdentifier(ctx, period.PeriodId, cmd.JournalType)
	if err != nil {
		return errors.Wrap(err, "unable to generate next number")
	}

	newEntry, err := journal_entry.New(
		cmd.SobId,
		cmd.EntryId,
		period.PeriodId,
		cmd.HeaderText,
		cmd.JournalType,
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

	return h.repo.CreateJournalEntry(ctx, newEntry)
}
