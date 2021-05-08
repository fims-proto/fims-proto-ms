package domain

import "github.com/pkg/errors"

type VoucherType uint

const (
	Invalid = VoucherType(iota)
	GeneralVoucher
)

func (t VoucherType) String() string {
	switch t {
	case GeneralVoucher:
		return "GENERAL_VOUCHER"
	default:
		panic("unknown voucher type")
	}
}

func NewVoucherTypeFromString(voucherType string) (VoucherType, error) {
	if voucherType != "GENERAL_VOUCHER" {
		return Invalid, errors.Errorf("voucher type %s not supported", voucherType)
	}
	return GeneralVoucher, nil
}
