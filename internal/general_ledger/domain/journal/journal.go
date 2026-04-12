package journal

import (
	"github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/period"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/transaction_date"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Journal struct {
	id                 uuid.UUID
	sobId              uuid.UUID
	periodId           uuid.UUID
	period             *period.Period
	headerText         string
	documentNumber     string
	journalType        JournalType
	referenceJournalId uuid.UUID
	attachmentQuantity int
	amount             decimal.Decimal
	creator            string
	reviewer           string
	auditor            string
	poster             string
	isReviewed         bool
	isAudited          bool
	isPosted           bool
	transactionDate    transaction_date.TransactionDate
	journalLines       []*JournalLine
}

func New(
	id uuid.UUID,
	sobId uuid.UUID,
	period *period.Period,
	headerText string,
	documentNumber string,
	journalType JournalType,
	referenceJournalId uuid.UUID,
	attachmentQuantity int,
	creator string,
	reviewer string,
	auditor string,
	poster string,
	isReviewed bool,
	isAudited bool,
	isPosted bool,
	transactionDate transaction_date.TransactionDate,
	journalLines []*JournalLine,
) (*Journal, error) {
	if id == uuid.Nil {
		return nil, errors.NewInternalError(errors.SlugJournalEmptyId)
	}

	if sobId == uuid.Nil {
		return nil, errors.NewInternalError(errors.SlugEmptySobId)
	}

	if period == nil {
		return nil, errors.NewInternalError(errors.SlugJournalEmptyPeriod)
	}

	if period.Id() == uuid.Nil {
		return nil, errors.NewInternalError(errors.SlugJournalEmptyPeriodId)
	}

	if headerText == "" {
		return nil, errors.NewInvalidInputError(errors.SlugJournalEmptyHeaderText)
	}

	if documentNumber == "" {
		return nil, errors.NewInternalError(errors.SlugJournalEmptyNumber)
	}

	if attachmentQuantity < 0 {
		return nil, errors.NewInvalidInputError(errors.SlugJournalInvalidAttachmentQty)
	}

	if isEmptyUser(creator) {
		return nil, errors.NewInternalError(errors.SlugJournalEmptyCreator)
	}

	if isReviewed && isEmptyUser(reviewer) {
		return nil, errors.NewInternalError(errors.SlugJournalEmptyReviewer)
	}

	if isAudited && isEmptyUser(auditor) {
		return nil, errors.NewInternalError(errors.SlugJournalEmptyAuditor)
	}

	if isPosted && isEmptyUser(poster) {
		return nil, errors.NewInternalError(errors.SlugJournalEmptyPoster)
	}

	if isPosted && (!isReviewed || !isAudited) {
		return nil, errors.NewInternalError(errors.SlugJournalInvalidPostStatus)
	}

	if transactionDate.IsZero() {
		return nil, errors.NewInternalError(errors.SlugJournalZeroTransactionDate)
	}

	if !journalType.IsValid() {
		return nil, errors.NewInternalError(errors.SlugJournalInvalidJournalType)
	}

	if journalType.RequiresReferenceJournal() && referenceJournalId == uuid.Nil {
		return nil, errors.NewInternalError(errors.SlugJournalMissingReferenceId)
	}

	if !journalType.RequiresReferenceJournal() && referenceJournalId != uuid.Nil {
		return nil, errors.NewInternalError(errors.SlugJournalUnexpectedReferenceId)
	}

	totalVal, err := sumJournalLines(journalLines)
	if err != nil {
		return nil, err
	}

	return &Journal{
		id:                 id,
		sobId:              sobId,
		periodId:           period.Id(),
		period:             period,
		headerText:         headerText,
		documentNumber:     documentNumber,
		journalType:        journalType,
		referenceJournalId: referenceJournalId,
		attachmentQuantity: attachmentQuantity,
		amount:             totalVal,
		creator:            creator,
		reviewer:           reviewer,
		auditor:            auditor,
		poster:             poster,
		isReviewed:         isReviewed,
		isAudited:          isAudited,
		isPosted:           isPosted,
		transactionDate:    transactionDate,
		journalLines:       journalLines,
	}, nil
}

func (j *Journal) SobId() uuid.UUID {
	return j.sobId
}

func (j *Journal) Id() uuid.UUID {
	return j.id
}

func (j *Journal) PeriodId() uuid.UUID {
	return j.periodId
}

func (j *Journal) Period() *period.Period {
	return j.period
}

func (j *Journal) HeaderText() string {
	return j.headerText
}

func (j *Journal) DocumentNumber() string {
	return j.documentNumber
}

func (j *Journal) JournalType() JournalType {
	return j.journalType
}

func (j *Journal) ReferenceJournalId() uuid.UUID {
	return j.referenceJournalId
}

func (j *Journal) AttachmentQuantity() int {
	return j.attachmentQuantity
}

func (j *Journal) Amount() decimal.Decimal {
	return j.amount
}

func (j *Journal) Creator() string {
	return j.creator
}

func (j *Journal) Reviewer() string {
	return j.reviewer
}

func (j *Journal) Auditor() string {
	return j.auditor
}

func (j *Journal) Poster() string {
	return j.poster
}

func (j *Journal) IsReviewed() bool {
	return j.isReviewed
}

func (j *Journal) IsAudited() bool {
	return j.isAudited
}

func (j *Journal) IsPosted() bool {
	return j.isPosted
}

func (j *Journal) TransactionDate() transaction_date.TransactionDate {
	return j.transactionDate
}

func (j *Journal) JournalLines() []*JournalLine {
	return j.journalLines
}
