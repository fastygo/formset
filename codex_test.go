package formset

import (
	"testing"

	"github.com/fastygo/codex/schema"
)

func TestFromCodexProjectsRecordForBinding(t *testing.T) {
	t.Parallel()
	record := FromCodex(schema.RecordType{
		ID: "article", Label: "Articles", Scope: schema.ScopeTenant,
		Fields: []schema.Field{
			{ID: "title", Label: "Title", Type: schema.FieldText, Required: true, Localized: true},
			{
				ID: "media", Label: "Media", Type: schema.FieldRelation, UIHint: "media",
			},
		},
		Relations: []schema.Relation{{
			ID: "media", Source: "article", Target: "media",
			Cardinality: schema.RelationOneToOne, DeleteBehavior: schema.DeleteNullify,
		}},
	})
	if record.ID != "article" || record.Fields[1].UIHint != "media" {
		t.Fatalf("projection: %#v", record)
	}
	if record.Relations[0].Target != "media" || record.Relations[0].Cardinality != RelationOneToOne {
		t.Fatalf("relation: %#v", record.Relations)
	}
	form, err := BindLocale(record, "en", map[string]any{"title": "Hello"})
	if err != nil {
		t.Fatalf("bind: %v", err)
	}
	if len(form.Issues) != 0 || form.Document("en")["title"] != "Hello" {
		t.Fatalf("bound form: %#v", form)
	}
}
