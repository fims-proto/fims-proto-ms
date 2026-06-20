package db

import (
	"context"
	"strings"
	"testing"
	"time"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/common/data/filterable"
	"github/fims-proto/fims-proto-ms/internal/common/data/pageable"
	"github/fims-proto/fims-proto-ms/internal/common/data/sortable"
	"github/fims-proto/fims-proto-ms/internal/common/datasource"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type readRepoTestDataSource struct {
	db *gorm.DB
}

func (d readRepoTestDataSource) GetConnection(context.Context) *gorm.DB {
	return d.db
}

func (d readRepoTestDataSource) EnableTransaction(ctx context.Context, txFn func(txCtx context.Context) error) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		return txFn(datasource.WrapInNewContext(ctx, tx))
	})
}

func TestLedgersByPeriodRange_AccountLevelIncludesInitializedOpeningBalance(t *testing.T) {
	ctx := context.Background()
	db := newReadRepoTestDB(t)
	repo := NewGeneralLedgerPostgresReadRepository(readRepoTestDataSource{db: db})

	sobId := uuid.New()
	accountId := uuid.New()
	periodId := uuid.New()
	insertReadRepoTestAccount(t, db, sobId, accountId, "1001")
	insertReadRepoTestPeriod(t, db, sobId, periodId, 2026, 1)
	insertReadRepoTestLedger(t, db, sobId, accountId, periodId, "100", "0", "0", "0", "100")

	page, err := repo.LedgersByPeriodRange(ctx, sobId, 2026, 1, 2026, 1, nil, readRepoTestPageRequest(t))
	require.NoError(t, err)

	ledgers := page.Content()
	require.Len(t, ledgers, 1)
	require.Equal(t, accountId, ledgers[0].AccountId)
	requireDecimalEqual(t, "100", ledgers[0].OpeningAmount)
	requireDecimalEqual(t, "0", ledgers[0].PeriodDebit)
	requireDecimalEqual(t, "0", ledgers[0].PeriodCredit)
	requireDecimalEqual(t, "0", ledgers[0].PeriodAmount)
	requireDecimalEqual(t, "100", ledgers[0].EndingAmount)
}

func TestLedgersByPeriodRange_AccountLevelUsesLedgerSnapshotsAcrossPeriods(t *testing.T) {
	ctx := context.Background()
	db := newReadRepoTestDB(t)
	repo := NewGeneralLedgerPostgresReadRepository(readRepoTestDataSource{db: db})

	sobId := uuid.New()
	accountId := uuid.New()
	firstPeriodId := uuid.New()
	secondPeriodId := uuid.New()
	insertReadRepoTestAccount(t, db, sobId, accountId, "1001")
	insertReadRepoTestPeriod(t, db, sobId, firstPeriodId, 2026, 1)
	insertReadRepoTestPeriod(t, db, sobId, secondPeriodId, 2026, 2)
	insertReadRepoTestLedger(t, db, sobId, accountId, firstPeriodId, "100", "50", "50", "0", "150")
	insertReadRepoTestLedger(t, db, sobId, accountId, secondPeriodId, "150", "-30", "0", "30", "120")

	page, err := repo.LedgersByPeriodRange(ctx, sobId, 2026, 1, 2026, 2, nil, readRepoTestPageRequest(t))
	require.NoError(t, err)

	ledgers := page.Content()
	require.Len(t, ledgers, 1)
	requireDecimalEqual(t, "100", ledgers[0].OpeningAmount)
	requireDecimalEqual(t, "50", ledgers[0].PeriodDebit)
	requireDecimalEqual(t, "30", ledgers[0].PeriodCredit)
	requireDecimalEqual(t, "20", ledgers[0].PeriodAmount)
	requireDecimalEqual(t, "120", ledgers[0].EndingAmount)
}

func TestLedgersByPeriodRange_DimensionOptionIgnoresInitializedOpeningBalance(t *testing.T) {
	ctx := context.Background()
	db := newReadRepoTestDB(t)
	repo := NewGeneralLedgerPostgresReadRepository(readRepoTestDataSource{db: db})

	sobId := uuid.New()
	accountId := uuid.New()
	periodId := uuid.New()
	dimensionOptionId := uuid.New()
	insertReadRepoTestAccount(t, db, sobId, accountId, "1001")
	insertReadRepoTestPeriod(t, db, sobId, periodId, 2026, 1)
	insertReadRepoTestLedger(t, db, sobId, accountId, periodId, "100", "0", "0", "0", "100")

	page, err := repo.LedgersByPeriodRange(ctx, sobId, 2026, 1, 2026, 1, &dimensionOptionId, readRepoTestPageRequest(t))
	require.NoError(t, err)
	require.Empty(t, page.Content())
}

func TestLedgersByPeriodRange_DimensionOptionUsesPriorTaggedActivityAsOpening(t *testing.T) {
	ctx := context.Background()
	db := newReadRepoTestDB(t)
	repo := NewGeneralLedgerPostgresReadRepository(readRepoTestDataSource{db: db})

	sobId := uuid.New()
	accountId := uuid.New()
	firstPeriodId := uuid.New()
	secondPeriodId := uuid.New()
	dimensionOptionId := uuid.New()
	insertReadRepoTestAccount(t, db, sobId, accountId, "1001")
	insertReadRepoTestPeriod(t, db, sobId, firstPeriodId, 2026, 1)
	insertReadRepoTestPeriod(t, db, sobId, secondPeriodId, 2026, 2)
	insertReadRepoTestLedger(t, db, sobId, accountId, firstPeriodId, "100", "40", "40", "0", "140")
	insertReadRepoTestLedger(t, db, sobId, accountId, secondPeriodId, "140", "0", "0", "0", "140")
	insertReadRepoTestJournalLine(t, db, sobId, firstPeriodId, accountId, dimensionOptionId, "40")

	page, err := repo.LedgersByPeriodRange(ctx, sobId, 2026, 2, 2026, 2, &dimensionOptionId, readRepoTestPageRequest(t))
	require.NoError(t, err)

	ledgers := page.Content()
	require.Len(t, ledgers, 1)
	require.Equal(t, accountId, ledgers[0].AccountId)
	requireDecimalEqual(t, "40", ledgers[0].OpeningAmount)
	requireDecimalEqual(t, "0", ledgers[0].PeriodDebit)
	requireDecimalEqual(t, "0", ledgers[0].PeriodCredit)
	requireDecimalEqual(t, "0", ledgers[0].PeriodAmount)
	requireDecimalEqual(t, "40", ledgers[0].EndingAmount)
}

func newReadRepoTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			NameReplacer: strings.NewReplacer("PO", ""),
		},
	})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(
		&accountPO{},
		&periodPO{},
		&ledgerPO{},
		&journalPO{},
		&journalLinePO{},
		&journalLineDimensionOptionPO{},
	))

	return db
}

func readRepoTestPageRequest(t *testing.T) data.PageRequest {
	t.Helper()

	p, err := pageable.New(1, 100)
	require.NoError(t, err)

	return data.NewPageRequest(p, sortable.Unsorted(), filterable.Unfiltered())
}

func insertReadRepoTestAccount(t *testing.T, db *gorm.DB, sobId, accountId uuid.UUID, rawAccountNumber string) {
	t.Helper()

	err := db.Create(&accountPO{
		Id:               accountId,
		SobId:            sobId,
		Title:            rawAccountNumber,
		RawAccountNumber: rawAccountNumber,
		Level:            1,
		IsLeaf:           true,
		BalanceDirection: "debit",
	}).Error
	require.NoError(t, err)
}

func insertReadRepoTestPeriod(t *testing.T, db *gorm.DB, sobId, periodId uuid.UUID, fiscalYear, periodNumber int) {
	t.Helper()

	err := db.Create(&periodPO{
		Id:           periodId,
		SobId:        sobId,
		FiscalYear:   fiscalYear,
		PeriodNumber: periodNumber,
	}).Error
	require.NoError(t, err)
}

func insertReadRepoTestLedger(
	t *testing.T,
	db *gorm.DB,
	sobId, accountId, periodId uuid.UUID,
	openingAmount, periodAmount, periodDebit, periodCredit, endingAmount string,
) {
	t.Helper()

	err := db.Create(&ledgerPO{
		Id:            uuid.New(),
		SobId:         sobId,
		AccountId:     accountId,
		PeriodId:      periodId,
		OpeningAmount: decimal.RequireFromString(openingAmount),
		PeriodAmount:  decimal.RequireFromString(periodAmount),
		PeriodDebit:   decimal.RequireFromString(periodDebit),
		PeriodCredit:  decimal.RequireFromString(periodCredit),
		EndingAmount:  decimal.RequireFromString(endingAmount),
	}).Error
	require.NoError(t, err)
}

func insertReadRepoTestJournalLine(
	t *testing.T,
	db *gorm.DB,
	sobId, periodId, accountId, dimensionOptionId uuid.UUID,
	amount string,
) {
	t.Helper()

	journalId := uuid.New()
	lineId := uuid.New()
	err := db.Create(&journalPO{
		Id:              journalId,
		SobId:           sobId,
		PeriodId:        periodId,
		DocumentNumber:  journalId.String(),
		IsPosted:        true,
		TransactionDate: time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
	}).Error
	require.NoError(t, err)

	err = db.Create(&journalLinePO{
		Id:        lineId,
		JournalId: journalId,
		AccountId: accountId,
		Text:      "tagged activity",
		Amount:    decimal.RequireFromString(amount),
	}).Error
	require.NoError(t, err)

	err = db.Create(&journalLineDimensionOptionPO{
		JournalLineId:     lineId,
		DimensionOptionId: dimensionOptionId,
	}).Error
	require.NoError(t, err)
}

func requireDecimalEqual(t *testing.T, expected string, actual decimal.Decimal) {
	t.Helper()
	require.True(t, decimal.RequireFromString(expected).Equal(actual), "expected %s, got %s", expected, actual.String())
}
