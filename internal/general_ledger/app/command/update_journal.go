package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/transaction_date"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/journal"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"

	"github.com/google/uuid"
)

type UpdateJournalCmd struct {
	JournalId       uuid.UUID
	HeaderText      string
	JournalLines    []JournalLineCmd
	TransactionDate transaction_date.TransactionDate
	Updater         uuid.UUID
}

type UpdateJournalHandler struct {
	repo             domain.Repository
	numberingService service.NumberingService
	dimensionService service.DimensionService
}

func NewUpdateJournalHandler(repo domain.Repository, numberingService service.NumberingService, dimensionService service.DimensionService) UpdateJournalHandler {
	if repo == nil {
		panic("nil repo")
	}

	if numberingService == nil {
		panic("nil numbering service")
	}

	if dimensionService == nil {
		panic("nil dimension service")
	}

	return UpdateJournalHandler{
		repo:             repo,
		numberingService: numberingService,
		dimensionService: dimensionService,
	}
}

func (h UpdateJournalHandler) Handle(ctx context.Context, cmd UpdateJournalCmd) error {
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		return h.updateJournal(txCtx, cmd)
	})
}

func (h UpdateJournalHandler) updateJournal(ctx context.Context, cmd UpdateJournalCmd) error {
	return h.repo.UpdateEntireJournal(
		ctx,
		cmd.JournalId,
		func(j *journal.Journal) (*journal.Journal, error) {
			// update journal lines
			if len(cmd.JournalLines) > 0 {
				journalLines, err := prepareJournalLines(ctx, h.repo, h.dimensionService, j.SobId(), cmd.JournalLines)
				if err != nil {
					return nil, fmt.Errorf("failed to prepare journal lines: %w", err)
				}

				if err = j.UpdateJournalLines(journalLines, cmd.Updater); err != nil {
					return nil, err
				}
			}

			// update transaction date (and period and document number, if needed)
			if !cmd.TransactionDate.IsZero() {
				p, err := readPeriodIdAndCheck(ctx, h.repo, h.numberingService, j.SobId(), cmd.TransactionDate)
				if err != nil {
					return nil, fmt.Errorf("failed to read or create period: %w", err)
				}

				if p.Id() != j.PeriodId() {
					// different period, need to regenerate journal id
					identifier, err := h.numberingService.GenerateIdentifier(ctx, p.Id(), j.JournalType().String())
					if err != nil {
						return nil, fmt.Errorf("failed to re-generate journal number: %w", err)
					}
					if err = j.UpdatePeriodAndDocumentNumber(p.Id(), identifier, cmd.Updater); err != nil {
						return nil, err
					}
				}

				if err = j.UpdateTransactionDate(cmd.TransactionDate, cmd.Updater); err != nil {
					return nil, err
				}
			}

			if cmd.HeaderText != "" {
				if err := j.UpdateHeaderText(cmd.HeaderText, cmd.Updater); err != nil {
					return nil, err
				}
			}

			return j, nil
		},
	)
}
