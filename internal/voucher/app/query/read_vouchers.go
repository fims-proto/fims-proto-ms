package query

import (
	"context"

	userQuery "github/fims-proto/fims-proto-ms/internal/user/app/query"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type void struct{}

var empty void

func toKeySlice[K comparable, V interface{}](set map[K]V) []K {
	keys := make([]K, len(set))
	i := 0
	for k := range set {
		keys[i] = k
		i++
	}
	return keys
}

type VouchersReadModel interface {
	ReadAllVouchers(ctx context.Context, sobId uuid.UUID, pageable data.Pageable) (data.Page[Voucher], error)
	ReadById(ctx context.Context, id uuid.UUID) (Voucher, error)
}

type ReadVouchersHandler struct {
	readModel      VouchersReadModel
	accountService AccountService
	userService    UserService
	ledgerService  LedgerService
}

func NewReadVouchersHandler(readModel VouchersReadModel, accountService AccountService, userService UserService, ledgerService LedgerService) ReadVouchersHandler {
	if readModel == nil {
		panic("nil readModel")
	}
	if accountService == nil {
		panic("nil account service")
	}
	if userService == nil {
		panic("nil user service")
	}
	if ledgerService == nil {
		panic("nil ledger service")
	}
	return ReadVouchersHandler{
		readModel:      readModel,
		accountService: accountService,
		userService:    userService,
		ledgerService:  ledgerService,
	}
}

func (h ReadVouchersHandler) HandleReadAll(ctx context.Context, sobId uuid.UUID, pageable data.Pageable) (data.Page[Voucher], error) {
	vouchersPage, err := h.readModel.ReadAllVouchers(ctx, sobId, pageable)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read all vouchers")
	}

	vouchers, err := h.enrichLineItemAccountNumber(ctx, vouchersPage.Content())
	if err != nil {
		return nil, errors.Wrap(err, "failed to enrich account number in vouchers")
	}

	vouchers, err = h.enrichUserName(ctx, vouchers)
	if err != nil {
		return nil, errors.Wrap(err, "failed to enrich user in vouchers")
	}

	vouchers, err = h.enrichPeriod(ctx, vouchers)
	if err != nil {
		return nil, errors.Wrap(err, "failed to enrich period in vouchers")
	}

	return data.NewPage(vouchers, pageable, vouchersPage.NumberOfElements())
}

func (h ReadVouchersHandler) HandleReadById(ctx context.Context, id uuid.UUID) (Voucher, error) {
	voucher, err := h.readModel.ReadById(ctx, id)
	if err != nil {
		return Voucher{}, errors.Wrap(err, "failed to read voucher")
	}

	singletonList, err := h.enrichLineItemAccountNumber(ctx, []Voucher{voucher})
	if err != nil {
		return Voucher{}, errors.Wrap(err, "failed to enrich account number in voucher")
	}

	singletonList, err = h.enrichUserName(ctx, singletonList)
	if err != nil {
		return Voucher{}, errors.Wrap(err, "failed to enrich user in voucher")
	}

	singletonList, err = h.enrichPeriod(ctx, singletonList)
	if err != nil {
		return Voucher{}, errors.Wrap(err, "failed to enrich period in voucher")
	}

	return singletonList[0], nil
}

func (h ReadVouchersHandler) enrichLineItemAccountNumber(ctx context.Context, vouchers []Voucher) ([]Voucher, error) {
	accountSet := make(map[uuid.UUID]void)
	for _, voucher := range vouchers {
		for _, item := range voucher.LineItems {
			accountSet[item.AccountId] = empty
		}
	}

	accounts, err := h.accountService.ReadAccountsByIds(ctx, toKeySlice[uuid.UUID, void](accountSet))
	if err != nil {
		return nil, errors.Wrap(err, "failed to read accounts by Ids")
	}

	for i := range vouchers {
		for j := range vouchers[i].LineItems {
			account, ok := accounts[vouchers[i].LineItems[j].AccountId]
			if !ok {
				return nil, errors.Errorf("account not found by id: %s", vouchers[i].LineItems[j].AccountId)
			}
			vouchers[i].LineItems[j].AccountNumber = account.AccountNumber
		}
	}

	return vouchers, nil
}

func (h ReadVouchersHandler) enrichUserName(ctx context.Context, vouchers []Voucher) ([]Voucher, error) {
	userSet := make(map[uuid.UUID]void)
	for _, voucher := range vouchers {
		if voucher.Creator.Id != uuid.Nil {
			userSet[voucher.Creator.Id] = empty
		}
		if voucher.Reviewer.Id != uuid.Nil {
			userSet[voucher.Reviewer.Id] = empty
		}
		if voucher.Auditor.Id != uuid.Nil {
			userSet[voucher.Auditor.Id] = empty
		}
	}

	users, err := h.userService.ReadUserByIds(ctx, toKeySlice[uuid.UUID, void](userSet))
	if err != nil {
		return nil, errors.Wrap(err, "failed to read users by Ids")
	}

	convertUser := func(user User, users map[uuid.UUID]userQuery.User) User {
		return User{
			Id:     user.Id,
			Traits: users[user.Id].Traits,
		}
	}

	for i := range vouchers {
		if vouchers[i].Creator.Id != uuid.Nil {
			vouchers[i].Creator = convertUser(vouchers[i].Creator, users)
		}
		if vouchers[i].Reviewer.Id != uuid.Nil {
			vouchers[i].Reviewer = convertUser(vouchers[i].Reviewer, users)
		}
		if vouchers[i].Auditor.Id != uuid.Nil {
			vouchers[i].Auditor = convertUser(vouchers[i].Auditor, users)
		}
	}

	return vouchers, nil
}

func (h ReadVouchersHandler) enrichPeriod(ctx context.Context, vouchers []Voucher) ([]Voucher, error) {
	periodSet := make(map[uuid.UUID]void)
	for _, voucher := range vouchers {
		periodSet[voucher.Period.Id] = empty
	}

	periods, err := h.ledgerService.ReadPeriodsByIds(ctx, toKeySlice[uuid.UUID, void](periodSet))
	if err != nil {
		return nil, errors.Wrap(err, "failed to read periods by Ids")
	}

	for i := range vouchers {
		period, ok := periods[vouchers[i].Period.Id]
		if !ok {
			return nil, errors.Errorf("period not found by id: %s", vouchers[i].Period.Id)
		}
		vouchers[i].Period = Period{
			Id:            period.Id,
			FinancialYear: period.FinancialYear,
			Number:        period.Number,
			OpeningTime:   period.OpeningTime,
			EndingTime:    period.EndingTime,
			IsClosed:      period.IsClosed,
		}
	}

	return vouchers, nil
}
