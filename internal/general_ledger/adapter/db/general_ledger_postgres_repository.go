package db

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github/fims-proto/fims-proto-ms/internal/common/data/converter"
	"github/fims-proto/fims-proto-ms/internal/common/utils"

	"github/fims-proto/fims-proto-ms/internal/common/datasource"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account/class"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/journal"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/ledger"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/ledger_entry"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/period"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GeneralLedgerPostgresRepository struct {
	dataSource datasource.DataSource
}

func NewGeneralLedgerPostgresRepository(dataSource datasource.DataSource) *GeneralLedgerPostgresRepository {
	if dataSource == nil {
		panic("nil data source")
	}

	return &GeneralLedgerPostgresRepository{
		dataSource: dataSource,
	}
}

func (r GeneralLedgerPostgresRepository) Migrate(ctx context.Context) error {
	db := r.dataSource.GetConnection(ctx)

	return db.AutoMigrate(
		&accountPO{},
		&periodPO{},
		&ledgerEntryPO{},
		&ledgerPO{},
		&journalPO{},
		&journalLinePO{},
	)
}

func (r GeneralLedgerPostgresRepository) EnableTx(ctx context.Context, txFn func(txCtx context.Context) error) error {
	return r.dataSource.EnableTransaction(ctx, txFn)
}

func (r GeneralLedgerPostgresRepository) InitialAccounts(ctx context.Context, accounts []*account.Account) error {
	if len(accounts) == 0 {
		return errors.New("empty Account list")
	}

	db := r.dataSource.GetConnection(ctx)

	// delete all within sob
	if err := db.Where("sob_id = ?", accounts[0].SobId()).Delete(&accountPO{}).Error; err != nil {
		return fmt.Errorf("failed initialize accounts: %w", err)
	}

	// create all
	pos := converter.BOsToPOs(accounts, accountBOToPO)
	return db.CreateInBatches(&pos, 100).Error
}

func (r GeneralLedgerPostgresRepository) CreateAccount(ctx context.Context, a *account.Account) error {
	db := r.dataSource.GetConnection(ctx)

	po := accountBOToPO(a)
	return db.Create(&po).Error
}

func (r GeneralLedgerPostgresRepository) UpdateAccount(
	ctx context.Context,
	accountId uuid.UUID,
	updateFn func(a *account.Account) (*account.Account, error),
) error {
	db := r.dataSource.GetConnection(ctx)

	po := accountPO{Id: accountId}
	if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&po).Error; err != nil {
		return err
	}

	bo, err := accountPOToBO(po)
	if err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	updatedBO, err := updateFn(bo)
	if err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	po = accountBOToPO(updatedBO)

	return db.Save(&po).Error
}

func (r GeneralLedgerPostgresRepository) ReadAllAccounts(ctx context.Context, sobId uuid.UUID) ([]*account.Account, error) {
	db := r.dataSource.GetConnection(ctx)

	var accountPOs []accountPO
	if err := db.Where(accountPO{SobId: sobId}).Find(&accountPOs).Error; err != nil {
		return nil, err
	}

	return converter.POsToBOs(accountPOs, accountPOToBO)
}

func (r GeneralLedgerPostgresRepository) ReadAccountByNumber(ctx context.Context, sobId uuid.UUID, accountNumber string) (*account.Account, error) {
	db := r.dataSource.GetConnection(ctx)

	var po accountPO
	if err := db.Where(accountPO{SobId: sobId, AccountNumber: accountNumber}).First(&po).Error; err != nil {
		return nil, err
	}

	return accountPOToBO(po)
}

func (r GeneralLedgerPostgresRepository) ReadAccountsByNumbers(ctx context.Context, sobId uuid.UUID, accountNumbers []string) ([]*account.Account, error) {
	db := r.dataSource.GetConnection(ctx)

	// unique account numbers
	accountNumbers = utils.Unique(accountNumbers)

	if len(accountNumbers) == 0 {
		return nil, nil
	}

	var pos []accountPO
	if err := db.Where("sob_id = ? AND account_number IN ?", sobId, accountNumbers).
		Find(&pos).Error; err != nil {
		return nil, err
	}

	if len(pos) != len(accountNumbers) {
		return nil, fmt.Errorf("not all accounts found for sob %s and account numbers %v", sobId, accountNumbers)
	}

	// check if all keys are found
	for _, accountNumber := range accountNumbers {
		found := false
		for _, po := range pos {
			if po.AccountNumber == accountNumber {
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("account with number %s not found for sob %s", accountNumber, sobId)
		}
	}

	return converter.POsToBOs(pos, accountPOToBO)
}

func (r GeneralLedgerPostgresRepository) ReadSuperiorAccountsById(ctx context.Context, accountId uuid.UUID) ([]*account.Account, error) {
	rawSql := `WITH RECURSIVE res AS (
		   SELECT *
		   FROM a_accounts
		   WHERE id = ?
		   UNION
		   SELECT a_accounts.*
		   FROM res
		   JOIN a_accounts ON a_accounts.id = res.superior_account_id
		)
		SELECT *
		FROM res
		WHERE id != ?`
	rawSql = strings.Join(strings.Fields(rawSql), " ") // normalize whitespaces

	db := r.dataSource.GetConnection(ctx)

	var accountPOs []accountPO
	if err := db.Raw(rawSql, accountId, accountId).Scan(&accountPOs).Error; err != nil {
		return nil, err
	}

	return converter.POsToBOs(accountPOs, accountPOToBO)
}

func (r GeneralLedgerPostgresRepository) ReadAccountsWithSuperiorsByIds(
	ctx context.Context,
	sobId uuid.UUID,
	accountIds []uuid.UUID,
) ([]*account.Account, error) {
	db := r.dataSource.GetConnection(ctx)

	// unique account ids
	accountIds = utils.Unique(accountIds)

	var accountPOs []accountPO
	if err := db.Where("sob_id = ? AND id IN ?", sobId, accountIds).Find(&accountPOs).Error; err != nil {
		return nil, err
	}

	// check if superiors exist, if yes get superiors first
	var superiorIds []uuid.UUID
	for _, po := range accountPOs {
		if po.SuperiorAccountId != nil {
			superiorIds = append(superiorIds, *po.SuperiorAccountId)
		}
	}
	var superiorAccountBOs []*account.Account
	if len(superiorIds) > 0 {
		var err error
		superiorAccountBOs, err = r.ReadAccountsWithSuperiorsByIds(ctx, sobId, superiorIds)
		if err != nil {
			return nil, err
		}
	}
	superiorAccountMap := utils.SliceToMap(superiorAccountBOs, func(e *account.Account) uuid.UUID {
		return e.Id()
	}, func(e *account.Account) *account.Account {
		return e
	})

	// convert
	var result []*account.Account
	for _, po := range accountPOs {
		var sa *account.Account
		if po.SuperiorAccountId != nil {
			var ok bool
			sa, ok = superiorAccountMap[*po.SuperiorAccountId]
			if !ok {
				return nil, fmt.Errorf("superior account not found for %s", po.AccountNumber)
			}
		}
		a, err := accountPOToBOWithSuperior(po, sa)
		if err != nil {
			return nil, err
		}
		result = append(result, a)
	}

	return result, nil
}

func (r GeneralLedgerPostgresRepository) ReadAllSubAccountsWithSuperiors(ctx context.Context, sobId uuid.UUID) ([]*account.Account, error) {
	db := r.dataSource.GetConnection(ctx)

	// leaf accounts
	var subAccountPOs []accountPO
	if err := db.Select("id").Where(accountPO{SobId: sobId, IsLeaf: true}).Find(&subAccountPOs).Error; err != nil {
		return nil, err
	}
	var subAccountIds []uuid.UUID
	for _, po := range subAccountPOs {
		subAccountIds = append(subAccountIds, po.Id)
	}

	return r.ReadAccountsWithSuperiorsByIds(ctx, sobId, subAccountIds)
}

func (r GeneralLedgerPostgresRepository) CreatePeriodIfNotExists(ctx context.Context, p *period.Period) (*period.Period, bool, error) {
	db := r.dataSource.GetConnection(ctx)

	var existedPeriod periodPO
	err := db.Where(periodPO{SobId: p.SobId(), FiscalYear: p.FiscalYear(), PeriodNumber: p.PeriodNumber()}).First(&existedPeriod).Error

	if err == nil {
		// found
		bo, err := periodPOToBO(existedPeriod)
		return bo, false, err
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// error
		return nil, false, err
	}

	po := periodBOToPO(*p)

	if po.IsCurrent {
		// make sure only 1 current period in one sob
		_, err = r.ReadCurrentPeriod(ctx, po.SobId)
		if err == nil {
			return nil, false, commonErrors.NewSlugError("period-duplicatedCurrent")
		} else if !errors.Is(err, commonErrors.ErrRecordNotFound()) {
			return nil, false, fmt.Errorf("failed to check current period: %w", err)
		}
	}

	return p, true, db.Save(&po).Error
}

func (r GeneralLedgerPostgresRepository) UpdatePeriod(
	ctx context.Context,
	periodId uuid.UUID,
	updateFn func(p *period.Period) (*period.Period, error),
) error {
	db := r.dataSource.GetConnection(ctx)

	po := periodPO{Id: periodId}
	if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&po).Error; err != nil {
		return err
	}

	bo, err := periodPOToBO(po)
	if err != nil {
		return fmt.Errorf("failed to update period: %w", err)
	}

	updatedBO, err := updateFn(bo)
	if err != nil {
		return fmt.Errorf("failed to update period: %w", err)
	}

	po = periodBOToPO(*updatedBO)

	return db.Save(&po).Error
}

func (r GeneralLedgerPostgresRepository) ReadCurrentPeriod(ctx context.Context, sobId uuid.UUID) (*period.Period, error) {
	db := r.dataSource.GetConnection(ctx)

	var po periodPO
	err := db.Where(periodPO{SobId: sobId, IsCurrent: true}).Take(&po).Error
	if err == nil {
		return periodPOToBO(po)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, commonErrors.ErrRecordNotFound()
	}
	return nil, err
}

func (r GeneralLedgerPostgresRepository) ReadPreviousPeriod(ctx context.Context, currentPeriodId uuid.UUID) (*period.Period, error) {
	db := r.dataSource.GetConnection(ctx)

	currentPO := periodPO{Id: currentPeriodId}
	if err := db.First(&currentPO).Error; err != nil {
		return nil, fmt.Errorf("failed to find period by id %s: %w", currentPeriodId, err)
	}

	currentBO, err := periodPOToBO(currentPO)
	if err != nil {
		return nil, fmt.Errorf("failed to read period: %w", err)
	}

	previousFiscalYear, previousNumber := currentBO.PreviousNumber()

	var previousPO periodPO
	err = db.Where(periodPO{SobId: currentBO.SobId(), FiscalYear: previousFiscalYear, PeriodNumber: previousNumber}).First(&previousPO).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, commonErrors.ErrRecordNotFound()
	} else if err != nil {
		return nil, err
	}

	return periodPOToBO(previousPO)
}

func (r GeneralLedgerPostgresRepository) ReadFirstPeriod(ctx context.Context, sobId uuid.UUID) (*period.Period, error) {
	db := r.dataSource.GetConnection(ctx)

	var po periodPO
	err := db.Order("fiscal_year asc, period_number asc").Where(periodPO{SobId: sobId}).First(&po).Error
	if err == nil {
		return periodPOToBO(po)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, commonErrors.ErrRecordNotFound()
	}
	return nil, err
}

func (r GeneralLedgerPostgresRepository) CreateLedgerEntries(ctx context.Context, entries []*ledger_entry.LedgerEntry) error {
	db := r.dataSource.GetConnection(ctx)

	return db.CreateInBatches(new(converter.BOsToPOs(entries, ledgerEntryBOToPOForCreate)), 100).Error
}

func (r GeneralLedgerPostgresRepository) CreateLedgers(ctx context.Context, ledgers []*ledger.Ledger) error {
	db := r.dataSource.GetConnection(ctx)

	return db.Omit("Account").CreateInBatches(new(converter.BOsToPOs(ledgers, ledgerBOToPO)), 100).Error
}

func (r GeneralLedgerPostgresRepository) UpdateLedgersByPeriodAndAccountIds(
	ctx context.Context,
	periodId uuid.UUID,
	accountIds []uuid.UUID,
	updateFn func(accounts []*ledger.Ledger) ([]*ledger.Ledger, error),
) error {
	db := r.dataSource.GetConnection(ctx)

	// unique account ids
	accountIds = utils.Unique(accountIds)

	var ledgerPOs []ledgerPO
	if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("period_id = ? AND account_id IN ?", periodId, accountIds).
		Preload("Account").
		Find(&ledgerPOs).Error; err != nil {
		return err
	}

	ledgerBOs, err := converter.POsToBOs(ledgerPOs, ledgerPOToBO)
	if err != nil {
		return fmt.Errorf("failed to update ledgers: %w", err)
	}

	updatedLedgers, err := updateFn(ledgerBOs)
	if err != nil {
		return fmt.Errorf("failed to update ledgers: %w", err)
	}

	return db.Omit("Account").Save(new(converter.BOsToPOs(updatedLedgers, ledgerBOToPO))).Error
}

func (r GeneralLedgerPostgresRepository) ReadLedgersByPeriod(ctx context.Context, periodId uuid.UUID) ([]*ledger.Ledger, error) {
	db := r.dataSource.GetConnection(ctx)

	var ledgerPOs []ledgerPO
	if err := db.Where(ledgerPO{PeriodId: periodId}).Preload("Account").Find(&ledgerPOs).Error; err != nil {
		return nil, err
	}

	return converter.POsToBOs(ledgerPOs, ledgerPOToBO)
}

func (r GeneralLedgerPostgresRepository) ExistsProfitAndLossLedgersHavingBalanceInPeriod(
	ctx context.Context,
	sobId, periodId uuid.UUID,
) (bool, error) {
	db := r.dataSource.GetConnection(ctx)

	var count int64
	err := db.Model(&ledgerPO{}).
		Where(ledgerPO{SobId: sobId, PeriodId: periodId}).
		Where("ending_amount <> 0").
		InnerJoins("Account", db.Where(accountPO{Class: int(class.ProfitsAndLosses)})).
		Count(&count).
		Error

	return count > 0, err
}

func (r GeneralLedgerPostgresRepository) ReadFirstLevelLedgersInPeriod(ctx context.Context, sobId, periodId uuid.UUID) ([]*ledger.Ledger, error) {
	db := r.dataSource.GetConnection(ctx)

	var ledgerPOs []ledgerPO
	if err := db.Where(ledgerPO{SobId: sobId, PeriodId: periodId}).
		InnerJoins("Account", db.Where(accountPO{Level: 1})).
		Find(&ledgerPOs).Error; err != nil {
		return nil, err
	}

	return converter.POsToBOs(ledgerPOs, ledgerPOToBO)
}

func (r GeneralLedgerPostgresRepository) CreateJournal(ctx context.Context, j *journal.Journal) error {
	db := r.dataSource.GetConnection(ctx)

	return db.Create(new(journalBOToPO(*j))).Error
}

func (r GeneralLedgerPostgresRepository) UpdateJournalHeader(
	ctx context.Context,
	journalId uuid.UUID,
	updateFn func(j *journal.Journal) (*journal.Journal, error),
) error {
	db := r.dataSource.GetConnection(ctx)

	po := journalPO{Id: journalId}
	if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("JournalLines.Account").
		Preload("Period").
		First(&po).Error; err != nil {
		return err
	}

	bo, err := journalPOToBO(po)
	if err != nil {
		return fmt.Errorf("failed to update journal header: %w", err)
	}

	updatedBO, err := updateFn(bo)
	if err != nil {
		return fmt.Errorf("failed to update journal header: %w", err)
	}

	updatedPO := journalBOToPO(*updatedBO)

	// Only persist the journal table columns; omit associations to avoid touching journal_lines.
	return db.Omit("JournalLines", "Period").Save(&updatedPO).Error
}

func (r GeneralLedgerPostgresRepository) UpdateEntireJournal(
	ctx context.Context,
	journalId uuid.UUID,
	updateFn func(j *journal.Journal) (*journal.Journal, error),
) error {
	db := r.dataSource.GetConnection(ctx)

	po := journalPO{Id: journalId}
	if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("JournalLines.Account").
		Preload("Period").
		First(&po).Error; err != nil {
		return err
	}

	bo, err := journalPOToBO(po)
	if err != nil {
		return fmt.Errorf("failed to update journal: %w", err)
	}

	updatedBO, err := updateFn(bo)
	if err != nil {
		return fmt.Errorf("failed to update journal: %w", err)
	}

	po = journalBOToPO(*updatedBO)

	// remove existing journal lines
	if err = db.Where("journal_id = ?", po.Id).Delete(&journalLinePO{}).Error; err != nil {
		return fmt.Errorf("failed to delete journal lines: %w", err)
	}

	return db.Save(&po).Error
}

func (r GeneralLedgerPostgresRepository) ExistsJournalsNotPostedInPeriod(ctx context.Context, sobId, periodId uuid.UUID) (bool, error) {
	db := r.dataSource.GetConnection(ctx)

	var count int64
	err := db.Model(&journalPO{}).
		Where("sob_id = ? AND period_id = ? AND is_posted = false", sobId, periodId).
		Count(&count).
		Error

	return count > 0, err
}
