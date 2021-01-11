package adapter

import (
	"context"
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/query"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/lineitem"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/voucher"
	"sync"
)

type VoucherMemoryRepository struct{
	sync.RWMutex
	data map[string]voucher.Voucher
}

func NewVoucherMemoryRepository() VoucherMemoryRepository {
	return VoucherMemoryRepository{data: make(map[string]voucher.Voucher)}
}

func (h *VoucherMemoryRepository) AddVoucher(ctx context.Context, voucher *voucher.Voucher) error {
	_, ok := h.data[voucher.UUID()]
	if ok {
		return errors.Errorf("voucher %s exists", voucher.UUID())
	}
	h.data[voucher.UUID()] = *voucher
	return nil
}

func (h *VoucherMemoryRepository) UpdateVoucher(ctx context.Context, voucherUUID string, updateFn func(v *voucher.Voucher) (*voucher.Voucher, error)) error {
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
		result = append(result, query.Voucher{
			UUID:               v.UUID(),
			Number:             v.Number(),
			CreatedAt:          v.CreatedAt(),
			AttachmentQuantity: v.AttachmentQuantity(),
			LineItems:          h.itemModelToQuery(v.LineItems()),
			Debit:              v.Debit().String(),
			Credit:             v.Credit().String(),
			CreatorUUID:        v.CreatorUUID(),
			ReviewerUUID:       v.ReviewerUUID(),
			IsReviewed:         v.IsReviewed(),
			AuditorUUID:        v.AuditorUUID(),
			IsAudited:          v.IsAudited(),
		})
	}
	return result, nil
}

func (h VoucherMemoryRepository) itemModelToQuery(items []lineitem.LineItem) []query.LineItem {
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
