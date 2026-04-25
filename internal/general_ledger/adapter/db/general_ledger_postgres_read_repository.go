package db

import (
	"context"
	"errors"
	"sort"
	"time"

	"github/fims-proto/fims-proto-ms/internal/common/data/converter"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/common/data/filterable"
	"github/fims-proto/fims-proto-ms/internal/common/data/pageable"
	"github/fims-proto/fims-proto-ms/internal/common/datasource"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/query"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account/class"
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
	return data.SearchEntities(ctx, pageRequest, journalPO{}, journalPOToDTO, r.dataSource.GetConnection(ctx).InnerJoins("Period"))
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

func (r GeneralLedgerPostgresReadRepository) ClosingJournalIdBySobAndPeriod(ctx context.Context, sobId uuid.UUID, fiscalYear, periodNumber int, journalType string) (*uuid.UUID, error) {
	db := r.dataSource.GetConnection(ctx)

	var row struct{ Id uuid.UUID }
	err := db.Model(&journalPO{}).
		Select("journals.id").
		InnerJoins("Period").
		Where("journals.sob_id = ? AND journals.journal_type = ?", sobId, journalType).
		Where("\"Period\".fiscal_year = ? AND \"Period\".period_number = ?", fiscalYear, periodNumber).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &row.Id, nil
}

func addSobFilter(sobId uuid.UUID, pageRequest data.PageRequest) {
	if sobId != uuid.Nil {
		sobIdFilter, _ := filterable.NewFilter("sobId", filterable.OptEq, sobId.String())
		pageRequest.AddAndFilterable(filterable.NewFilterableAtom(sobIdFilter))
	}
}

// LedgersByPeriodRange aggregates journal line amounts grouped by account for a SoB within a period range.
// Always queries from journal_lines (the authoritative source) — never from the ledgers snapshot table.
// Two queries are executed and merged in Go:
//  1. Opening balances (all posted lines strictly before fromPeriod)
//  2. Period activity (posted lines within [fromPeriod, toPeriod])
//
// When dimensionOptionId is non-nil, only journal lines tagged with that dimension option are included.
// Account details are fetched in a single batch query after merging.
func (r GeneralLedgerPostgresReadRepository) LedgersByPeriodRange(
	ctx context.Context,
	sobId uuid.UUID,
	fromFiscalYear, fromPeriodNumber, toFiscalYear, toPeriodNumber int,
	dimensionOptionId *uuid.UUID,
	pageRequest data.PageRequest,
) (data.Page[query.Ledger], error) {
	db := r.dataSource.GetConnection(ctx)

	commonJoins := func(q *gorm.DB) *gorm.DB {
		q = q.
			Joins("INNER JOIN journals ON journal_lines.journal_id = journals.id").
			Joins("INNER JOIN periods ON journals.period_id = periods.id").
			Where("journals.sob_id = ?", sobId).
			Where("journals.is_posted = ?", true)
		if dimensionOptionId != nil {
			q = q.Joins(
				"INNER JOIN journal_line_dimension_options jldo ON jldo.journal_line_id = journal_lines.id AND jldo.dimension_option_id = ?",
				*dimensionOptionId,
			)
		}
		return q
	}

	// Query 1: opening balances — all periods strictly before fromPeriod
	var openingRows []dimensionOptionLedgerOpeningRow
	openingQ := commonJoins(db.Model(&journalLinePO{})).
		Select("journal_lines.account_id, SUM(journal_lines.amount) AS opening_amount").
		Where(
			"(periods.fiscal_year < ? OR (periods.fiscal_year = ? AND periods.period_number < ?))",
			fromFiscalYear, fromFiscalYear, fromPeriodNumber,
		).
		Group("journal_lines.account_id")
	if err := openingQ.Scan(&openingRows).Error; err != nil {
		return nil, err
	}

	// Query 2: period activity — lines within [fromPeriod, toPeriod]
	var periodRows []dimensionOptionLedgerPeriodRow
	periodQ := commonJoins(db.Model(&journalLinePO{})).
		Select(
			"journal_lines.account_id, "+
				"SUM(CASE WHEN journal_lines.amount > 0 THEN journal_lines.amount ELSE 0 END) AS period_debit, "+
				"SUM(CASE WHEN journal_lines.amount < 0 THEN ABS(journal_lines.amount) ELSE 0 END) AS period_credit, "+
				"SUM(journal_lines.amount) AS period_amount",
		).
		Where(
			"(periods.fiscal_year > ? OR (periods.fiscal_year = ? AND periods.period_number >= ?)) AND "+
				"(periods.fiscal_year < ? OR (periods.fiscal_year = ? AND periods.period_number <= ?))",
			fromFiscalYear, fromFiscalYear, fromPeriodNumber,
			toFiscalYear, toFiscalYear, toPeriodNumber,
		).
		Group("journal_lines.account_id")
	if err := periodQ.Scan(&periodRows).Error; err != nil {
		return nil, err
	}

	// Merge results in Go, keyed by account ID
	type mergedItem struct {
		opening decimal.Decimal
		debit   decimal.Decimal
		credit  decimal.Decimal
		period  decimal.Decimal
	}
	merged := make(map[uuid.UUID]*mergedItem)

	for _, row := range openingRows {
		merged[row.AccountId] = &mergedItem{opening: row.OpeningAmount}
	}
	for _, row := range periodRows {
		item, ok := merged[row.AccountId]
		if !ok {
			item = &mergedItem{}
			merged[row.AccountId] = item
		}
		item.debit = row.PeriodDebit
		item.credit = row.PeriodCredit
		item.period = row.PeriodAmount
	}

	if len(merged) == 0 {
		return data.NewPage([]query.Ledger{}, pageRequest, 0)
	}

	// Batch fetch account details
	accountIds := make([]uuid.UUID, 0, len(merged))
	for id := range merged {
		accountIds = append(accountIds, id)
	}
	var accountPos []accountPO
	if err := db.Where("id IN ?", accountIds).Find(&accountPos).Error; err != nil {
		return nil, err
	}
	accountMap := make(map[uuid.UUID]query.Account, len(accountPos))
	for _, po := range accountPos {
		accountMap[po.Id] = accountPOToDTO(po)
	}

	// Build result slice
	dtos := make([]query.Ledger, 0, len(merged))
	for id, item := range merged {
		acc := accountMap[id]
		dtos = append(dtos, query.Ledger{
			SobId:         sobId,
			AccountId:     id,
			Account:       acc,
			OpeningAmount: item.opening,
			PeriodDebit:   item.debit,
			PeriodCredit:  item.credit,
			PeriodAmount:  item.period,
			EndingAmount:  item.opening.Add(item.period),
		})
	}
	sort.Slice(dtos, func(i, j int) bool {
		return dtos[i].Account.RawAccountNumber < dtos[j].Account.RawAccountNumber
	})

	// Apply pagination in Go
	total := len(dtos)
	offset := pageRequest.Offset()
	if offset >= total {
		return data.NewPage([]query.Ledger{}, pageRequest, total)
	}
	end := offset + pageRequest.PageSize()
	if end > total {
		end = total
	}
	return data.NewPage(dtos[offset:end], pageRequest, total)
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

// LedgerEntriesByPeriodRange queries journal lines for a specific account within a fiscal year/period number range.
// Only journal lines from posted journals are included, matching the behaviour of the former ledger entry table.
// Filters are applied at SQL level for performance and supports pagination.
func (r GeneralLedgerPostgresReadRepository) LedgerEntriesByPeriodRange(
	ctx context.Context,
	sobId uuid.UUID,
	accountId *uuid.UUID,
	fromFiscalYear, fromPeriodNumber, toFiscalYear, toPeriodNumber int,
	dimensionOptionId *uuid.UUID,
	pageRequest data.PageRequest,
) (data.Page[query.LedgerEntry], error) {
	db := r.dataSource.GetConnection(ctx)

	baseQ := db.Table("journal_lines").
		Select("journal_lines.id, journal_lines.journal_id, journal_lines.text, journal_lines.amount, "+
			"journal_lines.created_at, journal_lines.updated_at, "+
			"journals.document_number, journals.transaction_date").
		Joins("INNER JOIN journals ON journal_lines.journal_id = journals.id").
		Joins("INNER JOIN periods ON journals.period_id = periods.id").
		Where("journals.sob_id = ?", sobId).
		Where("journals.is_posted = ?", true).
		Where(
			"(periods.fiscal_year > ? OR (periods.fiscal_year = ? AND periods.period_number >= ?)) AND "+
				"(periods.fiscal_year < ? OR (periods.fiscal_year = ? AND periods.period_number <= ?))",
			fromFiscalYear, fromFiscalYear, fromPeriodNumber,
			toFiscalYear, toFiscalYear, toPeriodNumber,
		)

	if accountId != nil {
		baseQ = baseQ.Where("journal_lines.account_id = ?", *accountId)
	}

	if dimensionOptionId != nil {
		baseQ = baseQ.Joins(
			"INNER JOIN journal_line_dimension_options jldo ON jldo.journal_line_id = journal_lines.id AND jldo.dimension_option_id = ?",
			*dimensionOptionId,
		)
	}

	baseQ = baseQ.Order("journal_lines.created_at asc")

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

type ledgerDimensionSummaryOpeningRow struct {
	DimensionOptionId uuid.UUID       `gorm:"column:dimension_option_id"`
	OpeningAmount     decimal.Decimal `gorm:"column:opening_amount"`
}

type ledgerDimensionSummaryPeriodRow struct {
	DimensionOptionId   uuid.UUID       `gorm:"column:dimension_option_id"`
	DimensionOptionName string          `gorm:"column:dimension_option_name"`
	PeriodDebit         decimal.Decimal `gorm:"column:period_debit"`
	PeriodCredit        decimal.Decimal `gorm:"column:period_credit"`
	PeriodAmount        decimal.Decimal `gorm:"column:period_amount"`
}

// LedgersByAccountAndDimensionOption aggregates journal line amounts grouped by dimension option
// for a specific account and dimension category within a period range.
// Two queries are executed and merged in Go:
//  1. Opening balances (all posted lines strictly before fromPeriod)
//  2. Period activity (posted lines within [fromPeriod, toPeriod])
func (r GeneralLedgerPostgresReadRepository) LedgersByAccountAndDimensionOption(
	ctx context.Context,
	sobId uuid.UUID,
	accountId *uuid.UUID,
	dimensionCategoryId uuid.UUID,
	fromFiscalYear, fromPeriodNumber, toFiscalYear, toPeriodNumber int,
	pageRequest data.PageRequest,
) (data.Page[query.LedgerDimensionSummaryItem], error) {
	db := r.dataSource.GetConnection(ctx)

	commonJoins := func(q *gorm.DB) *gorm.DB {
		q = q.
			Joins("INNER JOIN journals ON journal_lines.journal_id = journals.id").
			Joins("INNER JOIN periods ON journals.period_id = periods.id").
			Joins("INNER JOIN journal_line_dimension_options ON journal_line_dimension_options.journal_line_id = journal_lines.id").
			Joins("INNER JOIN dimension_options ON dimension_options.id = journal_line_dimension_options.dimension_option_id").
			Where("journals.sob_id = ?", sobId).
			Where("dimension_options.category_id = ?", dimensionCategoryId).
			Where("journals.is_posted = ?", true)
		if accountId != nil {
			q = q.Where("journal_lines.account_id = ?", *accountId)
		}
		return q
	}

	// Query 1: opening balances — all periods strictly before fromPeriod
	var openingRows []ledgerDimensionSummaryOpeningRow
	openingQ := commonJoins(db.Model(&journalLinePO{})).
		Select("journal_line_dimension_options.dimension_option_id, SUM(journal_lines.amount) AS opening_amount").
		Where(
			"(periods.fiscal_year < ? OR (periods.fiscal_year = ? AND periods.period_number < ?))",
			fromFiscalYear, fromFiscalYear, fromPeriodNumber,
		).
		Group("journal_line_dimension_options.dimension_option_id")
	if err := openingQ.Scan(&openingRows).Error; err != nil {
		return nil, err
	}

	// Query 2: period activity — lines within [fromPeriod, toPeriod]
	var periodRows []ledgerDimensionSummaryPeriodRow
	periodQ := commonJoins(db.Model(&journalLinePO{})).
		Select(
			"journal_line_dimension_options.dimension_option_id, "+
				"dimension_options.name AS dimension_option_name, "+
				"SUM(CASE WHEN journal_lines.amount > 0 THEN journal_lines.amount ELSE 0 END) AS period_debit, "+
				"SUM(CASE WHEN journal_lines.amount < 0 THEN ABS(journal_lines.amount) ELSE 0 END) AS period_credit, "+
				"SUM(journal_lines.amount) AS period_amount",
		).
		Where(
			"(periods.fiscal_year > ? OR (periods.fiscal_year = ? AND periods.period_number >= ?)) AND "+
				"(periods.fiscal_year < ? OR (periods.fiscal_year = ? AND periods.period_number <= ?))",
			fromFiscalYear, fromFiscalYear, fromPeriodNumber,
			toFiscalYear, toFiscalYear, toPeriodNumber,
		).
		Group("journal_line_dimension_options.dimension_option_id, dimension_options.name")
	if err := periodQ.Scan(&periodRows).Error; err != nil {
		return nil, err
	}

	// Merge results in Go, keyed by dimension option ID
	type mergedItem struct {
		name    string
		opening decimal.Decimal
		debit   decimal.Decimal
		credit  decimal.Decimal
		period  decimal.Decimal
	}
	merged := make(map[uuid.UUID]*mergedItem)

	for _, row := range openingRows {
		merged[row.DimensionOptionId] = &mergedItem{opening: row.OpeningAmount}
	}
	for _, row := range periodRows {
		item, ok := merged[row.DimensionOptionId]
		if !ok {
			item = &mergedItem{}
			merged[row.DimensionOptionId] = item
		}
		item.name = row.DimensionOptionName
		item.debit = row.PeriodDebit
		item.credit = row.PeriodCredit
		item.period = row.PeriodAmount
	}

	// For opening-only options, fetch the name from dimension_options table
	namelessIds := make([]uuid.UUID, 0)
	for id, item := range merged {
		if item.name == "" {
			namelessIds = append(namelessIds, id)
		}
	}
	if len(namelessIds) > 0 {
		type nameRow struct {
			Id   uuid.UUID `gorm:"column:id"`
			Name string    `gorm:"column:name"`
		}
		var nameRows []nameRow
		if err := db.Table("dimension_options").
			Select("id, name").
			Where("id IN ?", namelessIds).
			Scan(&nameRows).Error; err != nil {
			return nil, err
		}
		for _, nr := range nameRows {
			if item, ok := merged[nr.Id]; ok {
				item.name = nr.Name
			}
		}
	}

	// Build sorted slice
	dtos := make([]query.LedgerDimensionSummaryItem, 0, len(merged))
	for id, item := range merged {
		dtos = append(dtos, query.LedgerDimensionSummaryItem{
			DimensionOptionId:   id,
			DimensionOptionName: item.name,
			OpeningAmount:       item.opening,
			PeriodDebit:         item.debit,
			PeriodCredit:        item.credit,
			PeriodAmount:        item.period,
			EndingAmount:        item.opening.Add(item.period),
		})
	}
	sort.Slice(dtos, func(i, j int) bool {
		return dtos[i].DimensionOptionName < dtos[j].DimensionOptionName
	})

	// Apply pagination in Go
	total := len(dtos)
	offset := pageRequest.Offset()
	if offset >= total {
		return data.NewPage([]query.LedgerDimensionSummaryItem{}, pageRequest, total)
	}
	end := offset + pageRequest.PageSize()
	if end > total {
		end = total
	}
	return data.NewPage(dtos[offset:end], pageRequest, total)
}

func (r GeneralLedgerPostgresReadRepository) ProfitAndLossLedgersHavingBalanceInPeriod(
	ctx context.Context,
	sobId, periodId uuid.UUID,
) ([]query.Ledger, error) {
	db := r.dataSource.GetConnection(ctx)

	var pos []ledgerPO
	err := db.Where(ledgerPO{SobId: sobId, PeriodId: periodId}).
		Where("ending_amount <> 0").
		InnerJoins("Account", db.Where(accountPO{Class: int(class.ProfitsAndLosses), IsLeaf: true})).
		Find(&pos).Error
	if err != nil {
		return nil, err
	}

	return converter.POsToDTOs(pos, ledgerPOToDTO), nil
}

func (r GeneralLedgerPostgresReadRepository) PeriodById(
	ctx context.Context,
	sobId, periodId uuid.UUID,
) (query.Period, error) {
	db := r.dataSource.GetConnection(ctx)

	var po periodPO
	if err := db.Where(periodPO{Id: periodId, SobId: sobId}).First(&po).Error; err != nil {
		return query.Period{}, err
	}

	return periodPOToDTO(po), nil
}

type dimensionOptionLedgerOpeningRow struct {
	AccountId     uuid.UUID       `gorm:"column:account_id"`
	OpeningAmount decimal.Decimal `gorm:"column:opening_amount"`
}

type dimensionOptionLedgerPeriodRow struct {
	AccountId    uuid.UUID       `gorm:"column:account_id"`
	PeriodDebit  decimal.Decimal `gorm:"column:period_debit"`
	PeriodCredit decimal.Decimal `gorm:"column:period_credit"`
	PeriodAmount decimal.Decimal `gorm:"column:period_amount"`
}

func (r GeneralLedgerPostgresReadRepository) LedgerByRawAccountNumberInPeriod(
	ctx context.Context,
	sobId uuid.UUID,
	rawAccountNumber string,
	periodId uuid.UUID,
) (*query.Ledger, error) {
	db := r.dataSource.GetConnection(ctx)

	var pos []ledgerPO
	err := db.Where(ledgerPO{SobId: sobId, PeriodId: periodId}).
		InnerJoins("Account", db.Where(accountPO{RawAccountNumber: rawAccountNumber})).
		Find(&pos).Error
	if err != nil {
		return nil, err
	}

	if len(pos) == 0 {
		return nil, nil
	}

	return new(ledgerPOToDTO(pos[0])), nil
}
