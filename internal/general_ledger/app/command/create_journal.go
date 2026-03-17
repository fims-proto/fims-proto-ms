package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/transaction_date"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/journal"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/period"

	"github.com/google/uuid"
)

type CreateJournalCmd struct {
	JournalId          uuid.UUID
	SobId              uuid.UUID
	HeaderText         string
	JournalType        string
	ReferenceJournalId uuid.UUID
	AttachmentQuantity int
	JournalLines       []JournalLineCmd
	Creator            uuid.UUID
	TransactionDate    transaction_date.TransactionDate
}

type CreateJournalHandler struct {
	repo             domain.Repository
	numberingService service.NumberingService
	dimensionService service.DimensionService
}

func NewCreateJournalHandler(repo domain.Repository, numberingService service.NumberingService, dimensionService service.DimensionService) CreateJournalHandler {
	if repo == nil {
		panic("nil repo")
	}

	if numberingService == nil {
		panic("nil numbering service")
	}

	if dimensionService == nil {
		panic("nil dimension service")
	}

	return CreateJournalHandler{
		repo:             repo,
		numberingService: numberingService,
		dimensionService: dimensionService,
	}
}

func (h CreateJournalHandler) Handle(ctx context.Context, cmd CreateJournalCmd) error {
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		p, err := readPeriodIdAndCheck(txCtx, h.repo, h.numberingService, cmd.SobId, cmd.TransactionDate)
		if err != nil {
			return fmt.Errorf("failed to read or create period: %w", err)
		}

		return h.createJournal(txCtx, cmd, p)
	})
}

func (h CreateJournalHandler) createJournal(ctx context.Context, cmd CreateJournalCmd, p *period.Period) error {
	journalType := journal.JournalType(cmd.JournalType)

	// Verify reference journal exists in the same SoB when required
	if journalType.RequiresReferenceJournal() {
		exists, err := h.repo.ExistsJournalById(ctx, cmd.SobId, cmd.ReferenceJournalId)
		if err != nil {
			return fmt.Errorf("failed to check reference journal: %w", err)
		}
		if !exists {
			return errors.NewSlugError("journal-referenceJournalNotFound")
		}
	}

	// prepare journal lines
	journalLines, err := prepareJournalLines(ctx, h.repo, h.dimensionService, cmd.SobId, cmd.JournalLines)
	if err != nil {
		return fmt.Errorf("failed to prepare journal lines: %w", err)
	}

	// get document number
	identifier, err := h.numberingService.GenerateIdentifier(ctx, p.Id())
	if err != nil {
		return fmt.Errorf("failed to generate next number: %w", err)
	}

	newJournal, err := journal.New(
		cmd.JournalId,
		cmd.SobId,
		p,
		cmd.HeaderText,
		identifier,
		journalType,
		cmd.ReferenceJournalId,
		cmd.AttachmentQuantity,
		cmd.Creator,
		uuid.Nil,
		uuid.Nil,
		uuid.Nil,
		false,
		false,
		false,
		cmd.TransactionDate,
		journalLines,
	)
	if err != nil {
		return err
	}

	return h.repo.CreateJournal(ctx, newJournal)
}
