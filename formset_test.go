package formset

import (
	"testing"
)

func TestBindProjectsLocaleDocumentsAndRoundTripsExtras(t *testing.T) {
	t.Parallel()
	record := RecordType{
		ID: "product", Label: "Product", Scope: ScopeWorkspace,
		Fields: []Field{
			{ID: "title", Label: "Title", Type: FieldText, Required: true, Localized: true},
			{ID: "price", Label: "Price", Type: FieldNumber, Required: true},
			{ID: "level", Label: "Level", Type: FieldSelect, Options: []Option{{Value: "advanced", Label: "Advanced"}}},
		},
	}
	form, err := BindDocuments(record, Documents{
		RU: map[string]any{"title": "Курс", "price": 39900.0, "level": "advanced", "kicker": "Backend"},
		EN: map[string]any{"title": "Course", "price": 39900.0, "level": "advanced"},
	})
	if err != nil {
		t.Fatalf("bind: %v", err)
	}
	if len(form.Issues) != 0 {
		t.Fatalf("unexpected issues: %#v", form.Issues)
	}
	if form.Values["ru"]["title"] != "Курс" || form.Extra["ru"]["kicker"] != "Backend" {
		t.Fatalf("values/extra mismatch: %#v %#v", form.Values, form.Extra)
	}
	payloads := form.PayloadDocuments()
	if payloads.RU["kicker"] != "Backend" || payloads.EN["title"] != "Course" {
		t.Fatalf("round-trip lost document keys: %#v", payloads)
	}
}

func TestBindReportsRequiredAndTypeIssues(t *testing.T) {
	t.Parallel()
	record := RecordType{
		ID: "product", Label: "Product", Scope: ScopeWorkspace,
		Fields: []Field{
			{ID: "title", Label: "Title", Type: FieldText, Required: true},
			{ID: "price", Label: "Price", Type: FieldNumber},
		},
	}
	form, err := Bind(record, map[string]map[string]any{
		"en": {"price": "free"},
	}, "en")
	if err != nil {
		t.Fatalf("bind: %v", err)
	}
	if len(form.Issues) != 2 {
		t.Fatalf("issues=%d want 2: %#v", len(form.Issues), form.Issues)
	}
}

func TestReviewAndDiff(t *testing.T) {
	t.Parallel()
	record := RecordType{
		ID: "lead", Label: "Lead", Scope: ScopeWorkspace,
		Fields: []Field{
			{ID: "title", Label: "Title", Type: FieldText},
			{ID: "title", Label: "Title Again", Type: FieldText},
		},
	}
	report := ReviewSchema([]RecordType{record}, nil)
	if !report.HasErrors() {
		t.Fatalf("expected duplicate field error, got %s", report.Summary())
	}
	diff := DiffForRecordTypes("k0", RecordType{ID: "lead", Label: "Lead", Scope: ScopeWorkspace})
	if err := diff.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestBindValidatesNestedObjectCollections(t *testing.T) {
	t.Parallel()
	record := RecordType{
		ID: "landing", Label: "Landing", Scope: ScopeWorkspace,
		Fields: []Field{
			{
				ID: "faq", Label: "FAQ", Type: FieldCollection,
				Items: &Field{
					ID: "item", Label: "Item", Type: FieldObject,
					Fields: []Field{
						{ID: "question", Label: "Question", Type: FieldText, Required: true},
						{ID: "answer", Label: "Answer", Type: FieldTextarea},
					},
				},
			},
			{
				ID: "author", Label: "Author", Type: FieldObject,
				Fields: []Field{
					{ID: "name", Label: "Name", Type: FieldText, Required: true},
					{
						ID: "links", Label: "Links", Type: FieldCollection,
						Items: &Field{ID: "item", Label: "Item", Type: FieldText},
					},
				},
			},
		},
	}
	form, err := BindDocuments(record, Documents{
		RU: map[string]any{
			"faq":    []any{map[string]any{"question": "Где?", "answer": "В IDE"}},
			"author": map[string]any{"name": "Анна", "links": []any{"https://example.test"}},
		},
		EN: map[string]any{
			"faq":    []any{map[string]any{"question": "Where?", "answer": "In the IDE"}},
			"author": map[string]any{"name": "Anna", "links": []any{"https://example.test"}},
		},
	})
	if err != nil {
		t.Fatalf("bind: %v", err)
	}
	if len(form.Issues) != 0 {
		t.Fatalf("issues: %#v", form.Issues)
	}
	broken, err := Bind(record, map[string]map[string]any{
		"en": {"faq": []any{map[string]any{"answer": "missing question"}}, "author": map[string]any{"name": 1}},
	}, "en")
	if err != nil {
		t.Fatalf("bind broken: %v", err)
	}
	if len(broken.Issues) == 0 {
		t.Fatal("expected nested type/required issues")
	}
}

func TestBindAcceptsRichTextAndMarkdownStrings(t *testing.T) {
	t.Parallel()
	record := RecordType{
		ID: "page", Label: "Page", Scope: ScopeWorkspace,
		Fields: []Field{
			{ID: "body", Label: "Body", Type: FieldRichText, Localized: true, UIHint: UIHintTipTap},
			{ID: "notes", Label: "Notes", Type: FieldMarkdown, UIHint: UIHintMarkdown},
		},
	}
	form, err := BindDocuments(record, Documents{
		RU: map[string]any{"body": "<p>Привет</p>", "notes": "# Заметки"},
		EN: map[string]any{"body": "<p>Hello</p>", "notes": "# Notes"},
	})
	if err != nil {
		t.Fatalf("bind: %v", err)
	}
	if len(form.Issues) != 0 {
		t.Fatalf("issues: %#v", form.Issues)
	}
	broken, err := Bind(record, map[string]map[string]any{
		"en": {"body": map[string]any{"type": "doc"}, "notes": 1},
	}, "en")
	if err != nil {
		t.Fatalf("bind broken: %v", err)
	}
	if len(broken.Issues) < 2 {
		t.Fatalf("expected type issues, got %#v", broken.Issues)
	}
}

func TestJSONSchemaMarksRichTextHint(t *testing.T) {
	t.Parallel()
	schema := JSONSchemaFromRecord(RecordType{
		ID: "page", Label: "Page", Scope: ScopeWorkspace,
		Fields: []Field{{ID: "body", Label: "Body", Type: FieldRichText, Required: true, UIHint: UIHintTipTap}},
	})
	properties, _ := schema["properties"].(map[string]any)
	body, _ := properties["body"].(map[string]any)
	if body["type"] != "string" || body["x-field-type"] != "richtext" || body["x-ui-hint"] != UIHintTipTap {
		t.Fatalf("unexpected schema: %#v", body)
	}
}

func TestJSONSchemaFromRecord(t *testing.T) {
	t.Parallel()
	schema := JSONSchemaFromRecord(RecordType{
		ID: "product", Label: "Product", Scope: ScopeWorkspace,
		Fields: []Field{{ID: "price", Label: "Price", Type: FieldNumber, Required: true}},
	})
	properties, _ := schema["properties"].(map[string]any)
	if properties["price"].(map[string]any)["type"] != "number" {
		t.Fatalf("unexpected schema: %#v", schema)
	}
	required, _ := schema["required"].([]string)
	if len(required) != 1 || required[0] != "price" {
		t.Fatalf("required=%v", required)
	}
}
