# AGENTS.md

## Purpose

AI harness wrapper for fims-proto-ms.

## Shared Harness References

- ../.github/instructions/shared-core.instructions.md
- ../.github/instructions/architecture-ddd-go.instructions.md
- ../.github/instructions/testing-harness.instructions.md

## Repo-Specific Commands

- make test
- make fmt
- make lint
- make swag
- make peg

## Repo-Specific Rules

- Preserve DDD and hexagonal boundaries (domain/app/port/adapter).
- Keep command/query responsibilities separated.
- Use transaction helpers in app command handlers.
- Keep slug-based error and i18n patterns consistent.
