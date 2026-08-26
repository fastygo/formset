# github.com/fastygo/formset

Schema kernel for typed admin forms. No renderer, BFF, or database.

```text
schema (RecordType)
  + values (payload_ru / payload_en)
  → Form
  → slot renders fields
  → save writes two JSON documents
```

## API

```go
form, err := formset.BindDocuments(record, formset.Documents{
    RU: payloadRU,
    EN: payloadEN,
})
payloads := form.PayloadDocuments()
schema := formset.JSONSchemaFromRecord(record)
```

`Form.Extra` keeps undeclared document keys so a save does not drop unknown fields.

## Scope

Included: record types, fields, relations, schema review/diff, form bind, JSON Schema.

Not included: Templ, Platform BFF, GraphQL, persistence. Those stay in the product and `fastygo/backend`.
