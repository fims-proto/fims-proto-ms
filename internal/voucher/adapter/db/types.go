package db

import (
	"strings"
	"time"

	"github/fims-proto/fims-proto-ms/internal/voucher/domain/line_item"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/voucher"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/query"
)

type voucherPO struct {
	SobId              uuid.UUID `gorm:"type:uuid;uniqueIndex:vouchers_sob_period_number_key"`
	Id                 uuid.UUID `gorm:"type:uuid;primaryKey"`
	PeriodId           uuid.UUID `gorm:"type:uuid;uniqueIndex:vouchers_sob_period_number_key"`
	VoucherType        string
	HeaderText         string
	DocumentNumber     string `gorm:"uniqueIndex:vouchers_sob_period_number_key"`
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
	LineItems          []lineItemPO `gorm:"foreignKey:VoucherId"`
	CreatedAt          time.Time    `gorm:"<-:create"`
	UpdatedAt          time.Time
}

type lineItemPO struct {
	VoucherId uuid.UUID `gorm:"type:uuid"`
	Id        uuid.UUID `gorm:"type:uuid;primaryKey"`
	AccountId uuid.UUID `gorm:"type:uuid"`
	Text      string
	Debit     decimal.Decimal
	Credit    decimal.Decimal
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

// table names

func (j voucherPO) TableName() string {
	return "a_vouchers"
}

func (l lineItemPO) TableName() string {
	return "a_line_items"
}

// schemas

func (j voucherPO) ResolveAssociation(entity string) (string, error) {
	if entity == "" {
		return j.TableName(), nil
	}
	if strings.EqualFold(entity, "LineItems") {
		return "LineItems", nil
	}
	return "", errors.Errorf("voucherPO doesn't have association named %s", entity)
}

// mappers

func voucherBOToPO(bo voucher.Voucher) voucherPO {
	var itemPOs []lineItemPO
	for _, item := range bo.LineItems() {
		itemPOs = append(itemPOs, lineItemBOToPO(item, bo.Id()))
	}

	return voucherPO{
		SobId:              bo.SobId(),
		Id:                 bo.Id(),
		PeriodId:           bo.PeriodId(),
		VoucherType:        bo.VoucherType().String(),
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

func voucherPOToBO(po voucherPO) (*voucher.Voucher, error) {
	var itemBOs []line_item.LineItem
	for _, item := range po.LineItems {
		itemBO, err := lineItemPOToBO(item)
		if err != nil {
			return nil, errors.Wrap(err, "failed to map line item")
		}

		itemBOs = append(itemBOs, *itemBO)
	}

	return voucher.New(
		po.SobId,
		po.Id,
		po.PeriodId,
		po.HeaderText,
		po.VoucherType,
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

func voucherPOToDTO(po voucherPO) (query.Voucher, error) {
	var itemDTOs []query.LineItem
	for _, item := range po.LineItems {
		itemDTOs = append(itemDTOs, lineItemPOToDTO(item))
	}

	return query.Voucher{
		SobId:              po.SobId,
		Id:                 po.Id,
		Period:             query.Period{PeriodId: po.PeriodId},
		VoucherType:        po.VoucherType,
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

func lineItemBOToPO(bo line_item.LineItem, voucherId uuid.UUID) lineItemPO {
	return lineItemPO{
		VoucherId: voucherId,
		Id:        bo.Id(),
		AccountId: bo.AccountId(),
		Text:      bo.Text(),
		Debit:     bo.Debit(),
		Credit:    bo.Credit(),
	}
}

func lineItemPOToBO(po lineItemPO) (*line_item.LineItem, error) {
	return line_item.New(
		po.Id,
		po.AccountId,
		po.Text,
		po.Debit,
		po.Credit,
	)
}

func lineItemPOToDTO(po lineItemPO) query.LineItem {
	return query.LineItem{
		Id:            po.Id,
		AccountId:     po.AccountId,
		AccountNumber: "",
		Text:          po.Text,
		Debit:         po.Debit,
		Credit:        po.Credit,
		CreatedAt:     po.CreatedAt,
		UpdatedAt:     po.UpdatedAt,
	}
}
