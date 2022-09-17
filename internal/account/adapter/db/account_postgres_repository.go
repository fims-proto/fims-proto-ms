package db

import (
	"context"
	"strings"
	"time"

	"github/fims-proto/fims-proto-ms/internal/common/datav3"
	"github/fims-proto/fims-proto-ms/internal/common/datav3/filterable"
	"github/fims-proto/fims-proto-ms/internal/common/datav3/pageable"
	"github/fims-proto/fims-proto-ms/internal/common/datav3/sortable"

	"github/fims-proto/fims-proto-ms/internal/account/domain/ledger"

	"github/fims-proto/fims-proto-ms/internal/account/domain/account"
	"github/fims-proto/fims-proto-ms/internal/account/domain/period"
	"gorm.io/gorm/clause"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/account/app/query"

	"gorm.io/gorm"
)

type AccountPostgresRepository struct{}

func NewAccountPostgresRepository() *AccountPostgresRepository {
	return &AccountPostgresRepository{}
}

func (r AccountPostgresRepository) Migrate(ctx context.Context) error {
	db := readDBFromCtx(ctx)

	return db.AutoMigrate(&accountPO{}, &periodPO{}, &ledgerPO{})
}

func (r AccountPostgresRepository) InitialAccounts(ctx context.Context, accounts []*account.Account) error {
	if len(accounts) == 0 {
		return errors.New("empty Account configuration list")
	}

	db := readDBFromCtx(ctx)

	return db.Transaction(func(tx *gorm.DB) error {
		// delete all within sob
		if err := tx.Where("sob_id = ?", accounts[0].SobId()).Delete(&accountPO{}).Error; err != nil {
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
		return tx.CreateInBatches(&accountPOs, 100).Error
	})
}

func (r AccountPostgresRepository) CreatePeriod(ctx context.Context, period *period.Period, txFn func() error) error {
	db := readDBFromCtx(ctx)

	po := periodBOToPO(*period)

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&po).Error; err != nil {
			return err
		}

		// other actions in current transaction
		return txFn()
	})
}

func (r AccountPostgresRepository) CreateLedgers(ctx context.Context, ledgers []*ledger.Ledger) error {
	db := readDBFromCtx(ctx)

	var ledgerPOs []ledgerPO
	for _, bo := range ledgers {
		ledgerPOs = append(ledgerPOs, ledgerBOToPO(*bo))
	}

	return db.Transaction(func(tx *gorm.DB) error {
		return tx.Omit("Account").CreateInBatches(&ledgerPOs, 500).Error
	})
}

func (r AccountPostgresRepository) UpdateLedgersByPeriodAndAccountIds(ctx context.Context, periodId uuid.UUID, accountIds []uuid.UUID, updateFn func(ledgers []*ledger.Ledger) ([]*ledger.Ledger, error)) error {
	db := readDBFromCtx(ctx)

	return db.Transaction(func(tx *gorm.DB) error {
		var ledgerPOs []ledgerPO
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
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
			return errors.Wrap(err, "failed to update ledgers in transaction")
		}

		var updatedPOs []ledgerPO
		for _, updatedLedger := range updatedLedgers {
			updatedPOs = append(updatedPOs, ledgerBOToPO(*updatedLedger))
		}

		return tx.Save(&updatedPOs).Error
	})
}

// queries

func (r AccountPostgresRepository) SearchAccounts(ctx context.Context, sobId uuid.UUID, pageRequest datav3.PageRequest) (datav3.Page[query.Account], error) {
	addSobFilter(sobId, "account.sobId", pageRequest)
	return datav3.SearchEntities(pageRequest, &accountPO{}, accountPOToDTO, resolveEntity, readDBFromCtx(ctx))
}

func (r AccountPostgresRepository) SearchPeriods(ctx context.Context, sobId uuid.UUID, pageRequest datav3.PageRequest) (datav3.Page[query.Period], error) {
	addSobFilter(sobId, "period.sobId", pageRequest)
	return datav3.SearchEntities(pageRequest, &periodPO{}, periodPOToDTO, resolveEntity, readDBFromCtx(ctx))
}

func (r AccountPostgresRepository) SearchLedgers(ctx context.Context, sobId uuid.UUID, pageRequest datav3.PageRequest) (datav3.Page[query.Ledger], error) {
	addSobFilter(sobId, "ledger.sobId", pageRequest)
	return datav3.SearchEntities(pageRequest, &ledgerPO{}, ledgerPOToDTO, resolveEntity, readDBFromCtx(ctx).Preload("Account"))
}

func (r AccountPostgresRepository) PagingLedgersByPeriod(ctx context.Context, sobId uuid.UUID, periodId uuid.UUID, pageRequest datav3.PageRequest) (datav3.Page[query.Ledger], error) {
	periodIdFilter, _ := filterable.NewFilter("ledger.periodId", "eq", periodId)
	pageRequest.AddFilter(periodIdFilter)
	return r.SearchLedgers(ctx, sobId, pageRequest)
}

func (r AccountPostgresRepository) LedgersInPeriod(ctx context.Context, sobId uuid.UUID, periodId uuid.UUID) ([]query.Ledger, error) {
	periodIdFilter, _ := filterable.NewFilter("ledger.periodId", "eq", periodId)
	pageRequest := datav3.NewPageRequest(
		pageable.Unpaged(),
		sortable.Unsorted(),
		filterable.New(periodIdFilter),
	)
	ledgers, err := r.SearchLedgers(ctx, sobId, pageRequest)
	if err != nil {
		return nil, err
	}

	return ledgers.Content(), nil
}

func (r AccountPostgresRepository) AllAccounts(ctx context.Context, sobId uuid.UUID) ([]query.Account, error) {
	accounts, err := r.SearchAccounts(
		ctx,
		sobId,
		datav3.NewPageRequest(pageable.Unpaged(), sortable.Unsorted(), filterable.Unfiltered()),
	)
	if err != nil {
		return nil, err
	}
	return accounts.Content(), nil
}

func (r AccountPostgresRepository) SuperiorAccounts(ctx context.Context, accountId uuid.UUID) ([]query.Account, error) {
	rawSql := `WITH RECURSIVE res AS (
		   SELECT *
		   FROM a_accounts
		   WHERE account_id = ?
		   UNION
		   SELECT a_accounts.*
		   FROM res
		   JOIN a_accounts ON a_accounts.account_id = res.superior_account_id
		)
		SELECT *
		FROM res
		WHERE account_id != ?`
	rawSql = strings.Join(strings.Fields(rawSql), " ") // normalize whitespaces

	db := readDBFromCtx(ctx)

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

func (r AccountPostgresRepository) AccountsByIds(ctx context.Context, accountIds []uuid.UUID) ([]query.Account, error) {
	accountIdFilter, _ := filterable.NewFilter("accountId", "in", accountIds)
	pageRequest := datav3.NewPageRequest(
		pageable.Unpaged(),
		sortable.Unsorted(),
		filterable.New(accountIdFilter),
	)
	accounts, err := r.SearchAccounts(ctx, uuid.Nil, pageRequest)
	if err != nil {
		return nil, err
	}
	return accounts.Content(), nil
}

func (r AccountPostgresRepository) AccountsByNumbers(ctx context.Context, sobId uuid.UUID, accountNumbers []string) ([]query.Account, error) {
	accountIdFilter, _ := filterable.NewFilter("accountNumber", "in", accountNumbers)
	pageRequest := datav3.NewPageRequest(
		pageable.Unpaged(),
		sortable.Unsorted(),
		filterable.New(accountIdFilter),
	)
	accounts, err := r.SearchAccounts(ctx, sobId, pageRequest)
	if err != nil {
		return nil, err
	}
	return accounts.Content(), nil
}

func (r AccountPostgresRepository) PeriodById(ctx context.Context, periodId uuid.UUID) (query.Period, error) {
	periodIdFilter, _ := filterable.NewFilter("period.id", "eq", periodId)
	pageRequest := datav3.NewPageRequest(
		pageable.Unpaged(),
		sortable.Unsorted(),
		filterable.New(periodIdFilter),
	)
	periods, err := r.SearchPeriods(ctx, uuid.Nil, pageRequest)
	if err != nil {
		return query.Period{}, err
	}
	if periods.NumberOfElements() != 1 {
		return query.Period{}, errors.Errorf("period not found by id: %s", periodId)
	}

	return periods.Content()[0], nil
}

func (r AccountPostgresRepository) PeriodsByIds(ctx context.Context, periodIds []uuid.UUID) ([]query.Period, error) {
	periodIdFilter, _ := filterable.NewFilter("period.id", "in", periodIds)
	pageRequest := datav3.NewPageRequest(
		pageable.Unpaged(),
		sortable.Unsorted(),
		filterable.New(periodIdFilter),
	)
	periods, err := r.SearchPeriods(ctx, uuid.Nil, pageRequest)
	if err != nil {
		return nil, err
	}

	return periods.Content(), nil
}

func (r AccountPostgresRepository) PeriodByTime(ctx context.Context, sobId uuid.UUID, timePoint time.Time) (query.Period, error) {
	db := readDBFromCtx(ctx)

	var periodPOs []periodPO
	if err := db.
		Where("sob_id = ? AND opening_time <= ? AND (ending_time > ? OR ending_time = ?)", sobId, timePoint, timePoint, time.Time{}).
		Find(&periodPOs).Error; err != nil {
		return query.Period{}, errors.Wrap(err, "find period by id failed")
	}

	if len(periodPOs) != 1 {
		return query.Period{}, errors.Errorf("expected 1 but %d periods found", len(periodPOs))
	}

	return periodPOToDTO(periodPOs[0])
}

func addSobFilter(sobId uuid.UUID, fieldName string, pageRequest datav3.PageRequest) {
	if sobId != uuid.Nil {
		sobIdFilter, _ := filterable.NewFilter(fieldName, "eq", sobId.String())
		pageRequest.AddFilter(sobIdFilter)
	}
}

func readDBFromCtx(ctx context.Context) *gorm.DB {
	return ctx.Value("db").(*gorm.DB)
}

func resolveEntity(entity string) (string, error) {
	switch entity {
	case "account":
		return accountPO{}.TableName(), nil
	case "period":
		return periodPO{}.TableName(), nil
	case "ledger":
		return ledgerPO{}.TableName(), nil
	default:
		return "", errors.Errorf("invalid entity name %s", entity)
	}
}
