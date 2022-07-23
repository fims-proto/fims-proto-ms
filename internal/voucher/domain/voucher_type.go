package domain

type VoucherType uint

const (
	Invalid = VoucherType(iota)
	GeneralVoucher
)

func (t VoucherType) String() string {
	switch t {
	case GeneralVoucher:
		return "general_voucher"
	default:
		return "unknown voucher type"
	}
}

func NewVoucherTypeFromString(voucherType string) (VoucherType, error) {
	if voucherType != "general_voucher" {
		return Invalid, newDomainErr(errVoucherTypeNotSupported, voucherType)
	}
	return GeneralVoucher, nil
}
