package db

import (
	"context"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/common/data/filterable"
	"github/fims-proto/fims-proto-ms/internal/common/database"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/query"
)

type GeneralLedgerPostgresReadRepository struct{}

func NewGeneralLedgerPostgresReadRepository() *GeneralLedgerPostgresReadRepository {
	return &GeneralLedgerPostgresReadRepository{}
}

func (r GeneralLedgerPostgresReadRepository) SearchAccounts(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[query.Account], error) {
	addSobFilter(sobId, pageRequest)
	return data.SearchEntities(ctx, pageRequest, accountPO{}, accountPOToDTO, database.ReadDBFromContext(ctx).Preload("AuxiliaryCategories"))
}

func (r GeneralLedgerPostgresReadRepository) SearchAuxiliaryCategories(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[query.AuxiliaryCategory], error) {
	addSobFilter(sobId, pageRequest)
	return data.SearchEntities(ctx, pageRequest, auxiliaryCategoryPO{}, auxiliaryCategoryPOToDTO, database.ReadDBFromContext(ctx))
}

func (r GeneralLedgerPostgresReadRepository) SearchAuxiliaryAccounts(ctx context.Context, pageRequest data.PageRequest) (data.Page[query.AuxiliaryAccount], error) {
	return data.SearchEntities(ctx, pageRequest, auxiliaryAccountPO{}, auxiliaryAccountPOToDTO, database.ReadDBFromContext(ctx).InnerJoins("Category"))
}

func (r GeneralLedgerPostgresReadRepository) SearchLedgers(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[query.Ledger], error) {
	addSobFilter(sobId, pageRequest)
	return data.SearchEntities(ctx, pageRequest, ledgerPO{}, ledgerPOToDTO, database.ReadDBFromContext(ctx).InnerJoins("Account"))
}

func (r GeneralLedgerPostgresReadRepository) SearchAuxiliaryLedgers(ctx context.Context, pageRequest data.PageRequest) (data.Page[query.AuxiliaryLedger], error) {
	return data.SearchEntities(ctx, pageRequest, auxiliaryLedgerPO{}, auxiliaryLedgerPOToDTO, database.ReadDBFromContext(ctx).Joins("AuxiliaryAccount.Category"))
}

func (r GeneralLedgerPostgresReadRepository) SearchPeriods(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[query.Period], error) {
	addSobFilter(sobId, pageRequest)
	return data.SearchEntities(ctx, pageRequest, periodPO{}, periodPOToDTO, database.ReadDBFromContext(ctx))
}

func (r GeneralLedgerPostgresReadRepository) SearchVouchers(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[query.Voucher], error) {
	addSobFilter(sobId, pageRequest)
	return data.SearchEntities(ctx, pageRequest, voucherPO{}, voucherPOToDTO, database.ReadDBFromContext(ctx).Preload("LineItems.Account").Joins("Period"))
}

func (r GeneralLedgerPostgresReadRepository) PagingLedgersByPeriod(ctx context.Context, sobId, periodId uuid.UUID, pageRequest data.PageRequest) (data.Page[query.Ledger], error) {
	periodIdFilter, _ := filterable.NewFilter("periodId", filterable.OptEq, periodId)
	pageRequest.AddAndFilterable(filterable.NewFilterableAtom(periodIdFilter))
	return r.SearchLedgers(ctx, sobId, pageRequest)
}

func (r GeneralLedgerPostgresReadRepository) CurrentPeriod(ctx context.Context, sobId uuid.UUID) (query.Period, error) {
	db := database.ReadDBFromContext(ctx)

	var po periodPO
	if err := db.Where(periodPO{SobId: sobId, IsCurrent: true}).
		First(&po).Error; err != nil {
		return query.Period{}, err
	}

	return periodPOToDTO(po), nil
}

func (r GeneralLedgerPostgresReadRepository) VoucherById(ctx context.Context, voucherId uuid.UUID) (query.Voucher, error) {
	db := database.ReadDBFromContext(ctx)

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
