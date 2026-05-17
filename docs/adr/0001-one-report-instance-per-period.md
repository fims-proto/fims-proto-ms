# ADR-0001: One Report Instance Per (SoB, Class, Period)

**Status**: Accepted

## Decision

Enforce exactly one report instance per (SoB, report class, accounting period) via a partial DB unique index on `(sob_id, class, period_id) WHERE template = false`. The generate operation becomes an idempotent upsert: if an instance already exists it is regenerated; otherwise it is created from the template.

## Context

The original model allowed unlimited instances per template per period. This created two problems:
1. No uniqueness guarantee — a user could generate duplicate balance sheets for the same period, making it unclear which was authoritative.
2. No connection to period closing — the statutory accounting workflow (close period → produce financial statements) was not enforced.

## Alternatives Considered

**Unlimited instances (original)**: Flexible for ad-hoc comparisons but lacks the uniqueness required for statutory reporting under 小企业会计准则. Does not support the period-close → snapshot flow.

**Application-only uniqueness (no DB constraint)**: Simpler migration, but leaves a gap for race conditions or direct DB access.

## Trade-offs

- Constraining to one instance per period removes the ability to create multiple "drafts" for the same period. This is acceptable — the template always holds the structural draft; the instance holds the canonical numbers for a given period.
- Report generation failure after period close is surfaced as a warning (`period-close-report-failed`) rather than rolled back — the close operation is an accounting event and must not be undone by a display concern.

## Consequences

- `POST /sob/{sobId}/report/{reportId}/generate` is now idempotent per period.
- Period closing (`ClosePeriodHandler`, `ClosePeriodsHandler`) triggers `GenerateForPeriod` for all report classes after the GL transaction commits.
- A partial unique index `(sob_id, class, period_id) WHERE template = false` is added via migration.
