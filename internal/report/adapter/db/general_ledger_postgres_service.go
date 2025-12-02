package db

import (
	"context"
	"errors"

	"github/fims-proto/fims-proto-ms/internal/common/data/converter"
	"github/fims-proto/fims-proto-ms/internal/common/datasource"
	"github/fims-proto/fims-proto-ms/internal/common/utils"
	"github/fims-proto/fims-proto-ms/internal/report/domain/general_ledger"

	"github.com/google/uuid"
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

func (s GeneralLedgerPostgresService) ReadPeriodIdByFiscalYearAndNumber(
	ctx context.Context,
	sobId uuid.UUID,
	fiscalYear int,
	number int,
) (uuid.UUID, error) {
	db := s.dataSource.GetConnection(ctx)

	var po periodPO
	if err := db.Where(periodPO{SobId: sobId, FiscalYear: fiscalYear, PeriodNumber: number}).
		First(&po).Error; err != nil {
		return uuid.Nil, err
	}

	return po.Id, nil
}

func (s GeneralLedgerPostgresService) ReadPeriodById(
	ctx context.Context,
	_ uuid.UUID,
	periodId uuid.UUID,
) (*general_ledger.Period, error) {
	db := s.dataSource.GetConnection(ctx)

	po := periodPO{Id: periodId}
	if err := db.First(&po).Error; err != nil {
		return nil, err
	}

	return periodPOToBO(po), nil
}

func (s GeneralLedgerPostgresService) ReadFirstPeriodOfTheYear(
	ctx context.Context,
	sobId uuid.UUID,
	fiscalYear int,
) (*general_ledger.Period, error) {
	db := s.dataSource.GetConnection(ctx)

	var po periodPO
	err := db.Where(periodPO{SobId: sobId, FiscalYear: fiscalYear}).Order("period_number asc").First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// not found
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return periodPOToBO(po), nil
}

func (s GeneralLedgerPostgresService) ReadAccountIdsByNumbers(
	ctx context.Context,
	sobId uuid.UUID,
	accountNumbers []string,
) (map[string]uuid.UUID, error) {
	db := s.dataSource.GetConnection(ctx)

	// unique account numbers
	accountNumbers = utils.Unique(accountNumbers)

	var pos []accountPO
	if err := db.Where("sob_id = ? AND account_number IN ?", sobId, accountNumbers).Find(&pos).Error; err != nil {
		return nil, err
	}

	return utils.SliceToMap(
		pos, func(po accountPO) string {
			return po.AccountNumber
		}, func(po accountPO) uuid.UUID {
			return po.Id
		},
	), nil
}

func (s GeneralLedgerPostgresService) ReadLedgersByAccountAndPeriodsOrderByPeriod(
	ctx context.Context,
	sobId uuid.UUID,
	accountId uuid.UUID,
	periods []*general_ledger.Period,
) ([]*general_ledger.Ledger, error) {
	db := s.dataSource.GetConnection(ctx)

	var periodConditions [][]int
	for _, period := range periods {
		periodConditions = append(periodConditions, []int{period.FiscalYear(), period.PeriodNumber()})
	}

	var pos []ledgerPO

	if err := db.InnerJoins("Account", db.Where(accountPO{Id: accountId})).
		Joins("Period", db.Where("(fiscal_year, period_number) IN ?", periodConditions)).
		Where(ledgerPO{SobId: sobId}).
		Order("fiscal_year, period_number ASC").
		Find(&pos).Error; err != nil {
		return nil, err
	}

	return converter.POsToBOs(pos, ledgerPOToBO)
}
