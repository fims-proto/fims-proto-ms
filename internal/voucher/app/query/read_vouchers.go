package query

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type void struct{}

var empty void

func toKeySlice[K comparable, V interface{}](set map[K]V) []K {
	keys := make([]K, len(set))
	i := 0
	for k := range set {
		keys[i] = k
		i++
	}
	return keys
}

type VouchersReadModel interface {
	ReadAllVouchers(ctx context.Context, sobId uuid.UUID, pageable data.Pageable) (data.Page[Voucher], error)
	ReadById(ctx context.Context, id uuid.UUID) (Voucher, error)
}

type ReadVouchersHandler struct {
	readModel      VouchersReadModel
	accountService AccountService
}

func NewReadVouchersHandler(readModel VouchersReadModel, accountService AccountService) ReadVouchersHandler {
	if readModel == nil {
		panic("nil readModel")
	}
	if accountService == nil {
		panic("nil account service")
	}
	return ReadVouchersHandler{
		readModel:      readModel,
		accountService: accountService,
	}
}

func (h ReadVouchersHandler) HandleReadAll(ctx context.Context, sobId uuid.UUID, pageable data.Pageable) (data.Page[Voucher], error) {
	vouchersPage, err := h.readModel.ReadAllVouchers(ctx, sobId, pageable)
	if err != nil {
		return data.Page[Voucher]{}, errors.Wrap(err, "failed to read all vouchers")
	}

	vouchers, err := h.populateLineItemAccountNumber(ctx, vouchersPage.Content)
	if err != nil {
		return data.Page[Voucher]{}, errors.Wrap(err, "failed to populate account number in vouchers")
	}

	return data.NewPage(vouchers, vouchersPage.Page, vouchersPage.Size, vouchersPage.NumberOfElements)
}

func (h ReadVouchersHandler) HandleReadById(ctx context.Context, id uuid.UUID) (Voucher, error) {
	voucher, err := h.readModel.ReadById(ctx, id)
	if err != nil {
		return Voucher{}, errors.Wrap(err, "failed to read voucher")
	}

	singletonList, err := h.populateLineItemAccountNumber(ctx, []Voucher{voucher})
	if err != nil {
		return Voucher{}, errors.Wrap(err, "failed to populate account number in voucher")
	}

	return singletonList[0], nil
}

func (h ReadVouchersHandler) populateLineItemAccountNumber(ctx context.Context, vouchers []Voucher) ([]Voucher, error) {
	accountSet := make(map[uuid.UUID]void)
	for _, voucher := range vouchers {
		for _, item := range voucher.LineItems {
			accountSet[item.AccountId] = empty
		}
	}

	accounts, err := h.accountService.ReadAccountsByIds(ctx, toKeySlice[uuid.UUID, void](accountSet))
	if err != nil {
		return nil, errors.Wrap(err, "failed to read accounts by Ids")
	}

	for i := range vouchers {
		for j := range vouchers[i].LineItems {
			account, ok := accounts[vouchers[i].LineItems[j].AccountId]
			if !ok {
				return nil, errors.Errorf("account not found by id: %s", vouchers[i].LineItems[j].AccountId)
			}
			vouchers[i].LineItems[j].AccountNumber = account.AccountNumber
		}
	}

	return vouchers, nil
}
