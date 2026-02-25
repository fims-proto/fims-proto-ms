# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Financial Information Management System (FIMS) - A multi-tenant accounting system built with Go, using hexagonal architecture with CQRS patterns. Core domains: Set of Books (SoB/账套), General Ledger (总账), Reports, Numbering, and User management.

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
go test ./internal/general_ledger/domain/voucher/... -count=1

# Run specific test
go test ./internal/general_ledger/domain/voucher/... -run TestVoucherPost -count=1
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
│                  # Domain methods enforce business rules (e.g., voucher.Post(), voucher.Audit())
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

Core domains: `general_ledger`, `sob`, `report`, `numbering`, `user`

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
- `multitenant-datasource/` - Multi-tenant routing via subdomain (stub)

Subdomain resolution: `ResolveSubdomain()` middleware extracts tenant from URL (`tenant.domain.com` → `tenant`)

## Critical Development Patterns

### Transaction Management

**ALWAYS** use `repo.EnableTx()` for write operations. Never manage transactions in HTTP handlers.

```go
return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
    // All repository operations use txCtx
    return h.repo.UpdateVoucher(txCtx, id, func(v *voucher.Voucher) (*voucher.Voucher, error) {
        return v, v.Post(poster)
    })
})
```

### Repository Update Pattern

Domain state changes use update callbacks with domain methods:

```go
repo.UpdateVoucher(ctx, voucherId, func(v *voucher.Voucher) (*voucher.Voucher, error) {
    if err := v.Post(poster); err != nil {  // Domain method enforces rules
        return nil, err
    }
    return v, nil
})
```

Never mutate domain objects directly - always use their methods (e.g., `v.Post()`, `v.Audit()`, `v.Review()`).

### Error Handling with i18n

Use slug-based errors for business logic:

```go
return errors.NewSlugError("voucher-post-notAudited")
```

When adding new slugs:

1. Use in code: `errors.NewSlugError("module-operation-reason", args...)`
2. Add to `i18n/zh-CN.json`: `"module-operation-reason": "本地化消息 {{.A}}"`
3. Middleware auto-maps slugs to localized responses

See `internal/common/errors/slug_err.go` and `internal/common/errors/gin_middleware.go`

## Domain-Specific Knowledge

### General Ledger Components

- **account** - Chart of accounts with hierarchical structure
- **period** - Accounting periods (monthly/quarterly/yearly)
- **ledger** - Account balance tracking per period
- **auxiliary_ledger** - Detailed balance tracking with auxiliary dimensions
- **voucher** - Journal entries with line items

### Signed Amount Model

**Important**: The codebase uses a **signed amount model** instead of separate debit/credit fields.

**Amount Convention**:
- **Positive** amounts represent **debits**
- **Negative** amounts represent **credits**

**Domain Model Changes**:

1. **LineItem**: Single `amount` field (signed) instead of separate `debitAmount`/`creditAmount`
2. **Voucher**: Single `amount` field (transaction total = sum of positive/debit amounts)
3. **LedgerEntry**: Single `amount` field (signed)
4. **Ledger/AuxiliaryLedger**:
   - `openingAmount` - Signed opening balance
   - `periodAmount` - Signed net movement
   - `periodDebit/periodCredit` - Positive values (retained for query performance)
   - `endingAmount` - Signed ending balance

**Balance Calculation**:
```
endingAmount = openingAmount + periodAmount
```

**Trial Balance Validation**:
- Sum of all signed line item amounts must equal zero
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

### Voucher Lifecycle

Workflow: Create → Review → Audit → Post (登账) → affects Ledgers

Business rules enforced in `internal/general_ledger/domain/voucher/`:

- Voucher must be reviewed AND audited before posting
- Creator ≠ Reviewer ≠ Auditor (segregation of duties)
- Cannot modify after audit/review
- Period must be current and open for posting
- Posting updates both `ledger` and `auxiliary_ledger` balances

### Posting Mechanics

File: `internal/general_ledger/app/command/post_voucher.go`

Process:

1. Call `v.Post(poster)` to validate and mark voucher
2. Build posting records for all line items + their superior accounts (each with signed amount)
3. Merge identical account records (sum signed amounts)
4. Batch update: `UpdateLedgersByPeriodAndAccountIds()` and `UpdateAuxiliaryLedgersByPeriodAndAccountIds()`

**Critical**: Maintain merge logic and batch updates to prevent inconsistencies.

### Report Generation (High Risk Area)

File: `internal/report/domain/generator/generator.go`

Reports use two data source types:

- **Sum** - Aggregates ledger balances by account filters
- **Formulas** - Four formula rules: `Net`, `Debit`, `Credit`, `Transaction`

`ledgersCache` reduces DB reads during generation. Changes require extensive testing.

**Any formula or aggregation changes MUST have unit tests.**

### Numbering Service

Generates sequential identifiers for vouchers: `numberingService.GenerateIdentifier(ctx, periodId, voucherType)`

Configurations define patterns with auto-increment counters per period/type.

### Set of Books (SoB)

A SoB (账套) represents a complete accounting entity with its own:

- Chart of accounts
- Accounting periods
- Vouchers and ledgers
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

- Domain logic: voucher state transitions, business rule enforcement
- Report generator: formula calculations
- Repository updates with transactions

### Integration Tests

Bruno API tests in `bruno_collection/fims local/`: voucher creation → review → audit → post → verify ledgers.

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

- `internal/general_ledger/domain/voucher/` - Voucher aggregate with business rules
- `internal/general_ledger/domain/ledger/` - Ledger balance tracking
- `internal/sob/domain/sob/` - Set of Books configuration

**Critical Implementations**:

- `internal/general_ledger/app/command/post_voucher.go` - Posting logic
- `internal/report/domain/generator/generator.go` - Report generation
- `internal/common/errors/gin_middleware.go` - Error translation

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
// @Tags        vouchers
// @Summary     Post voucher to ledgers
// @Param       sobId path string true "Sob ID"
// @Router      /sob/{sobId}/voucher/{voucherId}/post [patch]
```

Main swagger config in `api/api.go`

## Code Style & Approach section

When asked to plan or implement changes, start with the simplest approach that follows existing patterns in the codebase. Do NOT over-engineer, create wrapper components, or introduce new abstractions unless explicitly requested. Look for existing patterns first and reuse them.

When asked for backend-only or frontend-only analysis, stay strictly within that boundary. Do not include suggestions or changes for the other side unless explicitly asked.

Before implementing a fix for a bug, create a brief plan and confirm the approach. Do not jump straight into coding a fix without understanding the root cause first. When debugging, avoid rapid-fire guessing — instead, methodically trace the data flow.
