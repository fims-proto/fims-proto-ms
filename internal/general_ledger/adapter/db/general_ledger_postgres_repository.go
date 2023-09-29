package db

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/common/database"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account"
	accountType "github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account/account_type"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_account"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_category"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_ledger"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/ledger"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/period"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/voucher"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GeneralLedgerPostgresRepository struct{}

func NewGeneralLedgerPostgresRepository() *GeneralLedgerPostgresRepository {
	return &GeneralLedgerPostgresRepository{}
}

func (r GeneralLedgerPostgresRepository) Migrate(ctx context.Context) error {
	db := database.ReadDBFromContext(ctx)

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
	db := database.ReadDBFromContext(ctx)

	return db.Transaction(func(tx *gorm.DB) error {
		return txFn(database.NewContextWithDB(ctx, tx))
	})
}

func (r GeneralLedgerPostgresRepository) InitialAccounts(ctx context.Context, accounts []*account.Account) error {
	if len(accounts) == 0 {
		return errors.New("empty Account list")
	}

	db := database.ReadDBFromContext(ctx)

	// delete all within sob
	if err := db.Where("sob_id = ?", accounts[0].SobId()).Delete(&accountPO{}).Error; err != nil {
		return fmt.Errorf("failed initialize accounts: %w", err)
	}

	// create all
	pos := bos2pos(accounts, accountBOToPO)
	return db.Omit("AuxiliaryCategories").CreateInBatches(&pos, 100).Error
}

func (r GeneralLedgerPostgresRepository) UpdateAccount(ctx context.Context, accountId uuid.UUID, updateFn func(a *account.Account) (*account.Account, error)) error {
	db := database.ReadDBFromContext(ctx)

	po := accountPO{Id: accountId}
	if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).Preload("AuxiliaryCategories").First(&po).Error; err != nil {
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

	po = accountBOToPO(*updatedBO)

	if err = db.Omit("AuxiliaryCategories").Save(&po).Error; err != nil {
		return err
	}

	if err = db.Model(&po).Omit("AuxiliaryCategories.*").Association("AuxiliaryCategories").Replace(po.AuxiliaryCategories); err != nil {
		return fmt.Errorf("failed to update auxiliary category associations: %w", err)
	}

	return nil
}

func (r GeneralLedgerPostgresRepository) ReadAllAccounts(ctx context.Context, sobId uuid.UUID) ([]*account.Account, error) {
	db := database.ReadDBFromContext(ctx)

	var accountPOs []accountPO
	if err := db.Where(accountPO{SobId: sobId}).Find(&accountPOs).Error; err != nil {
		return nil, err
	}

	return pos2bos(accountPOs, accountPOToBO)
}

func (r GeneralLedgerPostgresRepository) ReadAccountsByNumbers(ctx context.Context, sobId uuid.UUID, accountNumbers []string) ([]*account.Account, error) {
	db := database.ReadDBFromContext(ctx)

	if len(accountNumbers) == 0 {
		return nil, nil
	}

	var accountPOs []accountPO
	if err := db.Where("sob_id = ? AND account_number IN ?", sobId, accountNumbers).Preload("AuxiliaryCategories").Find(&accountPOs).Error; err != nil {
		return nil, err
	}

	return pos2bos(accountPOs, accountPOToBO)
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

	db := database.ReadDBFromContext(ctx)

	var accountPOs []accountPO
	if err := db.Raw(rawSql, accountId, accountId).Scan(&accountPOs).Error; err != nil {
		return nil, err
	}

	return pos2bos(accountPOs, accountPOToBO)
}

func (r GeneralLedgerPostgresRepository) CreatePeriodIfNotExists(ctx context.Context, p *period.Period) (*period.Period, bool, error) {
	db := database.ReadDBFromContext(ctx)

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

func (r GeneralLedgerPostgresRepository) UpdatePeriod(ctx context.Context, periodId uuid.UUID, updateFn func(p *period.Period) (*period.Period, error)) error {
	db := database.ReadDBFromContext(ctx)

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
	db := database.ReadDBFromContext(ctx)

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
	db := database.ReadDBFromContext(ctx)

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

func (r GeneralLedgerPostgresRepository) CreateLedgers(ctx context.Context, ledgers []*ledger.Ledger) error {
	db := database.ReadDBFromContext(ctx)

	pos := bos2pos(ledgers, ledgerBOToPO)

	return db.Omit("Account").CreateInBatches(&pos, 100).Error
}

func (r GeneralLedgerPostgresRepository) UpdateLedgersByPeriodAndAccountIds(ctx context.Context, periodId uuid.UUID, accountIds []uuid.UUID, updateFn func(accounts []*ledger.Ledger) ([]*ledger.Ledger, error)) error {
	db := database.ReadDBFromContext(ctx)

	var ledgerPOs []ledgerPO
	if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("period_id = ? AND account_id IN ?", periodId, accountIds).
		Preload("Account").
		Find(&ledgerPOs).Error; err != nil {
		return err
	}

	ledgerBOs, err := pos2bos(ledgerPOs, ledgerPOToBO)
	if err != nil {
		return fmt.Errorf("failed to update ledgers: %w", err)
	}

	updatedLedgers, err := updateFn(ledgerBOs)
	if err != nil {
		return fmt.Errorf("failed to update ledgers: %w", err)
	}

	updatedPOs := bos2pos(updatedLedgers, ledgerBOToPO)

	return db.Omit("Account").Save(&updatedPOs).Error
}

func (r GeneralLedgerPostgresRepository) ReadLedgersByPeriod(ctx context.Context, periodId uuid.UUID) ([]*ledger.Ledger, error) {
	db := database.ReadDBFromContext(ctx)

	var ledgerPOs []ledgerPO
	if err := db.Where(ledgerPO{PeriodId: periodId}).Preload("Account").Find(&ledgerPOs).Error; err != nil {
		return nil, err
	}

	return pos2bos(ledgerPOs, ledgerPOToBO)
}

func (r GeneralLedgerPostgresRepository) ExistsProfitAndLossLedgersHavingBalanceInPeriod(ctx context.Context, sobId, periodId uuid.UUID) (bool, error) {
	db := database.ReadDBFromContext(ctx)

	var count int64
	err := db.Model(&ledgerPO{}).
		Where(ledgerPO{SobId: sobId, PeriodId: periodId}).
		Where("ending_balance <> '0'").
		Joins("Account", db.Where(accountPO{AccountType: accountType.ProfitAndLoss.String()})).
		Count(&count).
		Error

	return count > 0, err
}

func (r GeneralLedgerPostgresRepository) ReadFirstLevelLedgersInPeriod(ctx context.Context, sobId, periodId uuid.UUID) ([]*ledger.Ledger, error) {
	db := database.ReadDBFromContext(ctx)

	var ledgerPOs []ledgerPO
	if err := db.Where(ledgerPO{SobId: sobId, PeriodId: periodId}).
		Joins("Account", db.Where(accountPO{Level: 1})).
		Find(&ledgerPOs).Error; err != nil {
		return nil, err
	}

	return pos2bos(ledgerPOs, ledgerPOToBO)
}

func (r GeneralLedgerPostgresRepository) CreateVoucher(ctx context.Context, v *voucher.Voucher) error {
	db := database.ReadDBFromContext(ctx)

	po := voucherBOToPO(*v)

	return db.Create(&po).Error
}

func (r GeneralLedgerPostgresRepository) UpdateVoucher(ctx context.Context, voucherId uuid.UUID, updateFn func(v *voucher.Voucher) (*voucher.Voucher, error)) error {
	db := database.ReadDBFromContext(ctx)

	po := voucherPO{Id: voucherId}
	if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("LineItems.Account.AuxiliaryCategories").
		Preload("LineItems.AuxiliaryAccounts.Category").
		First(&po).Error; err != nil {
		return err
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

	return db.Save(&po).Error
}

func (r GeneralLedgerPostgresRepository) ExistsVouchersNotPostedInPeriod(ctx context.Context, sobId, periodId uuid.UUID) (bool, error) {
	db := database.ReadDBFromContext(ctx)

	var count int64
	err := db.Model(&voucherPO{}).
		Where("sob_id = ? AND period_id = ? AND is_posted = false", sobId, periodId).
		Count(&count).
		Error

	return count > 0, err
}

func (r GeneralLedgerPostgresRepository) CreateAuxiliaryCategories(ctx context.Context, categories []*auxiliary_category.AuxiliaryCategory) error {
	db := database.ReadDBFromContext(ctx)

	pos := bos2pos(categories, auxiliaryCategoryBOToPO)

	return db.CreateInBatches(&pos, 100).Error
}

func (r GeneralLedgerPostgresRepository) ReadAuxiliaryCategoryByKey(ctx context.Context, key string) (*auxiliary_category.AuxiliaryCategory, error) {
	db := database.ReadDBFromContext(ctx)

	var po auxiliaryCategoryPO
	if err := db.Where(auxiliaryCategoryPO{Key: key}).First(&po).Error; err != nil {
		return nil, err
	}

	return auxiliaryCategoryPOToBO(po)
}

func (r GeneralLedgerPostgresRepository) CreateAuxiliaryAccounts(ctx context.Context, accounts []*auxiliary_account.AuxiliaryAccount) error {
	db := database.ReadDBFromContext(ctx)

	pos := bos2pos(accounts, auxiliaryAccountBOToPO)

	return db.Omit("Category").CreateInBatches(&pos, 100).Error
}

func (r GeneralLedgerPostgresRepository) ReadAuxiliaryAccountsByPairs(ctx context.Context, sobId uuid.UUID, pairs []auxiliary_account.AuxiliaryPair) ([]*auxiliary_account.AuxiliaryAccount, error) {
	db := database.ReadDBFromContext(ctx)

	if len(pairs) == 0 {
		return nil, nil
	}

	dbOr := db.Session(&gorm.Session{NewDB: true})
	for _, pair := range pairs {
		dbOr = dbOr.Or(`"Category"."key" = ? AND "a_auxiliary_accounts"."key" = ?`, pair.CategoryKey, pair.AccountKey)
	}

	var auxiliaryAccountPOs []auxiliaryAccountPO
	if err := db.InnerJoins("Category", db.Where(&auxiliaryCategoryPO{SobId: sobId})).Where(dbOr).Find(&auxiliaryAccountPOs).Error; err != nil {
		return nil, err
	}

	return pos2bos(auxiliaryAccountPOs, auxiliaryAccountPOToBO)
}

func (r GeneralLedgerPostgresRepository) ReadAllAuxiliaryAccounts(ctx context.Context, sobId uuid.UUID) ([]*auxiliary_account.AuxiliaryAccount, error) {
	db := database.ReadDBFromContext(ctx)

	var auxiliaryAccountPOs []auxiliaryAccountPO
	if err := db.InnerJoins("Category", db.Where(&auxiliaryCategoryPO{SobId: sobId})).Find(&auxiliaryAccountPOs).Error; err != nil {
		return nil, err
	}

	return pos2bos(auxiliaryAccountPOs, auxiliaryAccountPOToBO)
}

func (r GeneralLedgerPostgresRepository) CreateAuxiliaryLedgers(ctx context.Context, ledgers []*auxiliary_ledger.AuxiliaryLedger) error {
	db := database.ReadDBFromContext(ctx)

	pos := bos2pos(ledgers, auxiliaryLedgerBOToPO)

	return db.Omit("AuxiliaryAccount").CreateInBatches(&pos, 100).Error
}

func (r GeneralLedgerPostgresRepository) UpdateAuxiliaryLedgersByPeriodAndAccountIds(ctx context.Context, periodId uuid.UUID, auxiliaryAccountIds []uuid.UUID, updateFn func(auxiliaryLedgers []*auxiliary_ledger.AuxiliaryLedger) ([]*auxiliary_ledger.AuxiliaryLedger, error)) error {
	db := database.ReadDBFromContext(ctx)

	var auxiliaryLedgerPOs []auxiliaryLedgerPO
	if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("period_id = ? AND auxiliary_account_id IN ?", periodId, auxiliaryAccountIds).
		Preload("AuxiliaryAccount.Category").
		Find(&auxiliaryLedgerPOs).Error; err != nil {
		return err
	}

	auxiliaryLedgerBOs, err := pos2bos(auxiliaryLedgerPOs, auxiliaryLedgerPOToBO)
	if err != nil {
		return fmt.Errorf("failed to update auxiliary ledgers: %w", err)
	}

	updated, err := updateFn(auxiliaryLedgerBOs)
	if err != nil {
		return fmt.Errorf("failed to update auxiliary ledgers: %w", err)
	}

	updatedPOs := bos2pos(updated, auxiliaryLedgerBOToPO)

	return db.Omit("AuxiliaryAccount").Save(&updatedPOs).Error
}

func (r GeneralLedgerPostgresRepository) ReadAuxiliaryLedgersByPeriod(ctx context.Context, periodId uuid.UUID) ([]*auxiliary_ledger.AuxiliaryLedger, error) {
	db := database.ReadDBFromContext(ctx)

	var auxiliaryLedgerPOs []auxiliaryLedgerPO
	if err := db.Where(auxiliaryLedgerPO{PeriodId: periodId}).Preload("AuxiliaryAccount").Find(&auxiliaryLedgerPOs).Error; err != nil {
		return nil, err
	}

	return pos2bos(auxiliaryLedgerPOs, auxiliaryLedgerPOToBO)
}
