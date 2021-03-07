package adapter

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/query"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type VoucherMemoryRepository struct {
	lock *sync.RWMutex
	data map[uuid.UUID]domain.Voucher
}

func NewVoucherMemoryRepository() VoucherMemoryRepository {
	return VoucherMemoryRepository{
		data: make(map[uuid.UUID]domain.Voucher),
		lock: &sync.RWMutex{},
	}
}

func (h *VoucherMemoryRepository) AddVoucher(ctx context.Context, voucher *domain.Voucher) (uuid.UUID, error) {
	h.lock.Lock()
	defer h.lock.Unlock()

	_, ok := h.data[voucher.UUID()]
	if ok {
		return uuid.Nil, errors.Errorf("voucher %s exists", voucher.UUID())
	}

	h.data[voucher.UUID()] = *voucher
	return voucher.UUID(), nil
}

func (h *VoucherMemoryRepository) UpdateVoucher(ctx context.Context, voucherUUID uuid.UUID, updateFn func(v *domain.Voucher) (*domain.Voucher, error)) error {
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
		result = append(result, mapFromDomainVoucher(v))
	}
	return result, nil
}

func (h VoucherMemoryRepository) VoucherByUUID(ctx context.Context, voucherUUID uuid.UUID) (query.Voucher, error) {
	v, ok := h.data[voucherUUID]
	if !ok {
		return query.Voucher{}, errors.Errorf("voucher %s not exists", voucherUUID)
	}
	return mapFromDomainVoucher(v), nil
}

func mapFromDomainVoucher(v domain.Voucher) query.Voucher {
	return query.Voucher{
		UUID:               v.UUID(),
		Number:             v.Number(),
		CreatedAt:          v.CreatedAt(),
		AttachmentQuantity: v.AttachmentQuantity(),
		LineItems:          mapFromDomainLineItem(v.LineItems()),
		Debit:              v.Debit().String(),
		Credit:             v.Credit().String(),
		Creator:            v.Creator(),
		Reviewer:           v.Reviewer(),
		IsReviewed:         v.IsReviewed(),
		Auditor:            v.Auditor(),
		IsAudited:          v.IsAudited(),
		IsPosted:           v.IsPosted(),
	}
}

func mapFromDomainLineItem(items []domain.LineItem) []query.LineItem {
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
