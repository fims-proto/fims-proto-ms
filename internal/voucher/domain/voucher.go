package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type Voucher struct {
	id                 uuid.UUID
	sobId              uuid.UUID
	voucherType        VoucherType
	number             string
	attachmentQuantity uint
	lineItems          []*LineItem
	debit              decimal.Decimal
	credit             decimal.Decimal
	creator            string
	reviewer           string
	isReviewed         bool
	auditor            string
	isAudited          bool
	isPosted           bool
	transactionTime    time.Time
}

func NewVoucher(id, sobId uuid.UUID, voucherType, number string, attachmentQuantity uint, items []*LineItem,
	creator, reviewer, auditor string, isReviewed, isAudited, isPosted bool, transactionTime time.Time,
) (*Voucher, error) {
	if id == uuid.Nil {
		return nil, errors.New("empty voucher id")
	}
	if sobId == uuid.Nil {
		return nil, errors.New("empty sobId")
	}
	if number == "" {
		return nil, errors.New("empty voucher number")
	}

	vt, err := NewVoucherTypeFromString(voucherType)
	if err != nil {
		return nil, errors.Wrap(err, "invalid voucher type")
	}

	if creator == "" {
		return nil, errors.New("empty creator")
	}

	if isReviewed && reviewer == "" {
		return nil, errors.New("empty reviewer")
	}

	if isAudited && auditor == "" {
		return nil, errors.New("empty auditor")
	}

	if isPosted && (!isReviewed || !isAudited) {
		return nil, errors.New("invalid posted status")
	}

	totalVal, err := sumItems(items)
	if err != nil {
		return nil, err
	}

	if transactionTime.IsZero() {
		return nil, errors.New("zero transaction time")
	}

	return &Voucher{
		sobId:              sobId,
		id:                 id,
		voucherType:        vt,
		number:             number,
		attachmentQuantity: attachmentQuantity,
		lineItems:          items,
		debit:              totalVal,
		credit:             totalVal,
		creator:            creator,
		reviewer:           reviewer,
		isReviewed:         isReviewed,
		auditor:            auditor,
		isAudited:          isAudited,
		isPosted:           isPosted,
		transactionTime:    transactionTime,
	}, nil
}

func sumItems(items []*LineItem) (decimal.Decimal, error) {
	if len(items) == 0 {
		return decimal.Decimal{}, errors.New("line item cannot be empty")
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

func (v Voucher) SobId() uuid.UUID {
	return v.sobId
}

func (v Voucher) Id() uuid.UUID {
	return v.id
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

func (v Voucher) LineItems() []*LineItem {
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

func (v Voucher) TransactionTime() time.Time {
	return v.transactionTime
}
