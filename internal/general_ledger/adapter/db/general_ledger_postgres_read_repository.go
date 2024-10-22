package db

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/common/data/filterable"
	"github/fims-proto/fims-proto-ms/internal/common/datasource"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/query"
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
	return data.SearchEntities(ctx, pageRequest, accountPO{}, accountPOToDTO, r.dataSource.GetConnection(ctx).Preload("AuxiliaryCategories"))
}

func (r GeneralLedgerPostgresReadRepository) SearchAuxiliaryCategories(
	ctx context.Context,
	sobId uuid.UUID,
	pageRequest data.PageRequest,
) (data.Page[query.AuxiliaryCategory], error) {
	addSobFilter(sobId, pageRequest)
	return data.SearchEntities(ctx, pageRequest, auxiliaryCategoryPO{}, auxiliaryCategoryPOToDTO, r.dataSource.GetConnection(ctx))
}

func (r GeneralLedgerPostgresReadRepository) SearchAuxiliaryAccounts(
	ctx context.Context,
	pageRequest data.PageRequest,
) (data.Page[query.AuxiliaryAccount], error) {
	return data.SearchEntities(ctx, pageRequest, auxiliaryAccountPO{}, auxiliaryAccountPOToDTO, r.dataSource.GetConnection(ctx).InnerJoins("Category"))
}

func (r GeneralLedgerPostgresReadRepository) SearchLedgers(
	ctx context.Context,
	sobId uuid.UUID,
	pageRequest data.PageRequest,
) (data.Page[query.Ledger], error) {
	addSobFilter(sobId, pageRequest)
	return data.SearchEntities(ctx, pageRequest, ledgerPO{}, ledgerPOToDTO, r.dataSource.GetConnection(ctx).InnerJoins("Account"))
}

func (r GeneralLedgerPostgresReadRepository) SearchAuxiliaryLedgers(
	ctx context.Context,
	pageRequest data.PageRequest,
) (data.Page[query.AuxiliaryLedger], error) {
	return data.SearchEntities(ctx, pageRequest, auxiliaryLedgerPO{}, auxiliaryLedgerPOToDTO, r.dataSource.GetConnection(ctx).Joins("AuxiliaryAccount.Category"))
}

func (r GeneralLedgerPostgresReadRepository) SearchPeriods(
	ctx context.Context,
	sobId uuid.UUID,
	pageRequest data.PageRequest,
) (data.Page[query.Period], error) {
	addSobFilter(sobId, pageRequest)
	return data.SearchEntities(ctx, pageRequest, periodPO{}, periodPOToDTO, r.dataSource.GetConnection(ctx))
}

func (r GeneralLedgerPostgresReadRepository) SearchVouchers(
	ctx context.Context,
	sobId uuid.UUID,
	pageRequest data.PageRequest,
) (data.Page[query.Voucher], error) {
	addSobFilter(sobId, pageRequest)
	return data.SearchEntities(ctx, pageRequest, voucherPO{}, voucherPOToDTO, r.dataSource.GetConnection(ctx).Preload("LineItems.Account").Joins("Period"))
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
	err := db.Order("opening_time asc").Where(periodPO{SobId: sobId}).First(&po).Error
	if err == nil {
		return periodPOToDTO(po), nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return query.Period{}, commonErrors.ErrRecordNotFound()
	}
	return query.Period{}, err
}

func (r GeneralLedgerPostgresReadRepository) VoucherById(ctx context.Context, voucherId uuid.UUID) (query.Voucher, error) {
	db := r.dataSource.GetConnection(ctx)

	po := voucherPO{Id: voucherId}
	if err := db.
		Preload("LineItems.Account.AuxiliaryCategories").
		Preload("LineItems.AuxiliaryAccounts.Category").
		Preload("Period").
		First(&po).Error; err != nil {
		return query.Voucher{}, err
	}

	return voucherPOToDTO(po), nil
}

func addSobFilter(sobId uuid.UUID, pageRequest data.PageRequest) {
	if sobId != uuid.Nil {
		sobIdFilter, _ := filterable.NewFilter("sobId", filterable.OptEq, sobId.String())
		pageRequest.AddAndFilterable(filterable.NewFilterableAtom(sobIdFilter))
	}
}
