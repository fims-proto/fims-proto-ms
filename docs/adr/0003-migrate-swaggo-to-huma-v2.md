# ADR-0003: Migrate from swaggo/swag (OAS2) to danielgtaylor/huma/v2 (OAS3)

## Status

Accepted

## Context

The project used `github.com/swaggo/swag` to generate Swagger 2.0 docs via comment annotations. This blocked adoption of OpenAPI 3.x features (Links, `oneOf`/`anyOf`, richer examples, response headers), which are needed for the next phase of API client development.

## Decision

Replace `swaggo/swag` with `github.com/danielgtaylor/huma/v2` using the `humagin` adapter to preserve the existing Gin router and middleware stack.

Key choices made:

- **Keep Gin**: `humaginv2`/`humagin` wraps the existing `*gin.Engine`. All existing Gin middleware (error handler, subdomain resolver) stays in place for private routes.
- **Public API only**: Only `/api/v1` routes are registered through huma. Private `/internal` routes remain as plain Gin handlers — they exist solely for DB migrations and will be removed in a future approach.
- **Error handling via `huma.NewErrorWithContext` override**: `SlugErr` is not a `huma.StatusError`, so huma's default path calls `NewErrorWithContext`. We override this var in `internal/common/errors/huma_error.go` to detect `SlugErr`, look up i18n using `ctx.Header("Accept-Language")`, and return a `humaSlugError` with the project's `{slug, message}` JSON format.
- **Pagination params renamed**: `$page`/`$size`/`$sort`/`$filter` → `page`/`size`/`sort`/`filter`. No existing Bruno tests used these params so the rename is safe.
- **`--gen-spec` flag**: `go run cmd/main.go --gen-spec` initialises the full DI stack, registers all routes, marshals `api.OpenAPI()` to JSON, writes `docs/swagger_generated/openapi.json`, and exits. Used by `make swag`.

## Consequences

- All public handler functions now have typed `(ctx context.Context, input *XInput) (*XOutput, error)` signatures instead of `func(*gin.Context)`.
- `swaggo/swag`, `swaggo/gin-swagger`, and `swaggo/files` are removed from `go.mod`.
- `api/api.go` (global swag annotations) and `docs/swagger_generated/` (generated output) are deleted.
- The spec is served at `/openapi.json` and `/openapi.yaml` in all environments (previously only in `dev`).
- `PagingResponseProcessor` and `NewPageRequestFromQuery` in `gin_helper.go` are no longer used by public handlers but remain for any future internal use.
