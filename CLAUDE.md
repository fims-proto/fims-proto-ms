# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Financial Information Management System (FIMS) - A multi-tenant accounting system built with Go, using hexagonal architecture with CQRS patterns. Core domains: Set of Books (SoB/账套), General Ledger (总账), Dimension (维度), Reports, Numbering, and User management.

## Development Commands

### Running Locally

```bash
# 1. Start infrastructure (in fims-dev-environment repo)
docker compose up --build  # Postgres, Ory Kratos, Oathkeeper

# 2. Start backend
go run cmd/main.go  # Port 5002

# 3. Gateway accessible at http://127.0.0.1:4455/
```

### Testing

```bash
make test          # Run all tests with -count=1
go test ./...      # Alternative without make

# Run single package tests
go test ./internal/general_ledger/domain/journal/... -count=1

# Run specific test
go test ./internal/general_ledger/domain/journal/... -run TestJournalPost -count=1
```

### Code Quality

```bash
make fmt           # Format code with swag fmt + gofumpt
make lint          # Run golangci-lint

# Manual formatting
gofumpt -l -w internal/ cmd/
```

### Swagger Documentation

```bash
make swag          # Regenerate from @Tags/@Summary comments
# Visit: http://127.0.0.1:4455/fims/swagger/index.html
```

### PEG Parser Generation

```bash
make peg           # Regenerate filterable query parser from .peg grammar
```

## Architecture

### Hexagonal Architecture (Ports & Adapters)

Each domain module follows this structure:

```
internal/{domain}/
├── domain/        # Business logic, value objects, aggregates
│                  # Domain methods enforce business rules (e.g., journal.Post(), journal.Audit())
├── app/           # Application services organized as Commands (writes) and Queries (reads)
│   ├── command/   # Write operations with handler pattern
│   ├── query/     # Read operations and read models
│   └── app.go     # Application composition with Inject() method
├── port/          # Interfaces to external world
│   ├── public/http/       # External API endpoints
│   ├── private/http/      # Internal/admin API endpoints
│   └── private/intraprocess/  # Cross-module communication interfaces
└── adapter/       # Infrastructure implementations
    ├── db/        # Repository implementations
    └── {module}/  # Cross-module adapters (e.g., numbering, sob)
```

Core domains: `general_ledger`, `sob`, `report`, `numbering`, `user`, `dimension`

### Dependency Injection

Applications are initialized in `cmd/main.go` with manual DI pattern:

1. Create repositories from datasource
2. Create empty Applications: `app := module.NewApplication()`
3. Create intraprocess interfaces for cross-module calls
4. Inject dependencies: `app.Inject(repo, readModel, externalServices...)`

**Critical**: Cross-module calls MUST use intraprocess adapters, never direct repository access.

### DataSource Abstraction

Always use `datasource.DataSource` interface, never raw GORM:

```go
type DataSource interface {
GetConnection(ctx context.Context) *gorm.DB
EnableTransaction(ctx context.Context, transactionalFn func(txCtx context.Context) error) error
}
```

Two implementations:

- `dedicated-datasource/` - Single database (production ready)
- `multitenant-datasource/` - Multi-tenant routing via subdomain (**stub only**, not production-ready; `app.multiTenancy` defaults to false)

Subdomain resolution: `ResolveSubdomain()` middleware extracts tenant from URL (`tenant.domain.com` → `tenant`)

## Critical Development Patterns

### Transaction Management

**ALWAYS** use `repo.EnableTx()` for write operations. Never manage transactions in HTTP handlers.

```go
return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
    // All repository operations use txCtx
    return h.repo.UpdateJournal(txCtx, id, func(j *journal.Journal) (*journal.Journal, error) {
        return j, j.Post(poster)
    })
})
```

### Repository Update Pattern

Domain state changes use update callbacks with domain methods:

```go
repo.UpdateJournal(ctx, journalId, func(j *journal.Journal) (*journal.Journal, error) {
    if err := j.Post(poster); err != nil {  // Domain method enforces rules
        return nil, err
    }
    return j, nil
})
```

Never mutate domain objects directly - always use their methods (e.g., `j.Post()`, `j.Audit()`, `j.Review()`).

### Read Model vs Repository

**Rule**: Where a read method lives depends on what it serves — not just whether it reads data.

- **`GeneralLedgerReadModel`** — all read methods that serve **query handlers** (GET endpoints). These return query DTOs and live in `internal/general_ledger/app/query/read_model.go`.
- **`domain.Repository`** — read methods that exist solely to **support write commands**. For example, `ExistsJournalsNotPostedInPeriod` is used by `ClosePeriodHandler` to gate a write operation. These return domain objects.

When adding a new read-only query, always add it to `GeneralLedgerReadModel`. Only add to `domain.Repository` if the read is needed inside a command handler.

### Error Handling with i18n

Use typed slug-based errors for business logic. Each constructor sets the HTTP status automatically:

```go
return commonErrors.NewInvalidInputError(commonErrors.SlugJournalPostNotAudited) // → HTTP 400
return commonErrors.NewNotFoundError(commonErrors.SlugJournalNotFound)           // → HTTP 404
return commonErrors.NewConflictError(commonErrors.SlugJournalDuplicateDocumentNumber) // → HTTP 409
return commonErrors.NewInternalError(commonErrors.SlugRecordNotFound)            // → HTTP 500
```

All slug string constants live in `internal/common/errors/slugs.go` — **never use inline string literals**.

When adding new slugs:

1. Add a constant to `internal/common/errors/slugs.go`: `SlugMyNewSlug = "module-operation-reason"`
2. Reference it with the appropriate typed constructor: `commonErrors.NewInvalidInputError(commonErrors.SlugMyNewSlug, args...)`
3. Add to `i18n/zh-CN.json`: `"module-operation-reason": "本地化消息 {{.A}}"`
4. Middleware auto-maps slugs to localized responses

See `internal/common/errors/slug_err.go`, `slugs.go`, and `internal/common/errors/gin_middleware.go`

## Domain-Specific Knowledge

### General Ledger Components

- **account** - Chart of accounts with hierarchical structure
- **period** - Accounting periods (monthly/quarterly/yearly)
- **ledger** - Account balance tracking per period
- **journal** - Journal entries with journal lines

### Signed Amount Model

**Important**: The codebase uses a **signed amount model** instead of separate debit/credit fields.

**Amount Convention**:
- **Positive** amounts represent **debits**
- **Negative** amounts represent **credits**

**Domain Model Changes**:

1. **JournalLine**: Single `amount` field (signed) instead of separate `debitAmount`/`creditAmount`
2. **Journal**: Single `amount` field (transaction total = sum of positive/debit amounts)
3. **LedgerEntry**: Single `amount` field (signed)
4. **Ledger**:
   - `openingAmount` - Signed opening balance
   - `periodAmount` - Signed net movement
   - `periodDebit/periodCredit` - Positive values (retained for query performance)
   - `endingAmount` - Signed ending balance

**Balance Calculation**:
```
endingAmount = openingAmount + periodAmount
```

**Trial Balance Validation**:
- Sum of all signed journal line amounts must equal zero
- Instead of checking `totalDebits == totalCredits`

**Posting Logic** (`ledger.update_balance`):
```go
periodAmount += amount          // Signed addition
if amount.IsPositive() {
    periodDebit += amount      // Positive value
} else {
    periodCredit += amount.Abs() // Positive absolute value
}
```

### Journal Lifecycle

Workflow: Create → Review → Audit → Post (登账) → affects Ledgers

(Note: `CancelReview` and `CancelAudit` commands exist to revert journal back to earlier states)

**Delete path**: Only `TypeClosing` and `TypeYearlyClosing` system journals can be deleted (`DELETE /sob/{sobId}/journal/{journalId}`). Deletion reverses the journal's ledger impact (negates all posted amounts — exact inverse of posting). Regular user-created journals cannot be deleted. See `internal/general_ledger/app/command/delete_system_journal.go`.

Business rules enforced in `internal/general_ledger/domain/journal/`:

- Journal must be reviewed AND audited before posting
- Creator ≠ Reviewer ≠ Auditor (segregation of duties)
- Cannot modify after audit/review
- Period must be current and open for posting
- Posting updates `ledger` balances

### Posting Mechanics

File: `internal/general_ledger/app/command/post_journal.go`

Process:

1. Call `j.Post(poster)` to validate and mark journal
2. Build posting records for all journal lines + their superior accounts (each with signed amount)
3. Merge identical account records (sum signed amounts)
4. Batch update: `UpdateLedgersByPeriodAndAccountIds()`

**Critical**: Maintain merge logic and batch updates to prevent inconsistencies.

### Report Generation (High Risk Area)

File: `internal/report/domain/generator/generator.go`

Reports use two data source types:

- **Sum** - Aggregates ledger balances by account filters
- **Formulas** - Four formula rules: `Net`, `Debit`, `Credit`, `Transaction`

**Why high risk:**
- `ledgersCache` is a shared in-memory map across all formula evaluations — stale or incorrect entries produce wrong financial figures with no runtime error
- Formula rules interact with `balance_direction` and `data_source` in non-obvious ways; bugs only surface in output numbers, not compile time
- Balance sheet (`report-balanceSheet-imbalance`) and income statement (`report-incomeStatement-profitMismatch`) validations only catch end-to-end failures, not intermediate calculation errors

**Any formula or aggregation changes MUST have unit tests.**

### Numbering Service

Generates sequential identifiers for journals: `numberingService.GenerateIdentifier(ctx, periodId)`

Configurations define patterns with auto-increment counters per period/type.

### Dimension Domain

`internal/dimension/` — Manages accounting dimensions (cost centers, departments, tags, etc.)

- **category** — Dimension categories (e.g., "Department", "Project"), scoped per SoB
- **option** — Dimension values within a category (e.g., "Engineering", "Marketing")

#### Account → Dimension Category Binding

Each GL account stores `dimensionCategoryIds []uuid.UUID` — the set of dimension categories that are **required** when tagging a journal line against that account. Persisted as a join table (`accountDimensionCategoryPO`) in `internal/general_ledger/adapter/db/types.go`.

#### Journal Line → Dimension Option Association

Each `JournalLine` stores `dimensionOptionIds []uuid.UUID` — the selected dimension option IDs for that line. Options are stored by ID only; full objects are resolved at query time (see Enrichment below).

#### Validation Flow (Write Path)

In `prepareJournalLines` (`internal/general_ledger/app/command/common_functions.go`), for each journal line:

```go
dimensionService.ValidateOptions(ctx, a.DimensionCategoryIds(), item.DimensionOptionIds)
```

`ValidateOptions` enforces (`internal/dimension/app/query/validate_options.go`):
1. Every required category (from the account) must have exactly one option provided.
2. No option may belong to a category not bound to the account.
3. No duplicate options for the same category.
4. All provided option IDs must exist.

#### Enrichment (Read Path)

On detail queries, raw IDs are hydrated to full objects by `internal/general_ledger/app/query/enricher.go`:

- `enrichJournalLineDimensionOptions` — batch-fetches all option objects for all lines in a journal, maps them to `JournalLine.DimensionOptions []DimensionOption`
- `enrichAccountDimensionCategories` — fetches category objects for `Account.DimensionCategories []DimensionCategory`

These fields are **not stored in the DB** — populated on read only. The raw `DimensionOptionIds` / `DimensionCategoryIds` fields in query types are marked *"internal: used by enricher, not exposed in HTTP response"*.

#### LedgerDimensionSummary Query

`internal/general_ledger/app/query/ledger_dimension_summary.go` — aggregates ledger entries by dimension option within a period range for a given account and dimension category. Returns `[]LedgerDimensionSummaryItem{DimensionOptionId, DimensionOptionName, TotalAmount}`.

#### Cross-Module Wiring

- `DimensionService` interface: `internal/general_ledger/app/service/services.go`
- Intraprocess adapter: `internal/dimension/port/private/intraprocess/dimension.go` (`DimensionInterface`)
- Injected in `cmd/main.go` as `generalLedgerDimensionAdapter`

### Set of Books (SoB)

A SoB (账套) represents a complete accounting entity with its own:

- Chart of accounts
- Accounting periods
- Journals and ledgers
- Reporting configurations

Multi-tenancy is achieved through SoB isolation.

## Common Utilities

### Filterable Query Language

`internal/common/data/filterable/` - PEG-based query parser for filtering API results

Grammar defined in `filterable.peg`. Generate parser with `make peg`.

### Pagination & Sorting

`internal/common/data/pageable/` - Pagination helpers with Gin integration

`internal/common/data/field/` - Dynamic field selection and GORM mapping

## Testing Strategy

### Unit Tests

Focus on:

- Domain logic: journal state transitions, business rule enforcement
- Report generator: formula calculations
- Repository updates with transactions

### Integration Tests

Bruno API tests in `bruno_collection/fims local/`: journal creation → review → audit → post → verify ledgers.

Use these flows as integration test templates.

## Configuration

- `config/application-dev.yaml` - Development config
- `config/application-production.yaml` - Production config
- `internal/common/config/default.go` - Config struct and defaults

## Key Files Reference

**Entry Point**: `cmd/main.go` - DI setup, router configuration, middleware registration

**Core Interfaces**:

- `internal/common/datasource/datasource.go` - Database abstraction
- `internal/*/domain/repository.go` - Domain repository contracts
- `internal/*/app/app.go` - Application service composition

**Domain Examples**:

- `internal/general_ledger/domain/journal/` - Journal aggregate with business rules
- `internal/general_ledger/domain/ledger/` - Ledger balance tracking
- `internal/sob/domain/sob/` - Set of Books configuration

**Critical Implementations**:

- `internal/general_ledger/app/command/post_journal.go` - Posting logic
- `internal/report/domain/generator/generator.go` - Report generation
- `internal/common/errors/gin_middleware.go` - Error translation

**Dev-only**: `internal/devops/jwt_handler.go` — JWT utility registered at `/devops/` route only when `profile` starts with `"dev"` (controlled via `cmd/main.go`)

## Common Pitfalls

❌ Managing transactions in HTTP handlers
❌ Mutating domain objects without their methods
❌ Forgetting i18n entries for new slug errors
❌ Direct cross-module repository calls
❌ Modifying report formulas without tests
❌ Skipping `swag init` after API changes

## Swagger Annotations

HTTP handlers require swagger annotations:

```go
// @Tags        journals
// @Summary     Post journal to ledgers
// @Param       sobId path string true "Sob ID"
// @Router      /sob/{sobId}/journal/{journalId}/post [patch]
```

Main swagger config in `api/api.go`

## Code Style & Approach section

When asked to plan or implement changes, start with the simplest approach that follows existing patterns in the codebase. Do NOT over-engineer, create wrapper components, or introduce new abstractions unless explicitly requested. Look for existing patterns first and reuse them.

When asked for backend-only or frontend-only analysis, stay strictly within that boundary. Do not include suggestions or changes for the other side unless explicitly asked.

Before implementing a fix for a bug, create a brief plan and confirm the approach. Do not jump straight into coding a fix without understanding the root cause first. When debugging, avoid rapid-fire guessing — instead, methodically trace the data flow.
