# GitHub Copilot Instructions for fims-proto-ms

## Project Overview

Financial Information Management System (FIMS) - A multi-tenant accounting system built with Go, using hexagonal architecture with CQRS patterns. Core domains: Set of Books (SoB/账套), General Ledger (总账), Reports, Numbering, and User management.

## Architecture

### Hexagonal Architecture (Ports & Adapters)

Each domain module follows this structure:

- **domain/** - Business logic, value objects, aggregates. Domain methods enforce business rules (e.g., `voucher.Post()`, `voucher.Audit()`)
- **app/** - Application services organized as Commands (writes) and Queries (reads) with handler pattern
- **port/** - Interfaces: `public/http` (external API), `private/http` (internal API), `private/intraprocess` (cross-module calls)
- **adapter/** - Infrastructure: `db` (repositories), cross-module adapters (e.g., `general_ledger/adapter/numbering`)

Example: `internal/general_ledger/`, `internal/sob/`, `internal/report/`

### Dependency Injection Pattern

Applications are initialized in `cmd/main.go` with manual DI:

1. Create repositories from datasource
2. Create empty Applications: `sobApp := sobApp.NewApplication()`
3. Create intraprocess interfaces for cross-module calls
4. Inject dependencies via `app.Inject(repo, readModel, externalServices...)`

Cross-module calls use intraprocess adapters, never direct repository access.

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
- `multitenant-datasource/` - Multi-tenant routing via subdomain (stub, not implemented)

Subdomain resolution: `ResolveSubdomain()` middleware extracts tenant from URL (`tenant.domain.com` → `tenant`)

## Domain-Specific Knowledge

### Voucher Lifecycle (General Ledger)

Workflow: Create → Review → Audit → Post (登账) → affects Ledgers

Business rules (see `internal/general_ledger/domain/voucher/`):

- Voucher must be reviewed AND audited before posting
- Creator ≠ Reviewer ≠ Auditor (segregation of duties)
- Cannot modify after audit/review
- Period must be current and open for posting
- Posting updates both `ledger` and `auxiliary_ledger` balances

### Posting Mechanics

File: `internal/general_ledger/app/command/post_voucher.go`

Process:

1. Call `v.Post(poster)` to validate and mark voucher
2. Build posting records for all line items + their superior accounts
3. Merge identical account records (sum debits/credits)
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

## Development Workflows

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
go test ./... -count=1

# Focus on:
# - Domain logic: voucher state transitions
# - Report generator: formula calculations
# - Repository updates with transactions
```

### Swagger Documentation

```bash
make swag          # Regenerate from @Tags/@Summary comments in port/public/http/
```

Visit: `http://127.0.0.1:4455/fims/swagger/index.html`

Swagger annotations in HTTP handlers:

```go
// @Tags        vouchers
// @Summary     Post voucher to ledgers
// @Param       sobId path string true "Sob ID"
// @Router      /sob/{sobId}/voucher/{voucherId}/post [patch]
```

### Code Formatting

```bash
make fmt           # swag fmt + gofumpt
make lint          # golangci-lint
```

## Key Files Reference

**Entry Point**: `cmd/main.go` - DI setup, router configuration, middleware registration

**Config**: `config/application-*.yaml`, `internal/common/config/default.go`

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

## PR Requirements

When making changes, ensure:

- [ ] Transactions use `repo.EnableTx()` pattern
- [ ] Domain state changes use VO methods, not direct mutation
- [ ] New business errors added to `i18n/zh-CN.json`
- [ ] Swagger annotations updated for API changes
- [ ] Tests cover business logic changes (especially voucher lifecycle, report formulas)
- [ ] No direct cross-module repository access (use intraprocess adapters)
- [ ] Schema changes documented if model fields modified

## Common Pitfalls

❌ Managing transactions in HTTP handlers
❌ Mutating domain objects without their methods
❌ Forgetting i18n entries for new slug errors
❌ Direct cross-module repository calls
❌ Modifying report formulas without tests
❌ Skipping `swag init` after API changes

## Testing Strategy (Bruno Collection)

API tests in `bruno_collection/fims local/`: voucher creation → review → audit → post → verify ledgers.

Use these flows as integration test templates.
