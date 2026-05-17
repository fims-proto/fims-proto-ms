package command

import (
	"context"
	"fmt"
	"sync"

	"github/fims-proto/fims-proto-ms/internal/report/domain"
	"github/fims-proto/fims-proto-ms/internal/report/domain/generator"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report"
	"github/fims-proto/fims-proto-ms/internal/report/domain/service"
	"github/fims-proto/fims-proto-ms/internal/report/domain/validator"

	"github.com/google/uuid"
)

type GenerateForPeriodCmd struct {
	SobId    uuid.UUID
	PeriodId uuid.UUID
}

type GenerateForPeriodHandler struct {
	repo                 domain.Repository
	generalLedgerService service.GeneralLedgerService
}

func NewGenerateForPeriodHandler(repo domain.Repository, generalLedgerService service.GeneralLedgerService) GenerateForPeriodHandler {
	if repo == nil {
		panic("nil repository")
	}
	if generalLedgerService == nil {
		panic("nil general ledger service")
	}
	return GenerateForPeriodHandler{
		repo:                 repo,
		generalLedgerService: generalLedgerService,
	}
}

func (h GenerateForPeriodHandler) Handle(ctx context.Context, cmd GenerateForPeriodCmd) error {
	templates, err := h.repo.ReadTemplatesBySobId(ctx, cmd.SobId)
	if err != nil {
		return fmt.Errorf("failed to read report templates: %w", err)
	}

	var (
		mu   sync.Mutex
		errs []error
		wg   sync.WaitGroup
	)

	for _, tmpl := range templates {
		wg.Add(1)
		go func(tmpl *report.Report) {
			defer wg.Done()
			if err = h.upsertInstanceForTemplate(ctx, tmpl, cmd.PeriodId); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
		}(tmpl)
	}

	wg.Wait()

	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

func (h GenerateForPeriodHandler) upsertInstanceForTemplate(ctx context.Context, tmpl *report.Report, periodId uuid.UUID) error {
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		existing, err := h.repo.ReadInstanceBySobClassAndPeriod(txCtx, tmpl.SobId(), tmpl.Class(), periodId)
		if err != nil {
			return fmt.Errorf("failed to check existing report instance: %w", err)
		}

		if existing != nil {
			return h.repo.UpdateReport(txCtx, existing.Id(), func(r *report.Report) (*report.Report, error) {
				v := validator.NewValidatorFactory(r.Class())
				g := generator.NewGenerator(r, h.generalLedgerService, v)
				if err := g.Regenerate(txCtx); err != nil {
					return nil, fmt.Errorf("failed to regenerate report: %w", err)
				}
				return g.Report(), nil
			})
		}

		v := validator.NewValidatorFactory(tmpl.Class())
		g := generator.NewGenerator(tmpl, h.generalLedgerService, v)
		newReport, err := g.Generate(txCtx, uuid.New(), periodId, "", nil)
		if err != nil {
			return fmt.Errorf("failed to generate report: %w", err)
		}
		return h.repo.CreateReports(txCtx, []*report.Report{newReport})
	})
}
