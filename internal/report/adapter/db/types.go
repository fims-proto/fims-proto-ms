package db

import (
	"time"

	"github/fims-proto/fims-proto-ms/internal/report/domain/template"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type reportPO struct {
	Id            uuid.UUID `gorm:"type:uuid;primaryKey"`
	PeriodID      uuid.UUID `gorm:"type:uuid"`
	RefTemplateId uuid.UUID
	InnerTemplate templatePO `gorm:"foreignKey"`
}

type templatePO struct {
	Id            uuid.UUID `gorm:"type:uuid;primaryKey"`
	SobId         uuid.UUID `gorm:"type:uuid"`
	Title         string
	IsReportInner bool
	Tables        []tablePO `gorm:"foreignKey:TemplateID"`
	CreatedAt     time.Time `gorm:"<-:create"`
	UpdatedAt     time.Time
}

type tablePO struct {
	Id        uuid.UUID    `gorm:"type:uuid;primaryKey"`
	Header    headerPO     `gorm:"foreignKey:TableID"`
	LineItems []lineItemPO `gorm:"foreignKey:TableID"`
}

type headerPO struct {
	Id      uuid.UUID `gorm:"type:uuid;primaryKey"`
	TableId uuid.UUID `gorm:"type:uuid"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

type lineItemPO struct {
	Id        uuid.UUID   `gorm:"type:uuid;primaryKey"`
	TableId   uuid.UUID   `gorm:"type:uuid"`
	Formulas  []formulaPO `gorm:"foreignKey"`
	CreatedAt time.Time   `gorm:"<-:create"`
	UpdatedAt time.Time
	SumFactor int
	Values    []decimal.Decimal
}

type formulaPO struct {
	Id               uuid.UUID `gorm:"type:uuid;primaryKey"`
	AccountId        uuid.UUID `gorm:"type:uuid"`
	LineItemId       uuid.UUID `gorm:"type:uuid"`
	IsAccountFormula bool
	SumFactor        int
	Rule             string
	CreatedAt        time.Time `gorm:"<-:create"`
	UpdatedAt        time.Time
}

// table names

func (r reportPO) TableName() string {
	return "a_reports"
}

func (t templatePO) TableName() string {
	return "a_trmplate"
}

func (t tablePO) TableName() string {
	return "a_report_tables"
}

func (l lineItemPO) TableName() string {
	return "a_report_line_items"
}

func (f formulaPO) TableName() string {
	return "a_report_formulas"
}

// schema interface do we need to impl it?

// mappers

func formulaBOToPO(bo template.Formula) formulaPO {
	return formulaPO{
		Id:               bo.Id(),
		AccountId:        bo.AccountId(),
		LineItemId:       bo.LineItemId(),
		IsAccountFormula: bo.IsAccountFormula(),
		SumFactor:        bo.SumFactor(),
		Rule:             bo.Rule().String(),
	}
}

func forumaPOToBO(po formulaPO) (*template.Formula, error) {
	return template.NewFormula(
		po.Id, po.AccountId,
		po.LineItemId,
		po.IsAccountFormula,
		po.SumFactor,
		po.Rule,
	)
}
