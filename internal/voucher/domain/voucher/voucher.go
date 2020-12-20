package voucher

import (
	"github.com/shopspring/decimal"
	"time"
)

type Voucher struct {
	uuid string
	// TODO 字号
	number             uint
	createdAt          time.Time
	attachmentQuantity uint
	lineItems          []LineItem
	debit              decimal.Decimal
	credit             decimal.Decimal
	accountantUUID     string
	cashier            struct {
		uuid      string
		isChecked bool
	}
	supervisor struct {
		uuid      string
		isAudited bool
	}
}

// TODO NewVoucher

func (v Voucher) UUID() string {
	return v.uuid
}

func (v Voucher) Number() uint {
	return v.number
}

func (v Voucher) CreatedAt() time.Time {
	return v.createdAt
}

func (v Voucher) AttachmentQuantity() uint {
	return v.attachmentQuantity
}

func (v Voucher) LineItems() []LineItem {
	return v.lineItems
}

func (v Voucher) Debit() decimal.Decimal {
	return v.debit
}

func (v Voucher) Credit() decimal.Decimal {
	return v.credit
}

func (v Voucher) AccountantUUID() string {
	return v.accountantUUID
}

func (v Voucher) CashierUUID() string {
	return v.cashier.uuid
}

func (v Voucher) SupervisorUUID() string {
	return v.supervisor.uuid
}

func (v Voucher) isCheckedByCashier() bool {
	return v.cashier.isChecked
}

func (v Voucher) isAuditedBySupervisor() bool {
	return v.supervisor.isAudited
}
