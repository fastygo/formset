# ADR-0001: Codex is the single source of truth for the record vocabulary

Status: accepted (2026-11-01)

## Context

formset kept a byte-for-byte copy of `codex/schema` (RecordType, Field,
Relation, ValidationRule, ...) and copied field by field in `FromCodex`.
The copy drifted (pinned to codex v0.3.0, `Severity` risk) and backend wrote
a third mapping by hand.

## Decision

- All record vocabulary types in formset are **type aliases** of
  `github.com/fastygo/codex/schema`. formset owns no struct for it.
- Scalar semantics (integer, decimal, money, date, uri, uuid, nullable,
  read-only, json-any) are Codex profile rules (`fastygo.codex/*`) on a base
  Codex field type. formset never adds first-class types for them.
- `Bind` validates the record with `schema.RecordType.Validate`,
  `schema.ValidateFieldProfile`, and value rules with
  `schema.CheckMetadataFields`. Codex is the authority; issue codes for rule
  violations are Codex codes (for example `manifest.field_invalid`).
- `FromCodex` stays the single documented entry point and is a deep copy.
- formset adds only form concerns: locale binding, `Extra` round-trip,
  JSON Schema for slots, schema review and diffs, editor UI hints.

## Consequences

- Breaking for v0.2.0 users that relied on distinct formset types; with
  aliases most code compiles unchanged.
- `CrossWorkspace*` constants are not exported by codex v0.4.2; formset keeps
  them typed as `schema.CrossWorkspaceMode` until Codex exports them.
