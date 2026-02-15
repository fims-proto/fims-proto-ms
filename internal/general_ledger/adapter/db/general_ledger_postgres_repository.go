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
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_account"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_category"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_ledger"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/ledger"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/period"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/voucher"

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
		&auxiliaryCategoryPO{},
		&auxiliaryAccountPO{},
		&periodPO{},
		&ledgerPO{},
		&auxiliaryLedgerPO{},
		&voucherPO{},
		&lineItemPO{},
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
	if err := db.Omit("AuxiliaryCategories").CreateInBatches(&pos, 100).Error; err != nil {
		return err
	}

	// save associations
	for _, po := range pos {
		if len(po.AuxiliaryCategories) > 0 {
			if err := db.Model(&po).Omit("AuxiliaryCategories.*").Association("AuxiliaryCategories").Replace(po.AuxiliaryCategories); err != nil {
				return fmt.Errorf("failed to save auxiliary category associations for account %s: %w", po.AccountNumber, err)
			}
		}
	}

	return nil
}

func (r GeneralLedgerPostgresRepository) CreateAccount(ctx context.Context, a *account.Account) error {
	db := r.dataSource.GetConnection(ctx)

	po := accountBOToPO(a)

	if err := db.Omit("AuxiliaryCategories").Create(&po).Error; err != nil {
		return err
	}

	// save associations
	if len(po.AuxiliaryCategories) > 0 {
		if err := db.Model(&po).Omit("AuxiliaryCategories.*").Association("AuxiliaryCategories").Replace(po.AuxiliaryCategories); err != nil {
			return fmt.Errorf("failed to save auxiliary category associations: %w", err)
		}
	}

	return nil
}

func (r GeneralLedgerPostgresRepository) UpdateAccount(
	ctx context.Context,
	accountId uuid.UUID,
	updateFn func(a *account.Account) (*account.Account, error),
) error {
	db := r.dataSource.GetConnection(ctx)

	po := accountPO{Id: accountId}
	if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("AuxiliaryCategories").
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

	if err = db.Omit("AuxiliaryCategories").Save(&po).Error; err != nil {
		return err
	}

	if err = db.Model(&po).Omit("AuxiliaryCategories.*").Association("AuxiliaryCategories").Replace(po.AuxiliaryCategories); err != nil {
		return fmt.Errorf("failed to update auxiliary category associations: %w", err)
	}

	return nil
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
		Preload("AuxiliaryCategories").Find(&pos).Error; err != nil {
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
	err := db.Order("opening_time asc").Where(periodPO{SobId: sobId}).First(&po).Error
	if err == nil {
		return periodPOToBO(po)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, commonErrors.ErrRecordNotFound()
	}
	return nil, err
}

func (r GeneralLedgerPostgresRepository) CreateLedgers(ctx context.Context, ledgers []*ledger.Ledger) error {
	db := r.dataSource.GetConnection(ctx)

	pos := converter.BOsToPOs(ledgers, ledgerBOToPO)

	return db.Omit("Account").CreateInBatches(&pos, 100).Error
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

	updatedPOs := converter.BOsToPOs(updatedLedgers, ledgerBOToPO)

	return db.Omit("Account").Save(&updatedPOs).Error
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
		Where("ending_debit_balance <> ending_credit_balance").
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

func (r GeneralLedgerPostgresRepository) CreateVoucher(ctx context.Context, v *voucher.Voucher) error {
	db := r.dataSource.GetConnection(ctx)

	po := voucherBOToPO(*v)

	return db.Create(&po).Error
}

func (r GeneralLedgerPostgresRepository) UpdateVoucher(
	ctx context.Context,
	voucherId uuid.UUID,
	updateFn func(v *voucher.Voucher) (*voucher.Voucher, error),
) error {
	db := r.dataSource.GetConnection(ctx)

	po := voucherPO{Id: voucherId}
	if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("LineItems.Account.AuxiliaryCategories").
		Preload("LineItems.AuxiliaryAccounts.Category").
		Preload("Period").
		First(&po).Error; err != nil {
		return err
	}

	// remove existing link between line item and auxiliary account
	for _, item := range po.LineItems {
		if err := db.Model(&item).Association("AuxiliaryAccounts").Clear(); err != nil {
			return fmt.Errorf("failed to remove line item auxiliary account associations: %w", err)
		}
	}

	bo, err := voucherPOToBO(po)
	if err != nil {
		return fmt.Errorf("failed to update voucher: %w", err)
	}

	updatedBO, err := updateFn(bo)
	if err != nil {
		return fmt.Errorf("failed to update voucher: %w", err)
	}

	po = voucherBOToPO(*updatedBO)

	// remove existing line items
	if err = db.Where("voucher_id = ?", po.Id).Delete(&lineItemPO{}).Error; err != nil {
		return fmt.Errorf("failed to delete line items: %w", err)
	}

	return db.Save(&po).Error
}

func (r GeneralLedgerPostgresRepository) ExistsVouchersNotPostedInPeriod(ctx context.Context, sobId, periodId uuid.UUID) (bool, error) {
	db := r.dataSource.GetConnection(ctx)

	var count int64
	err := db.Model(&voucherPO{}).
		Where("sob_id = ? AND period_id = ? AND is_posted = false", sobId, periodId).
		Count(&count).
		Error

	return count > 0, err
}

func (r GeneralLedgerPostgresRepository) CreateAuxiliaryCategories(ctx context.Context, categories []*auxiliary_category.AuxiliaryCategory) error {
	db := r.dataSource.GetConnection(ctx)

	pos := converter.BOsToPOs(categories, auxiliaryCategoryBOToPO)

	return db.CreateInBatches(&pos, 100).Error
}

func (r GeneralLedgerPostgresRepository) ReadAuxiliaryCategoryByKey(ctx context.Context, sobId uuid.UUID, key string) (*auxiliary_category.AuxiliaryCategory, error) {
	db := r.dataSource.GetConnection(ctx)

	var po auxiliaryCategoryPO
	if err := db.Where(auxiliaryCategoryPO{SobId: sobId, Key: key}).First(&po).Error; err != nil {
		return nil, err
	}

	return auxiliaryCategoryPOToBO(po)
}

func (r GeneralLedgerPostgresRepository) ReadAuxiliaryCategoriesByKeys(ctx context.Context, sobId uuid.UUID, keys []string) ([]*auxiliary_category.AuxiliaryCategory, error) {
	db := r.dataSource.GetConnection(ctx)

	// unique keys
	keys = utils.Unique(keys)

	if len(keys) == 0 {
		return nil, nil
	}

	var pos []auxiliaryCategoryPO
	if err := db.Where("sob_id = ? AND key IN ?", sobId, keys).Find(&pos).Error; err != nil {
		return nil, err
	}

	if len(pos) != len(keys) {
		return nil, fmt.Errorf("not all auxiliary categories found for sob %s and keys %v", sobId, keys)
	}

	// check if all keys are found
	for _, key := range keys {
		found := false
		for _, po := range pos {
			if po.Key == key {
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("auxiliary category with key %s not found for sob %s", key, sobId)
		}
	}

	return converter.POsToBOs(pos, auxiliaryCategoryPOToBO)
}

func (r GeneralLedgerPostgresRepository) CreateAuxiliaryAccounts(ctx context.Context, accounts []*auxiliary_account.AuxiliaryAccount) error {
	db := r.dataSource.GetConnection(ctx)

	pos := converter.BOsToPOs(accounts, auxiliaryAccountBOToPO)

	return db.Omit("Category").CreateInBatches(&pos, 100).Error
}

func (r GeneralLedgerPostgresRepository) ReadAuxiliaryAccountsByPairs(
	ctx context.Context,
	sobId uuid.UUID,
	pairs []auxiliary_account.AuxiliaryPair,
) ([]*auxiliary_account.AuxiliaryAccount, error) {
	db := r.dataSource.GetConnection(ctx)

	if len(pairs) == 0 {
		return nil, nil
	}

	dbOr := db.Session(&gorm.Session{NewDB: true})
	for _, pair := range pairs {
		dbOr = dbOr.Or(`"Category"."key" = ? AND "a_auxiliary_accounts"."key" = ?`, pair.CategoryKey, pair.AccountKey)
	}

	var auxiliaryAccountPOs []auxiliaryAccountPO
	if err := db.InnerJoins("Category", db.Where(&auxiliaryCategoryPO{SobId: sobId})).
		Where(dbOr).
		Find(&auxiliaryAccountPOs).Error; err != nil {
		return nil, err
	}

	return converter.POsToBOs(auxiliaryAccountPOs, auxiliaryAccountPOToBO)
}

func (r GeneralLedgerPostgresRepository) ReadAllAuxiliaryAccounts(ctx context.Context, sobId uuid.UUID) (
	[]*auxiliary_account.AuxiliaryAccount,
	error,
) {
	db := r.dataSource.GetConnection(ctx)

	var auxiliaryAccountPOs []auxiliaryAccountPO
	if err := db.InnerJoins("Category", db.Where(&auxiliaryCategoryPO{SobId: sobId})).Find(&auxiliaryAccountPOs).Error; err != nil {
		return nil, err
	}

	return converter.POsToBOs(auxiliaryAccountPOs, auxiliaryAccountPOToBO)
}

func (r GeneralLedgerPostgresRepository) CreateAuxiliaryLedgers(ctx context.Context, ledgers []*auxiliary_ledger.AuxiliaryLedger) error {
	db := r.dataSource.GetConnection(ctx)

	pos := converter.BOsToPOs(ledgers, auxiliaryLedgerBOToPO)

	return db.Omit("AuxiliaryAccount", "AuxiliaryCategory", "Account").CreateInBatches(&pos, 100).Error
}

func (r GeneralLedgerPostgresRepository) UpdateAuxiliaryLedgersByPeriodAndAccounts(
	ctx context.Context,
	periodId uuid.UUID,
	accountId uuid.UUID,
	auxiliaryCategoryIds []uuid.UUID,
	auxiliaryAccountIds []uuid.UUID,
	updateFn func(auxiliaryLedgers []*auxiliary_ledger.AuxiliaryLedger) ([]*auxiliary_ledger.AuxiliaryLedger, error),
) error {
	db := r.dataSource.GetConnection(ctx)

	// unique IDs
	auxiliaryCategoryIds = utils.Unique(auxiliaryCategoryIds)
	auxiliaryAccountIds = utils.Unique(auxiliaryAccountIds)

	var auxiliaryLedgerPOs []auxiliaryLedgerPO
	if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("period_id = ? AND account_id = ? AND auxiliary_category_id IN ? AND auxiliary_account_id IN ?",
			periodId, accountId, auxiliaryCategoryIds, auxiliaryAccountIds).
		Preload("AuxiliaryAccount.Category").
		Preload("AuxiliaryCategory").
		Preload("Account").
		Find(&auxiliaryLedgerPOs).Error; err != nil {
		return err
	}

	auxiliaryLedgerBOs, err := converter.POsToBOs(auxiliaryLedgerPOs, auxiliaryLedgerPOToBO)
	if err != nil {
		return fmt.Errorf("failed to convert auxiliary ledgers: %w", err)
	}

	updated, err := updateFn(auxiliaryLedgerBOs)
	if err != nil {
		return fmt.Errorf("failed to update auxiliary ledgers: %w", err)
	}

	updatedPOs := converter.BOsToPOs(updated, auxiliaryLedgerBOToPO)

	return db.Omit("AuxiliaryAccount", "AuxiliaryCategory", "Account").Save(&updatedPOs).Error
}

func (r GeneralLedgerPostgresRepository) ReadAuxiliaryLedgersByPeriod(ctx context.Context, periodId uuid.UUID) (
	[]*auxiliary_ledger.AuxiliaryLedger,
	error,
) {
	db := r.dataSource.GetConnection(ctx)

	var auxiliaryLedgerPOs []auxiliaryLedgerPO
	if err := db.
		Where(auxiliaryLedgerPO{PeriodId: periodId}).
		InnerJoins("AuxiliaryAccount").
		InnerJoins("AuxiliaryAccount.Category").
		InnerJoins("AuxiliaryCategory").
		InnerJoins("Account").
		Find(&auxiliaryLedgerPOs).
		Error; err != nil {
		return nil, err
	}

	return converter.POsToBOs(auxiliaryLedgerPOs, auxiliaryLedgerPOToBO)
}

func (r GeneralLedgerPostgresRepository) ReadAuxiliaryLedgersByAccountAndPeriod(
	ctx context.Context,
	accountId uuid.UUID,
	periodId uuid.UUID,
) ([]*auxiliary_ledger.AuxiliaryLedger, error) {
	db := r.dataSource.GetConnection(ctx)

	var auxiliaryLedgerPOs []auxiliaryLedgerPO
	if err := db.
		Where("account_id = ? AND period_id = ?", accountId, periodId).
		InnerJoins("AuxiliaryAccount").
		InnerJoins("AuxiliaryAccount.Category").
		InnerJoins("AuxiliaryCategory").
		InnerJoins("Account").
		Find(&auxiliaryLedgerPOs).
		Error; err != nil {
		return nil, err
	}

	return converter.POsToBOs(auxiliaryLedgerPOs, auxiliaryLedgerPOToBO)
}
