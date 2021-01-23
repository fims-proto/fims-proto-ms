package query

import "context"

type ValidateAccountsModel interface {
	ValidateExistence(ctx context.Context, accNumbers []string) error
}

type ValidateAccountsHandler struct {
	validateModel ValidateAccountsModel
}

func NewValidateAccountsHandler(validateModel ValidateAccountsModel) ValidateAccountsHandler {
	if validateModel == nil {
		panic("nil validateModel")
	}
	return ValidateAccountsHandler{validateModel: validateModel}
}

func (h ValidateAccountsHandler) HandleValidateExistence(ctx context.Context, accNumbers []string) error {
	return h.validateModel.ValidateExistence(ctx, accNumbers)
}
