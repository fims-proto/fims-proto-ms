package db

import (
	"context"
	"errors"
	"time"

	"github/fims-proto/fims-proto-ms/internal/common/data/converter"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/common/data/filterable"
	"github/fims-proto/fims-proto-ms/internal/common/data/pageable"
	"github/fims-proto/fims-proto-ms/internal/common/datasource"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/query"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/transaction_date"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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
	return data.SearchEntities(ctx, pageRequest, accountPO{}, accountPOToDTO, r.dataSource.GetConnection(ctx).Preload("DimensionCategories"))
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
		Preload("JournalLines.Account.DimensionCategories").
		Preload("JournalLines.DimensionOptions").
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

// LedgerEntriesByPeriodRange queries journal lines for a specific account within a fiscal year/period number range.
// Only journal lines from posted journals are included, matching the behaviour of the former ledger entry table.
// Filters are applied at SQL level for performance and supports pagination.
func (r GeneralLedgerPostgresReadRepository) LedgerEntriesByPeriodRange(
	ctx context.Context,
	sobId, accountId uuid.UUID,
	fromFiscalYear, fromPeriodNumber, toFiscalYear, toPeriodNumber int,
	pageRequest data.PageRequest,
) (data.Page[query.LedgerEntry], error) {
	db := r.dataSource.GetConnection(ctx)

	type journalLineRow struct {
		Id              uuid.UUID       `gorm:"column:id"`
		JournalId       uuid.UUID       `gorm:"column:journal_id"`
		Text            string          `gorm:"column:text"`
		Amount          decimal.Decimal `gorm:"column:amount"`
		CreatedAt       time.Time       `gorm:"column:created_at"`
		UpdatedAt       time.Time       `gorm:"column:updated_at"`
		DocumentNumber  string          `gorm:"column:document_number"`
		TransactionDate time.Time       `gorm:"column:transaction_date"`
	}

	baseQ := db.Table("a_journal_lines").
		Select("a_journal_lines.id, a_journal_lines.journal_id, a_journal_lines.text, a_journal_lines.amount, "+
			"a_journal_lines.created_at, a_journal_lines.updated_at, "+
			"a_journals.document_number, a_journals.transaction_date").
		Joins("INNER JOIN a_journals ON a_journal_lines.journal_id = a_journals.id").
		Joins("INNER JOIN a_periods ON a_journals.period_id = a_periods.id").
		Where("a_journals.sob_id = ?", sobId).
		Where("a_journal_lines.account_id = ?", accountId).
		Where("a_journals.is_posted = ?", true).
		Where(
			"(a_periods.fiscal_year > ? OR (a_periods.fiscal_year = ? AND a_periods.period_number >= ?)) AND "+
				"(a_periods.fiscal_year < ? OR (a_periods.fiscal_year = ? AND a_periods.period_number <= ?))",
			fromFiscalYear, fromFiscalYear, fromPeriodNumber,
			toFiscalYear, toFiscalYear, toPeriodNumber,
		).
		Order("a_journal_lines.created_at asc")

	var count int64
	if err := baseQ.Session(&gorm.Session{}).Count(&count).Error; err != nil {
		return nil, err
	}

	var rows []journalLineRow
	if err := baseQ.Session(&gorm.Session{}).
		Scopes(pageable.Paging(pageRequest)).
		Find(&rows).Error; err != nil {
		return nil, err
	}

	dtos := make([]query.LedgerEntry, 0, len(rows))
	for _, row := range rows {
		td := transaction_date.TransactionDate{
			Year:  row.TransactionDate.Year(),
			Month: int(row.TransactionDate.Month()),
			Day:   row.TransactionDate.Day(),
		}
		dtos = append(dtos, query.LedgerEntry{
			JournalId:       row.JournalId,
			JournalNumber:   row.DocumentNumber,
			TransactionDate: td,
			Text:            row.Text,
			Amount:          row.Amount,
			CreatedAt:       row.CreatedAt,
			UpdatedAt:       row.UpdatedAt,
		})
	}

	return data.NewPage(dtos, pageRequest, int(count))
}

type ledgerDimensionSummaryRow struct {
	DimensionOptionId uuid.UUID       `gorm:"column:dimension_option_id"`
	TotalAmount       decimal.Decimal `gorm:"column:total_amount"`
}

// LedgerDimensionSummary aggregates journal line amounts grouped by dimension option
// for a specific account and dimension category within a period range.
func (r GeneralLedgerPostgresReadRepository) LedgerDimensionSummary(
	ctx context.Context,
	sobId, accountId, dimensionCategoryId uuid.UUID,
	fromFiscalYear, fromPeriodNumber, toFiscalYear, toPeriodNumber int,
) ([]query.LedgerDimensionSummaryItem, error) {
	db := r.dataSource.GetConnection(ctx)

	var rows []ledgerDimensionSummaryRow
	err := db.Model(&journalLinePO{}).
		Select("a_journal_line_dimension_options.dimension_option_id, SUM(a_journal_lines.amount) AS total_amount").
		Joins("INNER JOIN a_journals ON a_journal_lines.journal_id = a_journals.id").
		Joins("INNER JOIN a_periods ON a_journals.period_id = a_periods.id").
		Joins("INNER JOIN a_journal_line_dimension_options ON a_journal_line_dimension_options.journal_line_id = a_journal_lines.id").
		Joins("INNER JOIN a_dimension_options ON a_dimension_options.id = a_journal_line_dimension_options.dimension_option_id").
		Where("a_journals.sob_id = ?", sobId).
		Where("a_journal_lines.account_id = ?", accountId).
		Where("a_dimension_options.category_id = ?", dimensionCategoryId).
		Where(
			"(a_periods.fiscal_year > ? OR (a_periods.fiscal_year = ? AND a_periods.period_number >= ?)) AND "+
				"(a_periods.fiscal_year < ? OR (a_periods.fiscal_year = ? AND a_periods.period_number <= ?))",
			fromFiscalYear, fromFiscalYear, fromPeriodNumber,
			toFiscalYear, toFiscalYear, toPeriodNumber,
		).
		Group("a_journal_line_dimension_options.dimension_option_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]query.LedgerDimensionSummaryItem, 0, len(rows))
	for _, row := range rows {
		result = append(result, query.LedgerDimensionSummaryItem{
			DimensionOptionId: row.DimensionOptionId,
			TotalAmount:       row.TotalAmount,
		})
	}

	return result, nil
}
