package domain

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type Voucher struct {
	sob                string
	uuid               uuid.UUID
	voucherType        VoucherType
	number             string
	attachmentQuantity uint
	lineItems          []LineItem
	debit              decimal.Decimal
	credit             decimal.Decimal
	creator            string
	reviewer           string
	isReviewed         bool
	auditor            string
	isAudited          bool
	isPosted           bool
}

func sumItems(items []LineItem) (decimal.Decimal, error) {
	if len(items) == 0 {
		return decimal.Decimal{}, errors.New("lineitem cannot be empty")
	}

	var debitInTotal decimal.Decimal
	var creditInTotal decimal.Decimal
	for _, item := range items {
		debitInTotal = debitInTotal.Add(item.Debit())
		creditInTotal = creditInTotal.Add(item.Credit())
	}

	if !debitInTotal.Equal(creditInTotal) {
		return decimal.Decimal{}, errors.New("debit and credit not equal")
	}
	return debitInTotal, nil
}

func NewVoucher(sob string, voucherUUID uuid.UUID, voucherType VoucherType, number string, attachmentQuantity uint, items []LineItem, creator string) (*Voucher, error) {
	if sob == "" {
		return nil, errors.New("empty sob")
	}
	if voucherUUID == uuid.Nil {
		return nil, errors.New("empty voucher uuid")
	}
	if number == "" {
		return nil, errors.New("empty voucher number")
	}

	totalVal, err := sumItems(items)
	if err != nil {
		return nil, err
	}

	return &Voucher{
		sob:                sob,
		uuid:               voucherUUID,
		voucherType:        voucherType,
		number:             number,
		attachmentQuantity: attachmentQuantity,
		lineItems:          items,
		debit:              totalVal,
		credit:             totalVal,
		creator:            creator,
		reviewer:           "",
		isReviewed:         false,
		auditor:            "",
		isAudited:          false,
		isPosted:           false,
	}, nil
}

func (v Voucher) Sob() string {
	return v.sob
}

func (v Voucher) UUID() uuid.UUID {
	return v.uuid
}

func (v Voucher) Type() VoucherType {
	return v.voucherType
}

func (v Voucher) Number() string {
	return v.number
}

func (v Voucher) AttachmentQuantity() uint {
	return v.attachmentQuantity
}

func (v Voucher) LineItems() []LineItem {
	return v.lineItems
}

func (v Voucher) Debit() decimal.Decimal {
	return v.debit
}

func (v Voucher) Credit() decimal.Decimal {
	return v.credit
}

func (v Voucher) Creator() string {
	return v.creator
}

func (v Voucher) Reviewer() string {
	return v.reviewer
}

func (v Voucher) Auditor() string {
	return v.auditor
}

func (v Voucher) IsReviewed() bool {
	return v.isReviewed
}

func (v Voucher) IsAudited() bool {
	return v.isAudited
}

func (v Voucher) IsPosted() bool {
	return v.isPosted
}
