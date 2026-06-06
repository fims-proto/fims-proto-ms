package db

import (
	"context"
	"errors"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/datasource"
	"github/fims-proto/fims-proto-ms/internal/common/utils"
	"github/fims-proto/fims-proto-ms/internal/report/domain/service"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type GeneralLedgerPostgresService struct {
	dataSource datasource.DataSource
}

func NewGeneralLedgerPostgresService(dataSource datasource.DataSource) *GeneralLedgerPostgresService {
	if dataSource == nil {
		panic("nil data source")
	}
	return &GeneralLedgerPostgresService{dataSource: dataSource}
}

type ledgerPO struct {
	Id            uuid.UUID `gorm:"type:uuid;primaryKey"`
	SobId         uuid.UUID `gorm:"type:uuid"`
	AccountId     uuid.UUID `gorm:"type:uuid"`
	PeriodId      uuid.UUID `gorm:"type:uuid"`
	OpeningAmount decimal.Decimal
	PeriodAmount  decimal.Decimal
	PeriodDebit   decimal.Decimal
	PeriodCredit  decimal.Decimal
	EndingAmount  decimal.Decimal

	Account accountPO `gorm:"foreignKey:AccountId"`
	Period  periodPO  `gorm:"foreignKey:PeriodId"`
}

type accountPO struct {
	Id               uuid.UUID `gorm:"type:uuid;primaryKey"`
	SobId            uuid.UUID `gorm:"type:uuid"`
	RawAccountNumber string
	BalanceDirection string
}

type cashFlowItemPO struct {
	Id    uuid.UUID `gorm:"type:uuid;primaryKey"`
	SobId uuid.UUID `gorm:"type:uuid"`
	Code  string
	Name  string
}

func (s GeneralLedgerPostgresService) ReadPeriodIdByFiscalYearAndNumber(ctx context.Context, sobId uuid.UUID, fiscalYear int, number int) (uuid.UUID, error) {
	db := s.dataSource.GetConnection(ctx)
	var po periodPO
	if err := db.Where(periodPO{SobId: sobId, FiscalYear: fiscalYear, PeriodNumber: number}).First(&po).Error; err != nil {
		return uuid.Nil, err
	}
	return po.Id, nil
}

func (s GeneralLedgerPostgresService) ReadPeriodById(ctx context.Context, _ uuid.UUID, periodId uuid.UUID) (*service.Period, error) {
	db := s.dataSource.GetConnection(ctx)
	var po periodPO
	if err := db.Where("id = ?", periodId).First(&po).Error; err != nil {
		return nil, err
	}
	return periodPOToService(po), nil
}

func (s GeneralLedgerPostgresService) ReadFirstPeriodOfTheYear(ctx context.Context, sobId uuid.UUID, fiscalYear int) (*service.Period, error) {
	db := s.dataSource.GetConnection(ctx)
	var po periodPO
	err := db.Where(periodPO{SobId: sobId, FiscalYear: fiscalYear}).Order("period_number asc").First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return periodPOToService(po), nil
}

func (s GeneralLedgerPostgresService) ReadAccountIdsByRawNumbers(ctx context.Context, sobId uuid.UUID, rawAccountNumbers []string) (map[string]uuid.UUID, error) {
	db := s.dataSource.GetConnection(ctx)
	rawAccountNumbers = utils.Unique(rawAccountNumbers)
	var pos []accountPO
	if err := db.Where("sob_id = ? AND raw_account_number IN ?", sobId, rawAccountNumbers).Find(&pos).Error; err != nil {
		return nil, err
	}
	return utils.SliceToMap(
		pos,
		func(po accountPO) string { return po.RawAccountNumber },
		func(po accountPO) uuid.UUID { return po.Id },
	), nil
}

func (s GeneralLedgerPostgresService) ReadCashFlowItemIdsByCodes(ctx context.Context, sobId uuid.UUID, codes []string) (map[string]uuid.UUID, error) {
	db := s.dataSource.GetConnection(ctx)
	codes = utils.Unique(codes)
	var pos []cashFlowItemPO
	if err := db.Where("sob_id = ? AND code IN ?", sobId, codes).Find(&pos).Error; err != nil {
		return nil, err
	}
	return utils.SliceToMap(
		pos,
		func(po cashFlowItemPO) string { return po.Code },
		func(po cashFlowItemPO) uuid.UUID { return po.Id },
	), nil
}

func (s GeneralLedgerPostgresService) ReadLedgersByAccountAndPeriodsOrderByPeriod(
	ctx context.Context,
	sobId uuid.UUID,
	accountId uuid.UUID,
	periods []*service.Period,
) ([]service.Ledger, error) {
	db := s.dataSource.GetConnection(ctx)

	periodConditions := periodConditions(periods)
	var pos []ledgerPO
	if err := db.InnerJoins("Account", db.Where(accountPO{Id: accountId})).
		InnerJoins("Period", db.Where("(fiscal_year, period_number) IN ?", periodConditions)).
		Where(ledgerPO{SobId: sobId}).
		Order("fiscal_year, period_number ASC").
		Find(&pos).Error; err != nil {
		return nil, err
	}

	ledgers := make([]service.Ledger, 0, len(pos))
	for _, po := range pos {
		ledgers = append(ledgers, service.Ledger{
			AccountId:        po.AccountId,
			BalanceDirection: po.Account.BalanceDirection,
			Period:           *periodPOToService(po.Period),
			OpeningAmount:    po.OpeningAmount,
			PeriodAmount:     po.PeriodAmount,
			PeriodDebit:      po.PeriodDebit,
			PeriodCredit:     po.PeriodCredit,
			EndingAmount:     po.EndingAmount,
		})
	}
	return ledgers, nil
}

func (s GeneralLedgerPostgresService) SumAbsJournalLineAmountsByCashFlowItemsAndPeriods(
	ctx context.Context,
	sobId uuid.UUID,
	cashFlowItemIds []uuid.UUID,
	periods []*service.Period,
) (map[uuid.UUID]decimal.Decimal, error) {
	if len(cashFlowItemIds) == 0 || len(periods) == 0 {
		return map[uuid.UUID]decimal.Decimal{}, nil
	}

	db := s.dataSource.GetConnection(ctx)
	var rows []struct {
		CashFlowItemId uuid.UUID
		Amount         decimal.Decimal
	}

	if err := db.Table("journal_lines").
		Select("journal_lines.cash_flow_item_id AS cash_flow_item_id, COALESCE(SUM(ABS(journal_lines.amount)), 0) AS amount").
		Joins("JOIN journals ON journals.id = journal_lines.journal_id").
		Joins("JOIN periods ON periods.id = journals.period_id").
		Where("journals.sob_id = ? AND journals.is_posted = true", sobId).
		Where("journal_lines.cash_flow_item_id IN ?", cashFlowItemIds).
		Where("(periods.fiscal_year, periods.period_number) IN ?", periodConditions(periods)).
		Group("journal_lines.cash_flow_item_id").
		Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("failed to aggregate cash flow journal lines: %w", err)
	}

	result := make(map[uuid.UUID]decimal.Decimal, len(rows))
	for _, row := range rows {
		result[row.CashFlowItemId] = row.Amount
	}
	return result, nil
}

func periodConditions(periods []*service.Period) [][]int {
	var conditions [][]int
	for _, period := range periods {
		conditions = append(conditions, []int{period.FiscalYear, period.PeriodNumber})
	}
	return conditions
}

func periodPOToService(po periodPO) *service.Period {
	return &service.Period{
		Id:           po.Id,
		FiscalYear:   po.FiscalYear,
		PeriodNumber: po.PeriodNumber,
	}
}
