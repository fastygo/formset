package formset_test

import (
	"fmt"

	"github.com/fastygo/codex/schema"
	"github.com/fastygo/formset"
)

func Example() {
	record := formset.FromCodex(schema.RecordType{
		ID: "offer", Label: "Offer", Scope: schema.ScopeWorkspace,
		Fields: []schema.Field{
			{ID: "title", Label: "Title", Type: schema.FieldText, Required: true, Localized: true},
			{ID: "price", Label: "Price", Type: schema.FieldNumber,
				Rules: []schema.ValidationRule{{Name: schema.RuleMoney}}},
		},
	})
	form, err := formset.BindLocale(record, "en", map[string]any{"title": "Lunch", "price": 12.5, "note": "kept"})
	if err != nil {
		panic(err)
	}
	for _, issue := range form.Issues {
		fmt.Println(issue.Field, issue.Code)
	}
	fmt.Println(form.Document("en")["note"])
	// Output:
	// price manifest.field_invalid
	// kept
}

func ExampleJSONSchemaFromRecord() {
	record := formset.RecordType{ID: "event", Label: "Event", Scope: formset.ScopeWorkspace, Fields: []formset.Field{
		{ID: "day", Label: "Day", Type: formset.FieldDateTime, Rules: []formset.ValidationRule{{Name: formset.RuleDate}}},
	}}
	day := formset.JSONSchemaFromRecord(record)["properties"].(map[string]any)["day"].(map[string]any)
	fmt.Println(day["type"], day["format"])
	// Output: string date
}
