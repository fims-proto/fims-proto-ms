package command

import (
	"context"
	"testing"
	"time"

	"github/fims-proto/fims-proto-ms/internal/voucher/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var userA = uuid.New()

func TestApp_HandleAuditVoucher(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		cancel      bool
		constructor func(t *testing.T) *domain.Voucher
		auditor     uuid.UUID
	}{
		{
			"audit_success",
			false,
			func(t *testing.T) *domain.Voucher {
				return createVoucherForAuditTest(t, uuid.Nil)
			},
			userA,
		},
		{
			"cancel_audit_success",
			true,
			func(t *testing.T) *domain.Voucher {
				return createVoucherForAuditTest(t, userA)
			},
			userA,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			v := test.constructor(t)
			repoMock := newVoucherRepoMock()
			repoMock.vouchers = map[uuid.UUID]domain.Voucher{
				v.Id(): *v,
			}

			handler := NewAuditVoucherHandler(repoMock)

			if test.cancel {
				err := handler.HandleCancel(context.Background(), AuditVoucherCmd{
					VoucherUUID: v.Id(),
					Auditor:     test.auditor,
				})
				assert.NoError(t, err)
				assert.False(t, repoMock.vouchers[v.Id()].IsAudited())
			} else {
				err := handler.Handle(context.Background(), AuditVoucherCmd{
					VoucherUUID: v.Id(),
					Auditor:     test.auditor,
				})
				assert.NoError(t, err)
				assert.True(t, repoMock.vouchers[v.Id()].IsAudited())
			}
		})
	}
}

func createVoucherForAuditTest(t *testing.T, auditor uuid.UUID) *domain.Voucher {
	v, err := domain.NewVoucher(uuid.New(), uuid.New(), "GENERAL_VOUCHER", "1", 0, prepareBalancedItems(), uuid.New(), uuid.Nil, uuid.Nil, false, false, false, time.Now())
	require.NoError(t, err)
	if auditor != uuid.Nil {
		err := v.Audit(auditor)
		require.NoError(t, err)
	}
	return v
}
