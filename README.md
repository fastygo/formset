# github.com/fastygo/formset

Schema kernel for typed admin forms. No renderer, BFF, or database.

**Codex is the single source of truth.** Every record type (`RecordType`,
`Field`, `Relation`, `ValidationRule`, ...) is a type alias of
[`github.com/fastygo/codex/schema`](https://github.com/fastygo/codex)
(compatible with `@ui8kit/contract` and schema.buildy.tw).
formset only adds form concerns on top.

| formset | requires codex |
|---|---|
| main (Unreleased) | v0.4.2 |
| v0.2.0 | v0.3.0 |

```text
codex schema.RecordType
  → formset.FromCodex (deep copy, no remapping)
  + locale + data document
  → formset.Bind → Form (values, extra, issues)
  → slot renders fields → save writes that locale document
```

## Install

```sh
go get github.com/fastygo/formset@main github.com/fastygo/codex@v0.4.2
```

## Usage

```go
record := formset.FromCodex(schema.RecordType{
	ID: "offer", Label: "Offer", Scope: schema.ScopeWorkspace,
	Fields: []schema.Field{
		{ID: "title", Label: "Title", Type: schema.FieldText, Required: true, Localized: true},
		// money is a Codex rule on a number field: integer minor units.
		{ID: "price", Label: "Price", Type: schema.FieldNumber,
			Rules: []schema.ValidationRule{{Name: schema.RuleMoney}}},
		{ID: "currency", Label: "Currency", Type: schema.FieldString},
	},
})

form, err := formset.BindLocale(record, "en", map[string]any{"title": "Lunch", "price": 1250.0})
// err: the record or a field profile is invalid (Codex rules).
// form.Issues: required / type / option, plus Codex value-rule codes.
document := form.Document("en")            // values + unknown keys (Extra)
jsonSchema := formset.JSONSchemaFromRecord(record)
report := formset.ReviewSchema(formset.FromCodexManifest(manifest), nil)
```

`Bind` accepts an explicit locale map when more than one document is in
memory. Codex reads/writes one locale per request (`?locale=` + fallback).
`Form.Extra` keeps undeclared keys so a save does not drop unknown fields.

## Types and rules

Field types are Codex field types (`text`, `string`, `number`, `boolean`,
`select`, `datetime`, `json`, `object`, `collection`, `relation`, ...).
Scalar semantics are Codex rules, re-exported for convenience:

| need | Codex declaration | JSON Schema |
|---|---|---|
| integer | `number` + `RuleInteger` | `integer` |
| decimal | `number` + `RuleDecimal` | `number` |
| money | `number` + `RuleMoney` (minor units, `currency_field` arg) | `integer`, `x-codex-money` |
| date | `datetime` + `RuleDate` | `format: date` |
| uri / uuid | `string` + `RuleURI` / `RuleUUID` | `format: uri` / `uuid` |
| enum | `select` + `Options` | `enum` |
| media | `relation` + `UIHint: UIHintMedia` | `x-ui-hint: media` |
| nullable / read-only | `RuleNullable` / `RuleReadOnly` | type union / `readOnly` |

## Scope

Included: form binding, Codex-rule validation of bound values, schema review
and diffs, JSON Schema for editor slots.

Not included: content ownership, Templ, a BFF, GraphQL, or persistence.

## Development

```sh
gofmt -l . && go mod tidy -diff && go vet ./... && go test -race ./... && golangci-lint run
```

See [CHANGELOG.md](CHANGELOG.md) and [docs/adr](docs/adr).
