package adapter

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/ledger/domain"
	"sync"

	"github.com/pkg/errors"
)

type LedgerMemoryRepository struct {
	lock *sync.RWMutex
	data map[string]domain.Ledger
}

func NewLedgerMemoryRepository() LedgerMemoryRepository {
	return LedgerMemoryRepository{
		lock: &sync.RWMutex{},
		data: make(map[string]domain.Ledger),
	}
}

func (r LedgerMemoryRepository) AddLedger(ctx context.Context, l *domain.Ledger) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	_, ok := r.data[l.Number()]
	if ok {
		return errors.Errorf("ledger exists with number %s", l.Number())
	}
	r.data[l.Number()] = *l
	return nil
}

func (r LedgerMemoryRepository) UpdateLedgers(
	ctx context.Context,
	ledgerNumbers []string,
	updateFn func(ledgers []*domain.Ledger) ([]*domain.Ledger, error),
) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	// fetch entries from db
	var ledgers []*domain.Ledger
	for _, inputNum := range ledgerNumbers {
		l, ok := r.data[inputNum]
		if !ok {
			return errors.Errorf("ledger number %s not exists", inputNum)
		}
		ledgers = append(ledgers, &l)
	}

	// call updateFn
	afterUpdateLedgers, err := updateFn(ledgers)
	if err != nil {
		return errors.Wrap(err, "ledger list update failed")
	}

	// update db
	for _, l := range afterUpdateLedgers {
		r.data[l.Number()] = *l
	}
	return nil
}

func (r LedgerMemoryRepository) Dataload(ctx context.Context, ledgers []*domain.Ledger) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	// clear map
	for i := range r.data {
		delete(r.data, i)
	}

	for _, ledger := range ledgers {
		r.data[ledger.Number()] = *ledger
	}
	return nil
}

// TODO: remove, test purpose
func (r LedgerMemoryRepository) AllLedgers(ctx context.Context) ([]domain.Ledger, error) {
	allLedgers := []domain.Ledger{}
	for _, v := range r.data {
		allLedgers = append(allLedgers, v)
	}
	return allLedgers, nil
}
