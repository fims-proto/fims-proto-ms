package adapter

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/query"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/lineitem"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/voucher"
	"sync"

	"github.com/pkg/errors"
)

type VoucherMemoryRepository struct {
	lock *sync.RWMutex
	data map[string]voucher.Voucher
}

func NewVoucherMemoryRepository() VoucherMemoryRepository {
	return VoucherMemoryRepository{
		data: make(map[string]voucher.Voucher),
		lock: &sync.RWMutex{},
	}
}

func (h *VoucherMemoryRepository) AddVoucher(ctx context.Context, voucher *voucher.Voucher) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	_, ok := h.data[voucher.UUID()]
	if ok {
		return errors.Errorf("voucher %s exists", voucher.UUID())
	}

	h.data[voucher.UUID()] = *voucher
	return nil
}

func (h *VoucherMemoryRepository) UpdateVoucher(ctx context.Context, voucherUUID string, updateFn func(v *voucher.Voucher) (*voucher.Voucher, error)) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	v, ok := h.data[voucherUUID]
	if !ok {
		return errors.Errorf("voucher %s not exists", voucherUUID)
	}

	updatedVoucher, err := updateFn(&v)
	if err != nil {
		return errors.Wrapf(err, "voucher %s updated failed", voucherUUID)
	}
	h.data[voucherUUID] = *updatedVoucher
	return nil
}

func (h VoucherMemoryRepository) AllVouchers(ctx context.Context) ([]query.Voucher, error) {
	var result []query.Voucher
	for _, v := range h.data {
		result = append(result, voucherModelToQuery(v))
	}
	return result, nil
}

func (h VoucherMemoryRepository) VoucherForUUID(ctx context.Context, uuid string) (query.Voucher, error) {
	v, ok := h.data[uuid]
	if !ok {
		return query.Voucher{}, errors.Errorf("voucher %s not exists", uuid)
	}
	return voucherModelToQuery(v), nil
}

func voucherModelToQuery(v voucher.Voucher) query.Voucher {
	return query.Voucher{
		UUID:               v.UUID(),
		Number:             v.Number(),
		CreatedAt:          v.CreatedAt(),
		AttachmentQuantity: v.AttachmentQuantity(),
		LineItems:          itemModelToQuery(v.LineItems()),
		Debit:              v.Debit().String(),
		Credit:             v.Credit().String(),
		Creator:            v.Creator(),
		Reviewer:           v.Reviewer(),
		IsReviewed:         v.IsReviewed(),
		Auditor:            v.Auditor(),
		IsAudited:          v.IsAudited(),
	}
}

func itemModelToQuery(items []lineitem.LineItem) []query.LineItem {
	var result []query.LineItem
	for _, item := range items {
		result = append(result, query.LineItem{
			Summary:       item.Summary(),
			AccountNumber: item.AccountNumber(),
			Debit:         item.Debit().String(),
			Credit:        item.Credit().String(),
		})
	}
	return result
}
