package voucher

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/common/utils"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_account"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_account_category"
)

type LineItem struct {
	id                uuid.UUID
	accountId         uuid.UUID
	account           *account.Account
	auxiliaryAccounts []*auxiliary_account.AuxiliaryAccount
	text              string
	debit             decimal.Decimal
	credit            decimal.Decimal
}

func NewLineItem(
	id uuid.UUID,
	accountId uuid.UUID,
	account *account.Account,
	auxiliaryAccounts []*auxiliary_account.AuxiliaryAccount,
	text string,
	debit decimal.Decimal,
	credit decimal.Decimal,
) (*LineItem, error) {
	if id == uuid.Nil {
		return nil, errors.NewSlugError("lineItem-emptyId")
	}

	if accountId == uuid.Nil {
		return nil, errors.NewSlugError("lineItem-emptyAccountId")
	}

	if account == nil {
		return nil, errors.NewSlugError("lineItem-nilAccount")
	}

	if len(auxiliaryAccounts) != len(account.AuxiliaryAccountCategories()) {
		return nil, errors.NewSlugError("lineItem-unmatchedAuxiliaryAccount")
	}

	for _, auxiliaryAccount := range auxiliaryAccounts {
		if auxiliaryAccount == nil {
			return nil, errors.NewSlugError("lineItem-nilAuxiliaryAccount")
		}
	}

	if text == "" {
		return nil, errors.NewSlugError("lineItem-emptyText")
	}

	if debit.IsZero() && credit.IsZero() {
		return nil, errors.NewSlugError("lineItem-emptyDebitCredit")
	}

	if !debit.IsZero() && !credit.IsZero() {
		return nil, errors.NewSlugError("lineItem-debitCreditDuplicated")
	}

	// validate each auxiliary account
	categorySet := utils.SliceToSet(account.AuxiliaryAccountCategories(), func(category *auxiliary_account_category.AuxiliaryAccountCategory) uuid.UUID {
		return category.Id()
	})
	for _, auxiliaryAccount := range auxiliaryAccounts {
		if _, ok := categorySet[auxiliaryAccount.Category().Id()]; !ok {
			return nil, errors.NewSlugError("lineItem-invalidAuxiliaryAccount", auxiliaryAccount.Title())
		}
	}

	return &LineItem{
		id:                id,
		accountId:         accountId,
		account:           account,
		auxiliaryAccounts: auxiliaryAccounts,
		text:              text,
		debit:             debit,
		credit:            credit,
	}, nil
}

func (i LineItem) Id() uuid.UUID {
	return i.id
}

func (i LineItem) AccountId() uuid.UUID {
	return i.accountId
}

func (i LineItem) Account() *account.Account {
	return i.account
}

func (i LineItem) AuxiliaryAccounts() []*auxiliary_account.AuxiliaryAccount {
	return i.auxiliaryAccounts
}

func (i LineItem) Text() string {
	return i.text
}

func (i LineItem) Debit() decimal.Decimal {
	return i.debit
}

func (i LineItem) Credit() decimal.Decimal {
	return i.credit
}
