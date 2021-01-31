package voucher

import (
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/lineitem"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type Voucher struct {
	uuid               string
	number             string
	createdAt          time.Time
	attachmentQuantity uint
	lineItems          []lineitem.LineItem
	debit              decimal.Decimal
	credit             decimal.Decimal
	creatorUUID        string
	reviewer           struct {
		uuid       string
		isReviewed bool
	}
	auditor struct {
		uuid      string
		isAudited bool
	}
}

func sumItems(items []lineitem.LineItem) (decimal.Decimal, error) {
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

func NewVoucher(uuid string, number string, createdAt time.Time, attachmentQuantity uint, items []lineitem.LineItem,
	creatorUUID string) (*Voucher, error) {
	if uuid == "" {
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
		uuid:               uuid,
		number:             number,
		createdAt:          createdAt,
		attachmentQuantity: attachmentQuantity,
		lineItems:          items,
		debit:              totalVal,
		credit:             totalVal,
		creatorUUID:        creatorUUID,
		reviewer: struct {
			uuid       string
			isReviewed bool
		}{
			"", false,
		},
		auditor: struct {
			uuid      string
			isAudited bool
		}{
			"", false,
		},
	}, nil
}

func (v Voucher) UUID() string {
	return v.uuid
}

func (v Voucher) Number() string {
	return v.number
}

func (v Voucher) CreatedAt() time.Time {
	return v.createdAt
}

func (v Voucher) AttachmentQuantity() uint {
	return v.attachmentQuantity
}

func (v Voucher) LineItems() []lineitem.LineItem {
	return v.lineItems
}

func (v Voucher) Debit() decimal.Decimal {
	return v.debit
}

func (v Voucher) Credit() decimal.Decimal {
	return v.credit
}

func (v Voucher) CreatorUUID() string {
	return v.creatorUUID
}

func (v Voucher) ReviewerUUID() string {
	return v.reviewer.uuid
}

func (v Voucher) AuditorUUID() string {
	return v.auditor.uuid
}

func (v Voucher) IsReviewed() bool {
	return v.reviewer.isReviewed
}

func (v Voucher) IsAudited() bool {
	return v.auditor.isAudited
}
