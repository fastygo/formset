# Changelog

All notable changes to `github.com/fastygo/formset`. Versions follow semver;
v0.x minors may break. No tag has been cut for Unreleased.

## Unreleased

### Changed (breaking)
- Depends on `github.com/fastygo/codex` **v0.4.2** (was v0.3.0).
- Record vocabulary types are now type aliases of `codex/schema`; the
  formset-owned copy (`types.go` structs, `validate.go`) is removed
  ([ADR-0001](docs/adr/0001-codex-single-source-of-truth.md)).
- `Bind` rejects fields whose Codex profile is invalid (unknown rule, rule on
  the wrong base type, required + nullable).
- `Bind` reports Codex value-rule violations (integer, money in minor units,
  decimal, date, uri, uuid, json-any) with Codex issue codes.

### Added
- `FromCodexResource`, `FromCodexManifest`.
- Re-exported Codex rule names `RuleInteger`, `RuleMoney`, ... and `UIHintMedia`.
- JSON Schema: `integer` for integer/money, `format` date/date-time/uri/uuid,
  `readOnly`, nullable type unions, `x-codex-rules`, `x-codex-money`.
- CI: gofmt, `go mod tidy -diff`, build, vet, `test -race`, golangci-lint.
- Runnable examples (`example_test.go`). Coverage 63.8% → 96.5%.

### Known Codex gaps
- codex v0.4.2 has no `CrossWorkspaceMode` constants; formset keeps
  `CrossWorkspaceForbidden` etc. typed as the Codex type.

## v0.2.0, v0.1.0
- Initial form kernel with its own copy of the Codex record vocabulary.
