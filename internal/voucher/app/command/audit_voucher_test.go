package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/voucher"
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
		constructor func(t *testing.T) *voucher.Voucher
		auditorUUID string
	}{
		{
			"normal_success",
			func(t *testing.T) *voucher.Voucher {
				return createVoucherForAuditTest(t, "")
			},
			"aud1_uuid",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			v := test.constructor(t)
			repoMock := newVoucherRepoMock()
			repoMock.vouchers = map[uuid.UUID]voucher.Voucher{
				v.UUID(): *v,
			}

			handler := NewAuditVoucherHandler(repoMock)

			err := handler.Handle(context.Background(), AuditVoucherCmd{
				VoucherUUID: v.UUID(),
				Auditor:     test.auditorUUID,
			})
			assert.NoError(t, err)

			assert.True(t, repoMock.vouchers[v.UUID()].IsAudited())
		})
	}
}

func createVoucherForAuditTest(t *testing.T, auditorUUID string) *voucher.Voucher {
	v, err := voucher.NewVoucher(uuid.New(), "1", time.Now(), 0, prepareBalancedItems(), "")
	require.NoError(t, err)
	if auditorUUID != "" {
		err := v.Audit(auditorUUID)
		require.NoError(t, err)
	}
	return v
}
