package query

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/errors"

	"github.com/google/uuid"
)

type periodRangeValidator struct {
	readModel GeneralLedgerReadModel
}

// NewPeriodRangeValidator creates a validator that uses the repository directly
// This avoids circular dependency between service and query packages
func newPeriodRangeValidator(readModel GeneralLedgerReadModel) periodRangeValidator {
	if readModel == nil {
		panic("nil read model")
	}
	return periodRangeValidator{readModel: readModel}
}

func (v periodRangeValidator) validate(ctx context.Context, sobId uuid.UUID, fromPeriod, toPeriod string) (int, int, int, int, error) {
	// Parse period format (e.g., "2026-01" -> fiscalYear=2026, periodNumber=1)
	fromFiscalYear, fromPeriodNumber, err := parsePeriodString(fromPeriod)
	if err != nil {
		return 0, 0, 0, 0, errors.NewSlugError("invalid-period-format", fromPeriod)
	}
	toFiscalYear, toPeriodNumber, err := parsePeriodString(toPeriod)
	if err != nil {
		return 0, 0, 0, 0, errors.NewSlugError("invalid-period-format", toPeriod)
	}

	// Validate from <= to
	if fromFiscalYear > toFiscalYear || (fromFiscalYear == toFiscalYear && fromPeriodNumber > toPeriodNumber) {
		return 0, 0, 0, 0, errors.NewSlugError("period-range-invalid", fromPeriod, toPeriod)
	}

	// Check period continuity using optimized SQL query (only queries from-to range)
	if err := v.readModel.CheckPeriodContinuity(ctx, sobId, fromFiscalYear, fromPeriodNumber, toFiscalYear, toPeriodNumber); err != nil {
		return 0, 0, 0, 0, errors.NewSlugError("period-range-not-continuous", fromPeriod, toPeriod)
	}

	return fromFiscalYear, fromPeriodNumber, toFiscalYear, toPeriodNumber, nil
}

func parsePeriodString(s string) (int, int, error) {
	if s == "" {
		return 0, 0, errors.NewSlugError("invalid-period-format", s)
	}

	// Parse "YYYY-MM" format
	var fy, pn int
	n, err := fmt.Sscanf(s, "%d-%d", &fy, &pn)
	if err != nil || n != 2 {
		return 0, 0, errors.NewSlugError("invalid-period-format", s)
	}
	if pn < 1 || pn > 12 {
		return 0, 0, errors.NewSlugError("invalid-period-format", s)
	}
	return fy, pn, nil
}
