package formset

// JSONSchema is a minimal JSON Schema object the slot or backend can publish.
type JSONSchema map[string]any

// JSONSchemaFromRecord describes one locale document for the record type.
func JSONSchemaFromRecord(record RecordType) JSONSchema {
	properties := map[string]any{}
	required := make([]string, 0)
	for _, field := range record.Fields {
		properties[string(field.ID)] = fieldSchema(field)
		if field.Required {
			required = append(required, string(field.ID))
		}
	}
	schema := JSONSchema{
		"$schema":    "https://json-schema.org/draft/2020-12/schema",
		"type":       "object",
		"title":      record.Label,
		"properties": properties,
	}
	if record.Description != "" {
		schema["description"] = record.Description
	}
	if len(required) > 0 {
		schema["required"] = required
	}
	return schema
}

func fieldSchema(field Field) map[string]any {
	schema := map[string]any{"title": field.Label, "x-field-type": string(field.Type)}
	if field.Description != "" {
		schema["description"] = field.Description
	}
	if field.UIHint != "" {
		schema["x-ui-hint"] = field.UIHint
	}
	switch field.Type {
	case FieldNumber:
		schema["type"] = "number"
	case FieldBoolean:
		schema["type"] = "boolean"
	case FieldJSON:
		schema["type"] = []string{"object", "array"}
	case FieldObject:
		schema["type"] = "object"
		properties := map[string]any{}
		for _, nested := range field.Fields {
			properties[string(nested.ID)] = fieldSchema(nested)
		}
		schema["properties"] = properties
	case FieldCollection:
		schema["type"] = "array"
	case FieldSelect:
		schema["type"] = "string"
		if len(field.Options) > 0 {
			enum := make([]string, 0, len(field.Options))
			for _, option := range field.Options {
				enum = append(enum, option.Value)
			}
			schema["enum"] = enum
		}
	default:
		schema["type"] = "string"
	}
	if field.Items != nil {
		schema["type"] = "array"
		schema["items"] = fieldSchema(*field.Items)
	}
	applyCodexRules(field, schema)
	return schema
}

// applyCodexRules maps fastygo.codex/* profile rules onto JSON Schema
// keywords. Rule names are also published verbatim in x-codex-rules.
func applyCodexRules(field Field, schema map[string]any) {
	if field.Type == FieldDateTime {
		schema["format"] = "date-time"
	}
	if len(field.Rules) == 0 {
		return
	}
	names := make([]string, 0, len(field.Rules))
	nullable := false
	for _, rule := range field.Rules {
		names = append(names, rule.Name)
		switch rule.Name {
		case RuleInteger:
			schema["type"] = "integer"
		case RuleMoney:
			schema["type"] = "integer"
			schema["x-codex-money"] = map[string]any{"unit": "minor", "currency_field": moneyCurrencyField(rule)}
		case RuleDate:
			schema["format"] = "date"
		case RuleURI:
			schema["format"] = "uri"
		case RuleUUID:
			schema["format"] = "uuid"
		case RuleJSONAny:
			delete(schema, "type")
		case RuleReadOnly:
			schema["readOnly"] = true
		case RuleNullable:
			nullable = true
		}
	}
	if typed, ok := schema["type"].(string); ok && nullable {
		schema["type"] = []string{typed, "null"}
	}
	schema["x-codex-rules"] = names
}

func moneyCurrencyField(rule ValidationRule) string {
	if rule.Args != nil && rule.Args["currency_field"] != "" {
		return rule.Args["currency_field"]
	}
	return "currency"
}
