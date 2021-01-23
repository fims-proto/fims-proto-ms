package command

import "context"

type AccountService interface {
	ValidateExistence(ctx context.Context, accNumbers []string) error
}
