package db

import (
	"strings"
	"time"

	"github/fims-proto/fims-proto-ms/internal/journal/domain/journal_entry"
	"github/fims-proto/fims-proto-ms/internal/journal/domain/line_item"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github/fims-proto/fims-proto-ms/internal/journal/app/query"
)

type journalEntryPO struct {
	SobId              uuid.UUID `gorm:"type:uuid;uniqueIndex:journal_entries_sob_period_number_key"`
	EntryId            uuid.UUID `gorm:"type:uuid;primaryKey"`
	PeriodId           uuid.UUID `gorm:"type:uuid;uniqueIndex:journal_entries_sob_period_number_key"`
	JournalType        string
	HeaderText         string
	DocumentNumber     string `gorm:"uniqueIndex:journal_entries_sob_period_number_key"`
	AttachmentQuantity int
	Debit              decimal.Decimal
	Credit             decimal.Decimal
	Creator            uuid.UUID `gorm:"type:uuid"`
	Reviewer           uuid.UUID `gorm:"type:uuid"`
	Auditor            uuid.UUID `gorm:"type:uuid"`
	Poster             uuid.UUID `gorm:"type:uuid"`
	IsReviewed         bool
	IsAudited          bool
	IsPosted           bool
	TransactionTime    time.Time
	LineItems          []lineItemPO `gorm:"foreignKey:EntryId"`
	CreatedAt          time.Time    `gorm:"<-:create"`
	UpdatedAt          time.Time
}

type lineItemPO struct {
	EntryId   uuid.UUID `gorm:"type:uuid"`
	ItemId    uuid.UUID `gorm:"type:uuid;primaryKey"`
	AccountId uuid.UUID `gorm:"type:uuid"`
	Text      string
	Debit     decimal.Decimal
	Credit    decimal.Decimal
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

// table names

func (j journalEntryPO) TableName() string {
	return "a_journal_entries"
}

func (l lineItemPO) TableName() string {
	return "a_line_items"
}

// schemas

func (j journalEntryPO) ResolveAssociation(entity string) (string, error) {
	if entity == "" {
		return j.TableName(), nil
	}
	if strings.ToLower(entity) == strings.ToLower("LineItems") {
		return "LineItems", nil
	}
	return "", errors.Errorf("journalEntryPO doesn't have association named %s", entity)
}

// mappers

func journalEntryBOToPO(bo journal_entry.JournalEntry) journalEntryPO {
	var itemPOs []lineItemPO
	for _, item := range bo.LineItems() {
		itemPOs = append(itemPOs, lineItemBOToPO(item, bo.EntryId()))
	}

	return journalEntryPO{
		SobId:              bo.SobId(),
		EntryId:            bo.EntryId(),
		PeriodId:           bo.PeriodId(),
		JournalType:        bo.JournalType().String(),
		HeaderText:         bo.HeaderText(),
		DocumentNumber:     bo.DocumentNumber(),
		AttachmentQuantity: bo.AttachmentQuantity(),
		Debit:              bo.Debit(),
		Credit:             bo.Credit(),
		Creator:            bo.Creator(),
		Reviewer:           bo.Reviewer(),
		Auditor:            bo.Auditor(),
		Poster:             bo.Poster(),
		IsReviewed:         bo.IsReviewed(),
		IsAudited:          bo.IsAudited(),
		IsPosted:           bo.IsPosted(),
		TransactionTime:    bo.TransactionTime(),
		LineItems:          itemPOs,
	}
}

func journalEntryPOToBO(po journalEntryPO) (*journal_entry.JournalEntry, error) {
	var itemBOs []line_item.LineItem
	for _, item := range po.LineItems {
		itemBO, err := lineItemPOToBO(item)
		if err != nil {
			return nil, errors.Wrap(err, "failed to map line item")
		}

		itemBOs = append(itemBOs, *itemBO)
	}

	return journal_entry.New(
		po.SobId,
		po.EntryId,
		po.PeriodId,
		po.HeaderText,
		po.JournalType,
		po.DocumentNumber,
		po.AttachmentQuantity,
		po.Creator,
		po.Reviewer,
		po.Auditor,
		po.Poster,
		po.IsReviewed,
		po.IsAudited,
		po.IsPosted,
		po.TransactionTime,
		itemBOs,
	)
}

func journalEntryPOToDTO(po journalEntryPO) (query.JournalEntry, error) {
	var itemDTOs []query.LineItem
	for _, item := range po.LineItems {
		itemDTOs = append(itemDTOs, lineItemPOToDTO(item))
	}

	return query.JournalEntry{
		SobId:              po.SobId,
		EntryId:            po.EntryId,
		Period:             query.Period{PeriodId: po.PeriodId},
		JournalType:        po.JournalType,
		HeaderText:         po.HeaderText,
		DocumentNumber:     po.DocumentNumber,
		AttachmentQuantity: po.AttachmentQuantity,
		Debit:              po.Debit,
		Credit:             po.Credit,
		Creator:            query.User{Id: po.Creator},
		Reviewer:           query.User{Id: po.Reviewer},
		Auditor:            query.User{Id: po.Auditor},
		Poster:             query.User{Id: po.Poster},
		IsReviewed:         po.IsReviewed,
		IsAudited:          po.IsAudited,
		IsPosted:           po.IsPosted,
		TransactionTime:    po.TransactionTime,
		LineItems:          itemDTOs,
		CreatedAt:          po.CreatedAt,
		UpdatedAt:          po.UpdatedAt,
	}, nil
}

func lineItemBOToPO(bo line_item.LineItem, entryId uuid.UUID) lineItemPO {
	return lineItemPO{
		EntryId:   entryId,
		ItemId:    bo.ItemId(),
		AccountId: bo.AccountId(),
		Text:      bo.Text(),
		Debit:     bo.Debit(),
		Credit:    bo.Credit(),
	}
}

func lineItemPOToBO(po lineItemPO) (*line_item.LineItem, error) {
	return line_item.New(
		po.ItemId,
		po.AccountId,
		po.Text,
		po.Debit,
		po.Credit,
	)
}

func lineItemPOToDTO(po lineItemPO) query.LineItem {
	return query.LineItem{
		ItemId:        po.ItemId,
		AccountId:     po.AccountId,
		AccountNumber: "",
		Text:          po.Text,
		Debit:         po.Debit,
		Credit:        po.Credit,
		CreatedAt:     po.CreatedAt,
		UpdatedAt:     po.UpdatedAt,
	}
}
