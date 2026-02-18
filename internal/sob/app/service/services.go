package service

import (
	"context"

	"github.com/google/uuid"
)

type GeneralLedgerService interface {
	InitializeGeneralLedger(ctx context.Context, sobId uuid.UUID) error
}

type ReportService interface {
	InitializeReport(ctx context.Context, sobId uuid.UUID) error
}
