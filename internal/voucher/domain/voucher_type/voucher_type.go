package voucher_type

import (
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/errors"
)

type VoucherType struct {
	slug string
}

func (t VoucherType) String() string {
	return t.slug
}

var (
	Unknown = VoucherType{""}
	General = VoucherType{"general_voucher"}
)

func FromString(s string) (VoucherType, error) {
	switch s {
	case General.slug:
		return General, nil
	}

	return Unknown, errors.NewSlugError("voucher-unknownVoucherType", fmt.Sprintf("unknown voucher type %s", s), s)
}
