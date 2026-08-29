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
	schema := map[string]any{"title": field.Label}
	if field.Description != "" {
		schema["description"] = field.Description
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
	return schema
}
