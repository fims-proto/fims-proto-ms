package http

import (
	"context"
	"net/http"
	"strconv"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/command"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account/class"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type ReadAllAccountsInput struct {
	SobId uuid.UUID `path:"sobId"`
}

type ReadAllAccountsOutput struct {
	Body []AccountSlimResponse
}

// ReadAllAccounts lists all accounts for a SoB.
func (h Handler) ReadAllAccounts(ctx context.Context, input *ReadAllAccountsInput) (*ReadAllAccountsOutput, error) {
	accounts, err := h.app.Queries.AllAccounts.Handle(ctx, input.SobId)
	if err != nil {
		return nil, err
	}
	vos := make([]AccountSlimResponse, len(accounts))
	for i, a := range accounts {
		vos[i] = accountDTOToSlimVO(a)
	}
	return &ReadAllAccountsOutput{Body: vos}, nil
}

type ReadAccountClassesInput struct {
	SobId uuid.UUID `path:"sobId"`
}

type ReadAccountClassesOutput struct {
	Body []AccountClass
}

// ReadAccountClasses returns supported account class and group mappings.
func (h Handler) ReadAccountClasses(_ context.Context, _ *ReadAccountClassesInput) (*ReadAccountClassesOutput, error) {
	var resp []AccountClass
	for _, c := range class.Classes {
		var groups []string
		for _, g := range c.Groups {
			groups = append(groups, strconv.Itoa(int(g)))
		}
		resp = append(resp, AccountClass{
			Class:  strconv.Itoa(int(c.Class)),
			Groups: groups,
		})
	}
	return &ReadAccountClassesOutput{Body: resp}, nil
}

type ReadAccountByIdInput struct {
	SobId     uuid.UUID `path:"sobId"`
	AccountId uuid.UUID `path:"accountId"`
}

type ReadAccountByIdOutput struct {
	Body AccountDetailResponse
}

// ReadAccountById returns one account by ID.
func (h Handler) ReadAccountById(ctx context.Context, input *ReadAccountByIdInput) (*ReadAccountByIdOutput, error) {
	v, err := h.app.Queries.AccountById.Handle(ctx, input.AccountId)
	if err != nil {
		return nil, err
	}
	if v.Id == uuid.Nil {
		return nil, huma.Error404NotFound("account not found")
	}
	return &ReadAccountByIdOutput{Body: accountDTOToDetailVO(v)}, nil
}

type CreateAccountInput struct {
	SobId uuid.UUID `path:"sobId"`
	Body  CreateAccountRequest
}

type CreateAccountOutput struct {
	Status int
	Body   AccountDetailResponse
}

// CreateAccount creates an account and returns created detail.
func (h Handler) CreateAccount(ctx context.Context, input *CreateAccountInput) (*CreateAccountOutput, error) {
	classReq, err := strconv.Atoi(input.Body.Class)
	if err != nil {
		return nil, huma.Error400BadRequest("invalid class value")
	}
	group, err := strconv.Atoi(input.Body.Group)
	if err != nil {
		return nil, huma.Error400BadRequest("invalid group value")
	}
	cmd := command.CreateAccountCmd{
		AccountId:                      uuid.New(),
		SobId:                          input.SobId,
		Title:                          input.Body.Title,
		LevelNumber:                    input.Body.LevelNumber,
		BalanceDirection:               input.Body.BalanceDirection,
		Class:                          classReq,
		Group:                          group,
		SuperiorRawAccountNumber:       input.Body.SuperiorRawAccountNumber,
		DimensionCategoryIds:           input.Body.DimensionCategoryIds,
		IsCashEquivalent:               input.Body.IsCashEquivalent,
		DefaultCashFlowItemIdForDebit:  input.Body.DefaultCashFlowItemIdForDebit,
		DefaultCashFlowItemIdForCredit: input.Body.DefaultCashFlowItemIdForCredit,
	}
	if err = h.app.Commands.CreateAccount.Handle(ctx, cmd); err != nil {
		return nil, err
	}
	createdAccount, err := h.app.Queries.AccountById.Handle(ctx, cmd.AccountId)
	if err != nil {
		return nil, err
	}
	return &CreateAccountOutput{Status: http.StatusCreated, Body: accountDTOToDetailVO(createdAccount)}, nil
}

type UpdateAccountInput struct {
	SobId     uuid.UUID `path:"sobId"`
	AccountId uuid.UUID `path:"accountId"`
	Body      UpdateAccountRequest
}

// UpdateAccount updates account settings and dimension bindings.
func (h Handler) UpdateAccount(ctx context.Context, input *UpdateAccountInput) (*struct{}, error) {
	group, err := strconv.Atoi(input.Body.Group)
	if err != nil {
		return nil, huma.Error400BadRequest("invalid group value")
	}
	cmd := command.UpdateAccountCmd{
		AccountId:                      input.AccountId,
		SobId:                          input.SobId,
		Title:                          input.Body.Title,
		LevelNumber:                    input.Body.LevelNumber,
		BalanceDirection:               input.Body.BalanceDirection,
		Group:                          group,
		DimensionCategoryIds:           input.Body.DimensionCategoryIds,
		IsCashEquivalent:               input.Body.IsCashEquivalent,
		DefaultCashFlowItemIdForDebit:  input.Body.DefaultCashFlowItemIdForDebit,
		DefaultCashFlowItemIdForCredit: input.Body.DefaultCashFlowItemIdForCredit,
		UpdateDefaultCashFlowItems:     input.Body.UpdateDefaultCashFlowItems,
	}
	if err = h.app.Commands.UpdateAccount.Handle(ctx, cmd); err != nil {
		return nil, err
	}
	return nil, nil
}

type DeleteAccountInput struct {
	SobId     uuid.UUID `path:"sobId"`
	AccountId uuid.UUID `path:"accountId"`
}

// DeleteAccount deletes an account from a SoB.
func (h Handler) DeleteAccount(ctx context.Context, input *DeleteAccountInput) (*struct{}, error) {
	cmd := command.DeleteAccountCmd{
		AccountId: input.AccountId,
		SobId:     input.SobId,
	}
	if err := h.app.Commands.DeleteAccount.Handle(ctx, cmd); err != nil {
		return nil, err
	}
	return nil, nil
}
