# github.com/fastygo/formset

Schema kernel for typed admin forms. No renderer, BFF, or database.

```text
schema (RecordType)
  + locale + data document
  → Form
  → slot renders fields
  → save writes that locale document
```

## API

```go
form, err := formset.BindLocale(record, "en", data)
document := form.Document("en")
schema := formset.JSONSchemaFromRecord(record)
```

`Bind` accepts an explicit locale map when more than one document is in
memory. Codex reads/writes one locale per request (`?locale=` + fallback).
`Form.Extra` keeps undeclared keys so a save does not drop unknown fields.

## Scope

Included: record types, fields, relations, schema review/diff, form bind, JSON Schema.

Not included: Templ, Platform BFF, GraphQL, persistence. Those stay in the product and `fastygo/backend`.
