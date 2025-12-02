package voucher

import (
	"time"

	"github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/period"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/voucher/voucher_type"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Voucher struct {
	id                 uuid.UUID
	sobId              uuid.UUID
	periodId           uuid.UUID
	period             *period.Period
	voucherType        voucher_type.VoucherType
	headerText         string
	documentNumber     string
	attachmentQuantity int
	debit              decimal.Decimal
	credit             decimal.Decimal
	creator            uuid.UUID
	reviewer           uuid.UUID
	auditor            uuid.UUID
	poster             uuid.UUID
	isReviewed         bool
	isAudited          bool
	isPosted           bool
	transactionTime    time.Time
	lineItems          []*LineItem
}

func New(
	id uuid.UUID,
	sobId uuid.UUID,
	period *period.Period,
	voucherType string,
	headerText string,
	documentNumber string,
	attachmentQuantity int,
	creator uuid.UUID,
	reviewer uuid.UUID,
	auditor uuid.UUID,
	poster uuid.UUID,
	isReviewed bool,
	isAudited bool,
	isPosted bool,
	transactionTime time.Time,
	lineItems []*LineItem,
) (*Voucher, error) {
	if id == uuid.Nil {
		return nil, errors.NewSlugError("voucher-emptyId")
	}

	if sobId == uuid.Nil {
		return nil, errors.NewSlugError("emptySobId")
	}

	if period == nil {
		return nil, errors.NewSlugError("voucher-emptyPeriod")
	}

	if period.Id() == uuid.Nil {
		return nil, errors.NewSlugError("voucher-emptyPeriodId")
	}

	if headerText == "" {
		return nil, errors.NewSlugError("voucher-emptyHeaderText")
	}

	dt, err := voucher_type.FromString(voucherType)
	if err != nil {
		return nil, err
	}

	if documentNumber == "" {
		return nil, errors.NewSlugError("voucher-emptyNumber")
	}

	if attachmentQuantity < 0 {
		return nil, errors.NewSlugError("voucher-invalidAttachmentQuantity")
	}

	if creator == uuid.Nil {
		return nil, errors.NewSlugError("voucher-emptyCreator")
	}

	if isReviewed && reviewer == uuid.Nil {
		return nil, errors.NewSlugError("voucher-emptyReviewer")
	}

	if isAudited && auditor == uuid.Nil {
		return nil, errors.NewSlugError("voucher-emptyAuditor")
	}

	if isPosted && poster == uuid.Nil {
		return nil, errors.NewSlugError("voucher-emptyPoster")
	}

	if isPosted && (!isReviewed || !isAudited) {
		return nil, errors.NewSlugError("voucher-invalidPostStatus")
	}

	if transactionTime.IsZero() {
		return nil, errors.NewSlugError("voucher-zeroTransactionTime")
	}

	totalVal, err := sumLineItems(lineItems)
	if err != nil {
		return nil, err
	}

	return &Voucher{
		id:                 id,
		sobId:              sobId,
		periodId:           period.Id(),
		period:             period,
		headerText:         headerText,
		voucherType:        dt,
		documentNumber:     documentNumber,
		attachmentQuantity: attachmentQuantity,
		debit:              totalVal,
		credit:             totalVal,
		creator:            creator,
		reviewer:           reviewer,
		auditor:            auditor,
		poster:             poster,
		isReviewed:         isReviewed,
		isAudited:          isAudited,
		isPosted:           isPosted,
		transactionTime:    transactionTime,
		lineItems:          lineItems,
	}, nil
}

func (v *Voucher) SobId() uuid.UUID {
	return v.sobId
}

func (v *Voucher) Id() uuid.UUID {
	return v.id
}

func (v *Voucher) PeriodId() uuid.UUID {
	return v.periodId
}

func (v *Voucher) Period() *period.Period {
	return v.period
}

func (v *Voucher) HeaderText() string {
	return v.headerText
}

func (v *Voucher) VoucherType() voucher_type.VoucherType {
	return v.voucherType
}

func (v *Voucher) DocumentNumber() string {
	return v.documentNumber
}

func (v *Voucher) AttachmentQuantity() int {
	return v.attachmentQuantity
}

func (v *Voucher) Debit() decimal.Decimal {
	return v.debit
}

func (v *Voucher) Credit() decimal.Decimal {
	return v.credit
}

func (v *Voucher) Creator() uuid.UUID {
	return v.creator
}

func (v *Voucher) Reviewer() uuid.UUID {
	return v.reviewer
}

func (v *Voucher) Auditor() uuid.UUID {
	return v.auditor
}

func (v *Voucher) Poster() uuid.UUID {
	return v.poster
}

func (v *Voucher) IsReviewed() bool {
	return v.isReviewed
}

func (v *Voucher) IsAudited() bool {
	return v.isAudited
}

func (v *Voucher) IsPosted() bool {
	return v.isPosted
}

func (v *Voucher) TransactionTime() time.Time {
	return v.transactionTime
}

func (v *Voucher) LineItems() []*LineItem {
	return v.lineItems
}
