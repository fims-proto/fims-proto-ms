package db

import (
	"context"
	"errors"

	"github/fims-proto/fims-proto-ms/internal/common/data/converter"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/common/data/filterable"
	"github/fims-proto/fims-proto-ms/internal/common/datasource"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/query"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GeneralLedgerPostgresReadRepository struct {
	dataSource datasource.DataSource
}

func NewGeneralLedgerPostgresReadRepository(dataSource datasource.DataSource) *GeneralLedgerPostgresReadRepository {
	return &GeneralLedgerPostgresReadRepository{
		dataSource: dataSource,
	}
}

func (r GeneralLedgerPostgresReadRepository) SearchAccounts(
	ctx context.Context,
	sobId uuid.UUID,
	pageRequest data.PageRequest,
) (data.Page[query.Account], error) {
	addSobFilter(sobId, pageRequest)
	return data.SearchEntities(ctx, pageRequest, accountPO{}, accountPOToDTO, r.dataSource.GetConnection(ctx))
}

func (r GeneralLedgerPostgresReadRepository) SearchLedgers(
	ctx context.Context,
	sobId uuid.UUID,
	pageRequest data.PageRequest,
) (data.Page[query.Ledger], error) {
	addSobFilter(sobId, pageRequest)
	return data.SearchEntities(ctx, pageRequest, ledgerPO{}, ledgerPOToDTO, r.dataSource.GetConnection(ctx).InnerJoins("Account"))
}

func (r GeneralLedgerPostgresReadRepository) SearchPeriods(
	ctx context.Context,
	sobId uuid.UUID,
	pageRequest data.PageRequest,
) (data.Page[query.Period], error) {
	addSobFilter(sobId, pageRequest)
	return data.SearchEntities(ctx, pageRequest, periodPO{}, periodPOToDTO, r.dataSource.GetConnection(ctx))
}

func (r GeneralLedgerPostgresReadRepository) SearchJournals(
	ctx context.Context,
	sobId uuid.UUID,
	pageRequest data.PageRequest,
) (data.Page[query.Journal], error) {
	addSobFilter(sobId, pageRequest)
	return data.SearchEntities(ctx, pageRequest, journalPO{}, journalPOToDTO, r.dataSource.GetConnection(ctx).Preload("JournalLines.Account").InnerJoins("Period"))
}

func (r GeneralLedgerPostgresReadRepository) CurrentPeriod(ctx context.Context, sobId uuid.UUID) (query.Period, error) {
	db := r.dataSource.GetConnection(ctx)

	var po periodPO
	if err := db.Where(periodPO{SobId: sobId, IsCurrent: true}).
		First(&po).Error; err != nil {
		return query.Period{}, err
	}

	return periodPOToDTO(po), nil
}

func (r GeneralLedgerPostgresReadRepository) FirstPeriod(ctx context.Context, sobId uuid.UUID) (query.Period, error) {
	db := r.dataSource.GetConnection(ctx)

	var po periodPO
	err := db.Order("fiscal_year asc, period_number asc").Where(periodPO{SobId: sobId}).First(&po).Error
	if err == nil {
		return periodPOToDTO(po), nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return query.Period{}, commonErrors.ErrRecordNotFound()
	}
	return query.Period{}, err
}

func (r GeneralLedgerPostgresReadRepository) JournalById(ctx context.Context, journalId uuid.UUID) (query.Journal, error) {
	db := r.dataSource.GetConnection(ctx)

	po := journalPO{Id: journalId}
	if err := db.
		Preload("JournalLines.Account").
		Preload("Period").
		First(&po).Error; err != nil {
		return query.Journal{}, err
	}

	return journalPOToDTO(po), nil
}

func addSobFilter(sobId uuid.UUID, pageRequest data.PageRequest) {
	if sobId != uuid.Nil {
		sobIdFilter, _ := filterable.NewFilter("sobId", filterable.OptEq, sobId.String())
		pageRequest.AddAndFilterable(filterable.NewFilterableAtom(sobIdFilter))
	}
}

// LedgersByPeriodRange queries ledgers for a specific account within a fiscal year/period number range
// Filters are applied at SQL level for performance
func (r GeneralLedgerPostgresReadRepository) LedgersByPeriodRange(
	ctx context.Context,
	sobId, accountId uuid.UUID,
	fromFiscalYear, fromPeriodNumber, toFiscalYear, toPeriodNumber int,
) ([]query.Ledger, error) {
	db := r.dataSource.GetConnection(ctx)

	var pos []ledgerPO
	q := db.Model(&ledgerPO{SobId: sobId, AccountId: accountId}).
		InnerJoins("Account").
		InnerJoins("Period").
		Where("a_ledgers.sob_id = ?", sobId).
		Where("a_ledgers.account_id = ?", accountId).
		Where(
			"(fiscal_year > ? OR (fiscal_year = ? AND period_number >= ?)) AND "+
				"(fiscal_year < ? OR (fiscal_year = ? AND period_number <= ?))",
			fromFiscalYear, fromFiscalYear, fromPeriodNumber, toFiscalYear, toFiscalYear, toPeriodNumber,
		).
		Order("fiscal_year asc, period_number asc")

	if err := q.Find(&pos).Error; err != nil {
		return nil, err
	}

	if len(pos) == 0 {
		return nil, commonErrors.ErrRecordNotFound()
	}

	result := converter.POsToDTOs(pos, ledgerPOToDTO)

	return result, nil
}

// AllLedgersByPeriodRange queries all ledgers across all accounts for a SoB within a fiscal year/period number range.
// Results are ordered by account_number asc, fiscal_year asc, period_number asc for aggregation by the caller.
func (r GeneralLedgerPostgresReadRepository) AllLedgersByPeriodRange(
	ctx context.Context,
	sobId uuid.UUID,
	fromFiscalYear, fromPeriodNumber, toFiscalYear, toPeriodNumber int,
) ([]query.Ledger, error) {
	db := r.dataSource.GetConnection(ctx)

	var pos []ledgerPO
	q := db.Model(&ledgerPO{}).
		InnerJoins("Account").
		InnerJoins("Period").
		Where("a_ledgers.sob_id = ?", sobId).
		Where(
			"(fiscal_year > ? OR (fiscal_year = ? AND period_number >= ?)) AND "+
				"(fiscal_year < ? OR (fiscal_year = ? AND period_number <= ?))",
			fromFiscalYear, fromFiscalYear, fromPeriodNumber, toFiscalYear, toFiscalYear, toPeriodNumber,
		).
		Order("account_number asc, fiscal_year asc, period_number asc")

	if err := q.Find(&pos).Error; err != nil {
		return nil, err
	}

	return converter.POsToDTOs(pos, ledgerPOToDTO), nil
}

// CheckPeriodContinuity verifies that all periods in the range [fromFiscalYear, fromPeriodNumber] to [toFiscalYear, toPeriodNumber] exist
// Returns error if any period in the range is missing
func (r GeneralLedgerPostgresReadRepository) CheckPeriodContinuity(
	ctx context.Context,
	sobId uuid.UUID,
	fromFiscalYear, fromPeriodNumber, toFiscalYear, toPeriodNumber int,
) error {
	db := r.dataSource.GetConnection(ctx)

	// Calculate expected number of periods in range
	expectedCount := (toFiscalYear-fromFiscalYear)*12 + (toPeriodNumber - fromPeriodNumber) + 1

	// Count actual periods in range using SQL
	var count int64
	db.Model(&periodPO{}).
		Where("sob_id = ?", sobId).
		Where(
			"(fiscal_year > ? OR (fiscal_year = ? AND period_number >= ?)) AND "+
				"(fiscal_year < ? OR (fiscal_year = ? AND period_number <= ?))",
			fromFiscalYear, fromFiscalYear, fromPeriodNumber, toFiscalYear, toFiscalYear, toPeriodNumber,
		).
		Count(&count)
	if err := db.Error; err != nil {
		return err
	}

	if count != int64(expectedCount) {
		return errors.New("periods in range are not continuous")
	}

	return nil
}

// LedgerEntriesByPeriodRange queries ledger entries for a specific account within a fiscal year/period number range
// Filters are applied at SQL level for performance and supports pagination
func (r GeneralLedgerPostgresReadRepository) LedgerEntriesByPeriodRange(
	ctx context.Context,
	sobId, accountId uuid.UUID,
	fromFiscalYear, fromPeriodNumber, toFiscalYear, toPeriodNumber int,
	pageRequest data.PageRequest,
) (data.Page[query.LedgerEntry], error) {
	db := r.dataSource.GetConnection(ctx)

	// Build base query joining ledger_entries with journals
	q := db.Model(&ledgerEntryPO{}).
		InnerJoins("Journal").
		InnerJoins("Period").
		Where("a_ledger_entries.sob_id = ?", sobId).
		Where("a_ledger_entries.account_id = ?", accountId).
		Where(
			"(fiscal_year > ? OR (fiscal_year = ? AND period_number >= ?)) AND "+
				"(fiscal_year < ? OR (fiscal_year = ? AND period_number <= ?))",
			fromFiscalYear, fromFiscalYear, fromPeriodNumber,
			toFiscalYear, toFiscalYear, toPeriodNumber,
		).
		Order("a_ledger_entries.created_at asc")

	// Apply pageable filters and return
	return data.SearchEntities(ctx, pageRequest, ledgerEntryPO{}, ledgerEntryPOToDTO, q)
}
