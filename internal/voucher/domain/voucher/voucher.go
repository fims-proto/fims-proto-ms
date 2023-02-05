package voucher

import (
	"time"

	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/line_item"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/voucher_type"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Voucher struct {
	sobId              uuid.UUID
	id                 uuid.UUID
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
	lineItems          []line_item.LineItem
}

func New(
	sobId, voucherId, periodId uuid.UUID,
	headerText, voucherType, documentNumber string,
	attachmentQuantity int,
	creator, reviewer, auditor, poster uuid.UUID,
	isReviewed, isAudited, isPosted bool,
	transactionTime time.Time,
	lineItems []line_item.LineItem,
) (*Voucher, error) {
	if sobId == uuid.Nil {
		return nil, commonErrors.NewSlugError("voucher-emptySobId", "empty sob id")
	}

	if voucherId == uuid.Nil {
		return nil, commonErrors.NewSlugError("voucher-emptyId", "empty voucher id")
	}

	if periodId == uuid.Nil {
		return nil, commonErrors.NewSlugError("voucher-emptyPeriodId", "empty period id")
	}

	if headerText == "" {
		return nil, commonErrors.NewSlugError("voucher-emptySummary", "empty header text")
	}

	dt, err := voucher_type.FromString(voucherType)
	if err != nil {
		return nil, err
	}

	if documentNumber == "" {
		return nil, commonErrors.NewSlugError("voucher-emptyNumber", "empty document number")
	}

	if attachmentQuantity < 0 {
		return nil, commonErrors.NewSlugError("voucher-emptyAttachment", "attachment quantity cannot lesser than 0")
	}

	if creator == uuid.Nil {
		return nil, commonErrors.NewSlugError("voucher-emptyCreator", "empty creator")
	}

	if isReviewed && reviewer == uuid.Nil {
		return nil, commonErrors.NewSlugError("voucher-emptyReviewer", "empty reviewer")
	}

	if isAudited && auditor == uuid.Nil {
		return nil, commonErrors.NewSlugError("voucher-emptyAuditor", "empty auditor")
	}

	if isPosted && poster == uuid.Nil {
		return nil, commonErrors.NewSlugError("voucher-emptyPoster", "empty poster")
	}

	if isPosted && (!isReviewed || !isAudited) {
		return nil, commonErrors.NewSlugError("voucher-invalidPostStatus", "invalid post status")
	}

	if transactionTime.IsZero() {
		return nil, commonErrors.NewSlugError("voucher-zeroTransactionTime", "zero transaction time")
	}

	if transactionTime.After(time.Now()) {
		return nil, commonErrors.NewSlugError("voucher-futureTransactionTime", "transaction time is in future")
	}

	totalVal, err := sumLineItems(lineItems)
	if err != nil {
		return nil, err
	}

	return &Voucher{
		sobId:              sobId,
		id:                 voucherId,
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

func (d *Voucher) LineItems() []line_item.LineItem {
	return d.lineItems
}
