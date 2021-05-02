package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApp_HandleAuditVoucher(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		cancel      bool
		constructor func(t *testing.T) *domain.Voucher
		auditor     string
	}{
		{
			"audit_success",
			false,
			func(t *testing.T) *domain.Voucher {
				return createVoucherForAuditTest(t, "")
			},
			"auditor1",
		},
		{
			"cancel_audit_success",
			true,
			func(t *testing.T) *domain.Voucher {
				return createVoucherForAuditTest(t, "auditor1")
			},
			"auditor1",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			v := test.constructor(t)
			repoMock := newVoucherRepoMock()
			repoMock.vouchers = map[uuid.UUID]domain.Voucher{
				v.UUID(): *v,
			}

			handler := NewAuditVoucherHandler(repoMock)

			if test.cancel {
				err := handler.HandleCancel(context.Background(), AuditVoucherCmd{
					VoucherUUID: v.UUID(),
					Auditor:     test.auditor,
				})
				assert.NoError(t, err)
				assert.False(t, repoMock.vouchers[v.UUID()].IsAudited())
			} else {
				err := handler.Handle(context.Background(), AuditVoucherCmd{
					VoucherUUID: v.UUID(),
					Auditor:     test.auditor,
				})
				assert.NoError(t, err)
				assert.True(t, repoMock.vouchers[v.UUID()].IsAudited())
			}
		})
	}
}

func createVoucherForAuditTest(t *testing.T, auditor string) *domain.Voucher {
	v, err := domain.NewVoucher(uuid.New(), domain.GeneralVoucher, "1", time.Now(), 0, prepareBalancedItems(), "")
	require.NoError(t, err)
	if auditor != "" {
		err := v.Audit(auditor)
		require.NoError(t, err)
	}
	return v
}
