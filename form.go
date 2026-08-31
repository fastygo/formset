package formset

import (
	"fmt"
	"slices"
)

type Issue struct {
	Locale  string `json:"locale,omitempty"`
	Field   string `json:"field,omitempty"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Form is the slot contract: schema fields plus one document per locale.
type Form struct {
	Record  RecordTypeID              `json:"record"`
	Locales []string                  `json:"locales"`
	Fields  []Field                   `json:"fields"`
	Values  map[string]map[string]any `json:"values"`
	Extra   map[string]map[string]any `json:"extra,omitempty"`
	Issues  []Issue                   `json:"issues,omitempty"`
}

// BindLocale projects one locale document. Codex and the admin bind the
// requested locale only; fallback is a storage/read concern, not FormSet.
func BindLocale(record RecordType, locale string, document map[string]any) (Form, error) {
	if locale == "" {
		return Form{}, fmt.Errorf("locale is required")
	}
	return Bind(record, map[string]map[string]any{locale: document}, locale)
}

// Bind projects a record type and locale documents into a renderable form.
// Unknown document keys are kept in Extra so a save can round-trip them.
// Locales come from the caller or from document keys. There is no default ru/en pair.
func Bind(record RecordType, documents map[string]map[string]any, locales ...string) (Form, error) {
	if err := record.Validate(); err != nil {
		return Form{}, err
	}
	if documents == nil {
		documents = map[string]map[string]any{}
	}
	resolvedLocales := locales
	if len(resolvedLocales) == 0 {
		for locale := range documents {
			resolvedLocales = append(resolvedLocales, locale)
		}
		slices.Sort(resolvedLocales)
	}
	declared := make(map[FieldID]Field, len(record.Fields))
	for _, field := range record.Fields {
		declared[field.ID] = field
	}
	form := Form{
		Record:  record.ID,
		Locales: append([]string(nil), resolvedLocales...),
		Fields:  append([]Field(nil), record.Fields...),
		Values:  map[string]map[string]any{},
		Extra:   map[string]map[string]any{},
	}
	for _, locale := range resolvedLocales {
		source := documents[locale]
		values := make(map[string]any, len(record.Fields))
		extra := map[string]any{}
		for _, field := range record.Fields {
			if source == nil {
				if field.Required {
					form.Issues = append(form.Issues, Issue{
						Locale: locale, Field: string(field.ID), Code: "required",
						Message: fmt.Sprintf("%s is required", field.Label),
					})
				}
				continue
			}
			value, exists := source[string(field.ID)]
			if !exists || value == nil || value == "" {
				if field.Required {
					form.Issues = append(form.Issues, Issue{
						Locale: locale, Field: string(field.ID), Code: "required",
						Message: fmt.Sprintf("%s is required", field.Label),
					})
				}
				continue
			}
			if issue := typeIssue(locale, field, value); issue != nil {
				form.Issues = append(form.Issues, *issue)
			}
			values[string(field.ID)] = value
		}
		for key, value := range source {
			if _, known := declared[FieldID(key)]; !known {
				extra[key] = value
			}
		}
		form.Values[locale] = values
		if len(extra) > 0 {
			form.Extra[locale] = extra
		}
	}
	return form, nil
}

// Documents reconstructs locale objects for storage, including Extra keys.
func (form Form) Documents() map[string]map[string]any {
	result := make(map[string]map[string]any, len(form.Locales))
	for _, locale := range form.Locales {
		result[locale] = form.Document(locale)
	}
	return result
}

// Document reconstructs one locale object, including Extra keys.
func (form Form) Document(locale string) map[string]any {
	document := map[string]any{}
	for key, value := range form.Values[locale] {
		document[key] = value
	}
	for key, value := range form.Extra[locale] {
		if _, exists := document[key]; !exists {
			document[key] = value
		}
	}
	return document
}

func typeIssue(locale string, field Field, value any) *Issue {
	switch field.Type {
	case FieldNumber:
		switch value.(type) {
		case float64, float32, int, int32, int64:
			return nil
		default:
			return &Issue{Locale: locale, Field: string(field.ID), Code: "type", Message: "value must be a number"}
		}
	case FieldBoolean:
		if _, ok := value.(bool); !ok {
			return &Issue{Locale: locale, Field: string(field.ID), Code: "type", Message: "value must be a boolean"}
		}
	case FieldSelect:
		text, ok := value.(string)
		if !ok {
			return &Issue{Locale: locale, Field: string(field.ID), Code: "type", Message: "value must be a string"}
		}
		if len(field.Options) == 0 {
			return nil
		}
		for _, option := range field.Options {
			if option.Value == text {
				return nil
			}
		}
		return &Issue{Locale: locale, Field: string(field.ID), Code: "option", Message: "value is not a declared option"}
	case FieldText, FieldString, FieldTextarea, FieldRichText, FieldMarkdown, FieldDateTime:
		if _, ok := value.(string); !ok {
			return &Issue{Locale: locale, Field: string(field.ID), Code: "type", Message: "value must be a string"}
		}
	case FieldObject:
		document, ok := value.(map[string]any)
		if !ok {
			return &Issue{Locale: locale, Field: string(field.ID), Code: "type", Message: "value must be an object"}
		}
		for _, nested := range field.Fields {
			item, exists := document[string(nested.ID)]
			if !exists || item == nil || item == "" {
				if nested.Required {
					return &Issue{
						Locale: locale, Field: string(nested.ID), Code: "required",
						Message: fmt.Sprintf("%s is required", nested.Label),
					}
				}
				continue
			}
			if issue := typeIssue(locale, nested, item); issue != nil {
				return issue
			}
		}
	case FieldCollection:
		if field.Items == nil {
			return nil
		}
		switch items := value.(type) {
		case []any:
			for _, item := range items {
				if issue := typeIssue(locale, *field.Items, item); issue != nil {
					return issue
				}
			}
		case []string:
			for _, item := range items {
				if issue := typeIssue(locale, *field.Items, item); issue != nil {
					return issue
				}
			}
		default:
			return &Issue{Locale: locale, Field: string(field.ID), Code: "type", Message: "value must be an array"}
		}
	}
	return nil
}
