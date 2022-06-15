package domain

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
		return "unknown voucher type"
	}
}

func NewVoucherTypeFromString(voucherType string) (VoucherType, error) {
	if voucherType != "GENERAL_VOUCHER" {
		return Invalid, newDomainErr(errVoucherTypeNotSupported, voucherType)
	}
	return GeneralVoucher, nil
}
