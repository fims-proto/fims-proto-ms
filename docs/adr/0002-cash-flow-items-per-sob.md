# ADR 0002: Cash Flow Items Are Per-SoB, Not Global

**Status**: Accepted

## Context

Cash Flow Items (OP_01, INV_01, FIN_01, etc.) are classification categories defined by an accounting standard (小企业会计准则). They look like global reference data — the same 16 items appear in every deployment. Journal lines hold a FK to a Cash Flow Item to classify the cash movement.

The question: should Cash Flow Items be a single global table (one row per item, shared by all SoBs) or a per-SoB table (rows seeded for each SoB at creation)?

## Decision

Cash Flow Items are **per-SoB**, seeded at SoB creation from the accounting standard's catalog (`dataload/xqykjzz/`).

## Rationale

Different accounting standards (小企业会计准则, 企业会计准则, etc.) define different Cash Flow Item sets. Multi-standard support is on the roadmap, and different SoBs may eventually use different standards. If items were global, adding multi-standard support would require retrofitting SoB-scoped filtering onto every query — a significant migration. By scoping items to SoB from the start, each SoB's chart of accounts, journal lines, and report templates all resolve Cash Flow Items within their own scope without cross-SoB joins.

This mirrors how accounts are handled: accounts are seeded per-SoB from the same `dataload/xqykjzz/accounts.csv` even though all SoBs using the same standard get the same account structure.

## Consequences

- A `sob_id` column is required on the `cash_flow_items` table.
- Journal line FKs reference the SoB-specific Cash Flow Item UUID (not a global code string).
- At SoB creation, Cash Flow Items must be seeded before journal lines can reference them.
- Report templates that reference Cash Flow Item codes resolve them to per-SoB UUIDs during template initialization or update.
- When multi-standard support is added, SoB gains an `accounting_standard` field and the seeding logic branches on it — no schema migration needed.
