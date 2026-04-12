package errors

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// uniqueConstraintSlugs maps PostgreSQL unique-constraint names to the conflict
// slug that should be returned to the client.  Update this map whenever a new
// unique index is added to the schema — that is the only maintenance needed.
var uniqueConstraintSlugs = map[string]string{
	// sob
	"UQ_Sobs_Name": SlugSobDuplicateName,
	// general_ledger
	"UQ_Accounts_SobId_RawAccountNumber":        SlugAccountDuplicateNumber,
	"UQ_Periods_SobId_FiscalYear_PeriodNumber":  SlugPeriodDuplicateNumber,
	"UQ_Journals_SobId_PeriodId_DocumentNumber": SlugJournalDuplicateDocumentNumber,
	// dimension
	"UQ_DimCategories_SobId_Name":   SlugDimCategoryDuplicateName,
	"UQ_DimOptions_CategoryId_Name": SlugDimOptionDuplicateName,
	// report
	"UQ_Reports_SobId_PeriodId_Title": SlugReportDuplicateTitle,
	// numbering
	"UQ_IdentifierConfigurations_TargetBusinessObject_PropertyMatchers": SlugNumberingConfigDuplicate,
	"UQ_Identifiers_IdentifierConfigurationId_Identifier":               SlugNumberingIdentifierDuplicate,
}

// TranslateDBError converts a raw GORM / pgx error into a typed SlugErr so that
// the Gin middleware can return the correct HTTP status code and an i18n message.
//
// Call this at the repository boundary for every write operation instead of
// wrapping the error with fmt.Errorf.
//
//   - gorm.ErrRecordNotFound          → 404 NewNotFoundError
//   - unique_violation (23505)        → 409 NewConflictError  (slug from uniqueConstraintSlugs)
//   - foreign_key_violation (23503)   → 400 NewInternalError
//   - check_violation (23514)         → 400 NewInternalError
//   - not_null_violation (23502)      → 400 NewInternalError
//   - any other error                 → wrapped as-is (middleware returns 500)
func TranslateDBError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NewNotFoundError(SlugRecordNotFound)
	}

	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		switch pgErr.Code {
		case "23505": // unique_violation
			slug, ok := uniqueConstraintSlugs[pgErr.ConstraintName]
			if !ok {
				slug = SlugDuplicateEntry
			}
			return NewConflictError(slug)
		case "23503": // foreign_key_violation
			return NewInternalError(SlugForeignKeyViolation)
		case "23514": // check_violation
			return NewInternalError(SlugCheckConstraintViolation)
		case "23502": // not_null_violation
			return NewInternalError(SlugNotNullViolation)
		}
	}

	// Unknown infrastructure error — wrap without translation.
	// The middleware will fall through to HTTP 500.
	return fmt.Errorf("db error: %w", err)
}
