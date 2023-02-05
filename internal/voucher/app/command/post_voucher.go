package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/voucher/app/service"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/voucher"

	"github/fims-proto/fims-proto-ms/internal/voucher/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type PostVoucherCmd struct {
	VoucherId uuid.UUID
	Poster    uuid.UUID
}

type PostVoucherHandler struct {
	repo           domain.Repository
	accountService service.AccountService
}

func NewPostVoucherHandler(repo domain.Repository, accountService service.AccountService) PostVoucherHandler {
	if repo == nil {
		panic("nil repo")
	}

	if accountService == nil {
		panic("nil account service")
	}

	return PostVoucherHandler{
		repo:           repo,
		accountService: accountService,
	}
}

func (h PostVoucherHandler) Handle(ctx context.Context, cmd PostVoucherCmd) error {
	return h.repo.UpdateVoucher(
		ctx,
		cmd.VoucherId,
		func(j *voucher.Voucher) (*voucher.Voucher, error) {
			if err := j.Post(cmd.Poster); err != nil {
				return nil, err
			}

			if err := h.accountService.PostVoucher(ctx, *j); err != nil {
				return nil, errors.Wrap(err, "failed to post voucher to accounts")
			}

			return j, nil
		},
	)
}
