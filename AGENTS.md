# Agent notes

FormSet binds editor forms. It does not own content records, storage, or HTTP.

Codex (`github.com/fastygo/codex`) owns fields and relations. Project them with `FromCodex`. Do not copy a second content vocabulary into this module, and do not use a local replace for Codex.

## Rules

1. Do not add a renderer, database, BFF, GraphQL, or product manifest.
2. A form record is a projection of a Codex record plus locale documents.
3. Unknown document keys stay in `Form.Extra` so a save can round-trip them.
4. Comments and documentation are written in English.
5. Pin Codex by a published tag.

Run before completion:

```text
go test ./...
go vet ./...
gofmt -w .
```
