package db

import (
	"context"
	"strings"

	accountType "github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account/account_type"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/common/data/filterable"
	"github/fims-proto/fims-proto-ms/internal/common/database"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/query"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account"
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
		&periodPO{},
		&ledgerPO{},
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
		return errors.Wrap(err, "accounts deletion failed")
	}

	// create all
	var accountPOs []accountPO
	for _, accountBO := range accounts {
		po, err := accountBOToPO(*accountBO)
		if err != nil {
			return errors.Wrap(err, "failed to map Account from BO to PO")
		}
		accountPOs = append(accountPOs, po)
	}
	return db.CreateInBatches(&accountPOs, 100).Error
}

func (r GeneralLedgerPostgresRepository) CreatePeriod(ctx context.Context, bo *period.Period) error {
	db := database.ReadDBFromContext(ctx)

	po := periodBOToPO(*bo)

	if po.IsCurrent {
		// make sure only 1 current period in one sob
		_, err := r.currentPeriod(db, po.SobId)
		if err == nil {
			return commonErrors.NewSlugError("period-duplicatedCurrent")
		} else if _, ok := err.(commonErrors.ObjectNotFoundErr); !ok {
			return errors.Wrap(err, "failed to check current period")
		}
	}

	return db.Create(&po).Error
}

func (r GeneralLedgerPostgresRepository) UpdatePeriod(ctx context.Context, periodId uuid.UUID, updateFn func(p *period.Period) (*period.Period, error)) error {
	db := database.ReadDBFromContext(ctx)

	po := periodPO{Id: periodId}
	if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&po).Error; err != nil {
		return err
	}

	bo, err := periodPOToBO(po)
	if err != nil {
		return errors.Wrap(err, "failed to map period")
	}

	updatedBO, err := updateFn(bo)
	if err != nil {
		return errors.Wrap(err, "update period failed")
	}

	po = periodBOToPO(*updatedBO)

	return db.Save(&po).Error
}

func (r GeneralLedgerPostgresRepository) CreateLedgers(ctx context.Context, ledgers []*ledger.Ledger) error {
	db := database.ReadDBFromContext(ctx)

	var ledgerPOs []ledgerPO
	for _, bo := range ledgers {
		ledgerPOs = append(ledgerPOs, ledgerBOToPO(*bo))
	}

	return db.Omit("Account").CreateInBatches(&ledgerPOs, 500).Error
}

func (r GeneralLedgerPostgresRepository) UpdateLedgersByPeriodAndAccountIds(ctx context.Context, periodId uuid.UUID, accountIds []uuid.UUID, updateFn func(ledgers []*ledger.Ledger) ([]*ledger.Ledger, error)) error {
	db := database.ReadDBFromContext(ctx)

	var ledgerPOs []ledgerPO
	if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("period_id = ? AND account_id IN ?", periodId, accountIds).
		Preload("Account").
		Find(&ledgerPOs).Error; err != nil {
		return err
	}

	var ledgerBOs []*ledger.Ledger
	for _, po := range ledgerPOs {
		bo, err := ledgerPOToBO(po)
		if err != nil {
			return errors.Wrap(err, "failed to map ledger")
		}
		ledgerBOs = append(ledgerBOs, bo)
	}

	updatedLedgers, err := updateFn(ledgerBOs)
	if err != nil {
		return errors.Wrap(err, "failed to update ledgers")
	}

	var updatedPOs []ledgerPO
	for _, updatedLedger := range updatedLedgers {
		updatedPOs = append(updatedPOs, ledgerBOToPO(*updatedLedger))
	}

	return db.Save(&updatedPOs).Error
}

func (r GeneralLedgerPostgresRepository) CreateVoucher(ctx context.Context, d *voucher.Voucher) error {
	db := database.ReadDBFromContext(ctx)

	po := voucherBOToPO(*d)

	return db.Create(&po).Error
}

func (r GeneralLedgerPostgresRepository) UpdateVoucher(ctx context.Context, voucherId uuid.UUID, updateFn func(d *voucher.Voucher) (*voucher.Voucher, error)) error {
	db := database.ReadDBFromContext(ctx)

	po := voucherPO{Id: voucherId}
	if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).Preload("LineItems").First(&po).Error; err != nil {
		return err
	}

	bo, err := voucherPOToBO(po)
	if err != nil {
		return errors.Wrap(err, "failed to map voucher")
	}

	updatedBO, err := updateFn(bo)
	if err != nil {
		return errors.Wrap(err, "update voucher failed")
	}

	po = voucherBOToPO(*updatedBO)

	// remove existing first
	if err = db.Where("voucher_id = ?", po.Id).Delete(&lineItemPO{}).Error; err != nil {
		return errors.Wrap(err, "delete voucher items failed")
	}

	return db.Save(&po).Error
}

// queries

func (r GeneralLedgerPostgresRepository) SearchAccounts(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[query.Account], error) {
	addSobFilter(sobId, pageRequest)
	return data.SearchEntities(ctx, pageRequest, accountPO{}, accountPOToDTO, database.ReadDBFromContext(ctx))
}

func (r GeneralLedgerPostgresRepository) SearchPeriods(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[query.Period], error) {
	addSobFilter(sobId, pageRequest)
	return data.SearchEntities(ctx, pageRequest, periodPO{}, periodPOToDTO, database.ReadDBFromContext(ctx))
}

func (r GeneralLedgerPostgresRepository) SearchLedgers(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[query.Ledger], error) {
	addSobFilter(sobId, pageRequest)
	return data.SearchEntities(ctx, pageRequest, ledgerPO{}, ledgerPOToDTO, database.ReadDBFromContext(ctx).Joins("Account"))
}

func (r GeneralLedgerPostgresRepository) SearchVouchers(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[query.Voucher], error) {
	addSobFilter(sobId, pageRequest)
	return data.SearchEntities(ctx, pageRequest, voucherPO{}, voucherPOToDTO, database.ReadDBFromContext(ctx).Preload("LineItems.Account").Joins("Period"))
}

func addSobFilter(sobId uuid.UUID, pageRequest data.PageRequest) {
	if sobId != uuid.Nil {
		sobIdFilter, _ := filterable.NewFilter("sobId", "eq", sobId.String())
		pageRequest.AddFilter(sobIdFilter)
	}
}

func (r GeneralLedgerPostgresRepository) PagingLedgersByPeriod(ctx context.Context, sobId, periodId uuid.UUID, pageRequest data.PageRequest) (data.Page[query.Ledger], error) {
	periodIdFilter, _ := filterable.NewFilter("periodId", "eq", periodId)
	pageRequest.AddFilter(periodIdFilter)
	return r.SearchLedgers(ctx, sobId, pageRequest)
}

func (r GeneralLedgerPostgresRepository) LedgersInPeriod(ctx context.Context, sobId, periodId uuid.UUID) ([]query.Ledger, error) {
	db := database.ReadDBFromContext(ctx)

	var ledgerPOs []ledgerPO
	if err := db.Where(ledgerPO{SobId: sobId, PeriodId: periodId}).Preload("Account").Find(&ledgerPOs).Error; err != nil {
		return nil, err
	}

	return toDTO(ledgerPOs, ledgerPOToDTO)
}

func (r GeneralLedgerPostgresRepository) FirstLevelLedgersInPeriod(ctx context.Context, sobId, periodId uuid.UUID) ([]query.Ledger, error) {
	db := database.ReadDBFromContext(ctx)

	var ledgerPOs []ledgerPO
	if err := db.Where(ledgerPO{SobId: sobId, PeriodId: periodId}).
		Joins("Account", db.Where(accountPO{Level: 1})).
		Find(&ledgerPOs).Error; err != nil {
		return nil, err
	}

	return toDTO(ledgerPOs, ledgerPOToDTO)
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

func (r GeneralLedgerPostgresRepository) AllAccounts(ctx context.Context, sobId uuid.UUID) ([]query.Account, error) {
	db := database.ReadDBFromContext(ctx)

	var accountPOs []accountPO
	if err := db.Where(accountPO{SobId: sobId}).Find(&accountPOs).Error; err != nil {
		return nil, err
	}

	return toDTO(accountPOs, accountPOToDTO)
}

func (r GeneralLedgerPostgresRepository) SuperiorAccounts(ctx context.Context, accountId uuid.UUID) ([]query.Account, error) {
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

	var accountDTOs []query.Account
	for _, po := range accountPOs {
		dto, err := accountPOToDTO(po)
		if err != nil {
			return nil, errors.Wrap(err, "failed to map account")
		}

		accountDTOs = append(accountDTOs, dto)
	}

	return accountDTOs, nil
}

func (r GeneralLedgerPostgresRepository) AccountsByIds(ctx context.Context, accountIds []uuid.UUID) ([]query.Account, error) {
	db := database.ReadDBFromContext(ctx)

	var accountPOs []accountPO
	if err := db.Where(accountIds).Find(&accountPOs).Error; err != nil {
		return nil, err
	}

	return toDTO(accountPOs, accountPOToDTO)
}

func (r GeneralLedgerPostgresRepository) AccountsByNumbers(ctx context.Context, sobId uuid.UUID, accountNumbers []string) ([]query.Account, error) {
	db := database.ReadDBFromContext(ctx)

	var accountPOs []accountPO
	if err := db.Where("sob_id = ? AND account_number IN ?", sobId, accountNumbers).Find(&accountPOs).Error; err != nil {
		return nil, err
	}

	return toDTO(accountPOs, accountPOToDTO)
}

func (r GeneralLedgerPostgresRepository) CurrentPeriod(ctx context.Context, sobId uuid.UUID) (query.Period, error) {
	db := database.ReadDBFromContext(ctx)

	currentPeriod, err := r.currentPeriod(db, sobId)
	if err != nil {
		return query.Period{}, err
	}

	return periodPOToDTO(currentPeriod)
}

func (r GeneralLedgerPostgresRepository) PeriodById(ctx context.Context, periodId uuid.UUID) (query.Period, error) {
	db := database.ReadDBFromContext(ctx)

	po := periodPO{Id: periodId}
	err := db.First(&po).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return query.Period{}, commonErrors.NewObjectNotFoundErr("period")
	} else if err != nil {
		return query.Period{}, err
	}

	return periodPOToDTO(po)
}

func (r GeneralLedgerPostgresRepository) PeriodsByIds(ctx context.Context, periodIds []uuid.UUID) ([]query.Period, error) {
	db := database.ReadDBFromContext(ctx)

	var periodPOs []periodPO
	if err := db.Where(periodIds).Find(&periodPOs).Error; err != nil {
		return nil, err
	}

	return toDTO(periodPOs, periodPOToDTO)
}

func (r GeneralLedgerPostgresRepository) PeriodByFiscalYearAndNumber(ctx context.Context, sobId uuid.UUID, fiscalYear, periodNumber int) (query.Period, error) {
	db := database.ReadDBFromContext(ctx)

	var po periodPO
	err := db.Where(periodPO{SobId: sobId, FiscalYear: fiscalYear, PeriodNumber: periodNumber}).First(&po).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return query.Period{}, commonErrors.NewObjectNotFoundErr("period")
	} else if err != nil {
		return query.Period{}, err
	}

	return periodPOToDTO(po)
}

func (r GeneralLedgerPostgresRepository) VoucherById(ctx context.Context, voucherId uuid.UUID) (query.Voucher, error) {
	db := database.ReadDBFromContext(ctx)

	po := voucherPO{Id: voucherId}
	err := db.Preload("LineItems.Account").Preload("Period").First(&po).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return query.Voucher{}, commonErrors.NewObjectNotFoundErr("voucher")
	} else if err != nil {
		return query.Voucher{}, err
	}

	return voucherPOToDTO(po)
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

func (r GeneralLedgerPostgresRepository) currentPeriod(db *gorm.DB, sobId uuid.UUID) (periodPO, error) {
	var periods []periodPO
	if err := db.Where(periodPO{SobId: sobId, IsCurrent: true}).
		Find(&periods).Error; err != nil {
		return periodPO{}, err
	}

	if len(periods) == 0 {
		return periodPO{}, commonErrors.NewObjectNotFoundErr("period")
	} else if len(periods) > 1 {
		return periodPO{}, errors.Errorf("expected 1 but %d periods found", len(periods))
	}

	return periods[0], nil
}

func toDTO[PO any, DTO any](pos []PO, convert func(po PO) (DTO, error)) ([]DTO, error) {
	var dtos []DTO
	for _, po := range pos {
		dto, err := convert(po)
		if err != nil {
			return nil, errors.Wrap(err, "failed to map entity to DTO")
		}
		dtos = append(dtos, dto)
	}
	return dtos, nil
}
