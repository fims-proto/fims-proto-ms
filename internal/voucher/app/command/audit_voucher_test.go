package command

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/lineitem"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/voucher"
	"testing"
	"time"
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
			deps := newAuditMockDeps()
			deps.repository.vouchers = map[string]voucher.Voucher{
				v.UUID(): *v,
			}

			err := deps.handler.Handle(context.Background(), AuditVoucherCmd{
				VoucherUUID: v.UUID(),
				AuditorUUID: test.auditorUUID,
			})
			assert.NoError(t, err)

			assert.True(t, deps.repository.vouchers[v.UUID()].IsAudited())
		})
	}
}

func createVoucherForAuditTest(t *testing.T, auditorUUID string) *voucher.Voucher {
	v, err := voucher.NewVoucher("test_uuid", 1, time.Now(), 0, []lineitem.LineItem{}, "")
	require.NoError(t, err)
	if auditorUUID != "" {
		err := v.Audit(auditorUUID)
		require.NoError(t, err)
	}
	return v
}

type auditMockDeps struct {
	repository *auditRepoMock
	handler    AuditVoucherHandler
}

func newAuditMockDeps() auditMockDeps {
	repository := &auditRepoMock{}
	return auditMockDeps{
		repository: repository,
		handler:    AuditVoucherHandler{repository},
	}
}

type auditRepoMock struct {
	vouchers map[string]voucher.Voucher
}

func (r *auditRepoMock) AddVoucher(ctx context.Context, v *voucher.Voucher) error {
	panic("implement me")
}

func (r *auditRepoMock) UpdateVoucher(
	ctx context.Context,
	voucherUUID string,
	updateFn func(v *voucher.Voucher) (*voucher.Voucher, error),
) error {
	v, ok := r.vouchers[voucherUUID]
	if !ok {
		return voucher.NotFoundError{VoucherUUID: voucherUUID}
	}

	updatedVoucher, err := updateFn(&v)
	if err != nil {
		return err
	}

	r.vouchers[voucherUUID] = *updatedVoucher
	return nil
}
