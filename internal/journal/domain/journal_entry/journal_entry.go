package journal_entry

import (
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/journal/domain/journal_type"
	"github/fims-proto/fims-proto-ms/internal/journal/domain/line_item"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type JournalEntry struct {
	sobId              uuid.UUID
	entryId            uuid.UUID
	periodId           uuid.UUID
	journalType        journal_type.JournalType
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
	sobId, entryId, periodId uuid.UUID,
	headerText, journalType, documentNumber string,
	attachmentQuantity int,
	creator, reviewer, auditor, poster uuid.UUID,
	isReviewed, isAudited, isPosted bool,
	transactionTime time.Time,
	lineItems []line_item.LineItem,
) (*JournalEntry, error) {
	if sobId == uuid.Nil {
		return nil, commonErrors.NewSlugError("journalEntry-emptySobId", "empty sob id")
	}

	if entryId == uuid.Nil {
		return nil, commonErrors.NewSlugError("journalEntry-emptyId", "empty entry id")
	}

	if periodId == uuid.Nil {
		return nil, commonErrors.NewSlugError("journalEntry-emptyPeriodId", "empty period id")
	}

	if headerText == "" {
		return nil, commonErrors.NewSlugError("journalEntry-emptySummary", "empty header text")
	}

	dt, err := journal_type.FromString(journalType)
	if err != nil {
		return nil, err
	}

	if documentNumber == "" {
		return nil, commonErrors.NewSlugError("journalEntry-emptyNumber", "empty document number")
	}

	if attachmentQuantity < 0 {
		return nil, commonErrors.NewSlugError("journalEntry-emptyAttachment", "attachment quantity cannot lesser than 0")
	}

	if creator == uuid.Nil {
		return nil, commonErrors.NewSlugError("journalEntry-emptyCreator", "empty creator")
	}

	if isReviewed && reviewer == uuid.Nil {
		return nil, commonErrors.NewSlugError("journalEntry-emptyReviewer", "empty reviewer")
	}

	if isAudited && auditor == uuid.Nil {
		return nil, commonErrors.NewSlugError("journalEntry-emptyAuditor", "empty auditor")
	}

	if isPosted && poster == uuid.Nil {
		return nil, commonErrors.NewSlugError("journalEntry-emptyPoster", "empty poster")
	}

	if isPosted && (!isReviewed || !isAudited) {
		return nil, commonErrors.NewSlugError("journalEntry-invalidPostStatus", "invalid post status")
	}

	if transactionTime.IsZero() {
		return nil, commonErrors.NewSlugError("journalEntry-zeroTransactionTime", "zero transaction time")
	}

	if transactionTime.After(time.Now()) {
		return nil, commonErrors.NewSlugError("journalEntry-futureTransactionTime", "transaction time is in future")
	}

	totalVal, err := sumLineItems(lineItems)
	if err != nil {
		return nil, err
	}

	return &JournalEntry{
		sobId:              sobId,
		entryId:            entryId,
		periodId:           periodId,
		headerText:         headerText,
		journalType:        dt,
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

func (d *JournalEntry) SobId() uuid.UUID {
	return d.sobId
}

func (d *JournalEntry) EntryId() uuid.UUID {
	return d.entryId
}

func (d *JournalEntry) PeriodId() uuid.UUID {
	return d.periodId
}

func (d *JournalEntry) HeaderText() string {
	return d.headerText
}

func (d *JournalEntry) JournalType() journal_type.JournalType {
	return d.journalType
}

func (d *JournalEntry) DocumentNumber() string {
	return d.documentNumber
}

func (d *JournalEntry) AttachmentQuantity() int {
	return d.attachmentQuantity
}

func (d *JournalEntry) Debit() decimal.Decimal {
	return d.debit
}

func (d *JournalEntry) Credit() decimal.Decimal {
	return d.credit
}

func (d *JournalEntry) Creator() uuid.UUID {
	return d.creator
}

func (d *JournalEntry) Reviewer() uuid.UUID {
	return d.reviewer
}

func (d *JournalEntry) Auditor() uuid.UUID {
	return d.auditor
}

func (d *JournalEntry) Poster() uuid.UUID {
	return d.poster
}

func (d *JournalEntry) IsReviewed() bool {
	return d.isReviewed
}

func (d *JournalEntry) IsAudited() bool {
	return d.isAudited
}

func (d *JournalEntry) IsPosted() bool {
	return d.isPosted
}

func (d *JournalEntry) TransactionTime() time.Time {
	return d.transactionTime
}

func (d *JournalEntry) LineItems() []line_item.LineItem {
	return d.lineItems
}
