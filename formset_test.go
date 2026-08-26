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
