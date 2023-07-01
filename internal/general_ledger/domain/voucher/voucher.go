package voucher

import (
	"time"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/voucher/voucher_type"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github/fims-proto/fims-proto-ms/internal/common/errors"
)

type Voucher struct {
	id                 uuid.UUID
	sobId              uuid.UUID
	periodId           uuid.UUID
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
	periodId uuid.UUID,
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

	if periodId == uuid.Nil {
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
		periodId:           periodId,
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

func (d *Voucher) SobId() uuid.UUID {
	return d.sobId
}

func (d *Voucher) Id() uuid.UUID {
	return d.id
}

func (d *Voucher) PeriodId() uuid.UUID {
	return d.periodId
}

func (d *Voucher) HeaderText() string {
	return d.headerText
}

func (d *Voucher) VoucherType() voucher_type.VoucherType {
	return d.voucherType
}

func (d *Voucher) DocumentNumber() string {
	return d.documentNumber
}

func (d *Voucher) AttachmentQuantity() int {
	return d.attachmentQuantity
}

func (d *Voucher) Debit() decimal.Decimal {
	return d.debit
}

func (d *Voucher) Credit() decimal.Decimal {
	return d.credit
}

func (d *Voucher) Creator() uuid.UUID {
	return d.creator
}

func (d *Voucher) Reviewer() uuid.UUID {
	return d.reviewer
}

func (d *Voucher) Auditor() uuid.UUID {
	return d.auditor
}

func (d *Voucher) Poster() uuid.UUID {
	return d.poster
}

func (d *Voucher) IsReviewed() bool {
	return d.isReviewed
}

func (d *Voucher) IsAudited() bool {
	return d.isAudited
}

func (d *Voucher) IsPosted() bool {
	return d.isPosted
}

func (d *Voucher) TransactionTime() time.Time {
	return d.transactionTime
}

func (d *Voucher) LineItems() []*LineItem {
	return d.lineItems
}
