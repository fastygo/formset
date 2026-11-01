package formset

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/fastygo/codex/schema"
)

func codexFixture() schema.RecordType {
	return schema.RecordType{
		ID: "offer", Label: "Offer", Scope: schema.ScopeWorkspace, SchemaVersion: "1",
		Capabilities: []schema.CapabilityID{"offer.edit"},
		Fields: []schema.Field{
			{ID: "title", Label: "Title", Type: schema.FieldText, Required: true, Localized: true, Searchable: true},
			{ID: "price", Label: "Price", Type: schema.FieldNumber, Rules: []schema.ValidationRule{
				{Name: schema.RuleMoney, Message: "minor units", Severity: "error", Args: map[string]string{"currency_field": "currency"}},
			}},
			{ID: "currency", Label: "Currency", Type: schema.FieldString},
			{ID: "level", Label: "Level", Type: schema.FieldSelect, Options: []schema.Option{{Value: "a", Label: "A"}}},
			{ID: "tags", Label: "Tags", Type: schema.FieldCollection, Items: &schema.Field{ID: "tag", Label: "Tag", Type: schema.FieldString}},
			{ID: "seo", Label: "SEO", Type: schema.FieldObject, Fields: []schema.Field{{ID: "title", Label: "SEO title", Type: schema.FieldString}}},
		},
		Relations: []schema.Relation{{
			ID: "owner", Source: "offer", Target: "user", Cardinality: schema.RelationOneToMany,
			Policy: schema.RelationPolicy{CrossWorkspaceMode: CrossWorkspaceForbidden, AllowedTargets: []string{"user"}},
		}},
	}
}

func TestFromCodexRoundTripIsLossless(t *testing.T) {
	t.Parallel()
	source := codexFixture()
	projected := FromCodex(source)
	if !reflect.DeepEqual(source, schema.RecordType(projected)) {
		t.Fatalf("projection is not lossless:\n%#v\n%#v", source, projected)
	}
	if projected.Fields[1].Rules[0].Severity != "error" {
		t.Fatalf("severity lost: %#v", projected.Fields[1].Rules)
	}
	want, _ := json.Marshal(source)
	got, _ := json.Marshal(projected)
	if string(want) != string(got) {
		t.Fatalf("json differs:\n%s\n%s", want, got)
	}
}

func TestFromCodexDeepCopies(t *testing.T) {
	t.Parallel()
	source := codexFixture()
	projected := FromCodex(source)
	projected.Fields[1].Rules[0].Args["currency_field"] = "x"
	projected.Fields[3].Options[0].Value = "x"
	projected.Fields[4].Items.ID = "x"
	projected.Fields[5].Fields[0].ID = "x"
	projected.Relations[0].Policy.AllowedTargets[0] = "x"
	projected.Capabilities[0] = "x"
	if !reflect.DeepEqual(source, codexFixture()) {
		t.Fatalf("source mutated: %#v", source)
	}
	if FromCodex(schema.RecordType{ID: "x"}).Fields != nil {
		t.Fatal("nil fields must stay nil")
	}
}

func TestFromCodexManifest(t *testing.T) {
	t.Parallel()
	manifest := schema.Manifest{Name: "m", Version: "1", Resources: []schema.Resource{{Record: codexFixture()}, {Record: schema.RecordType{ID: "user", Label: "User", Scope: schema.ScopeWorkspace}}}}
	records := FromCodexManifest(manifest)
	if len(records) != 2 || records[0].ID != "offer" || records[1].ID != "user" {
		t.Fatalf("records: %#v", records)
	}
	if FromCodexResource(manifest.Resources[1]).ID != "user" {
		t.Fatal("resource projection")
	}
	if report := ReviewSchema(records, nil); report.HasErrors() {
		t.Fatalf("review: %s", report.Summary())
	}
}

func TestBindEnforcesCodexRules(t *testing.T) {
	t.Parallel()
	record := FromCodex(schema.RecordType{
		ID: "r", Label: "R", Scope: schema.ScopeWorkspace,
		Fields: []schema.Field{
			{ID: "count", Label: "Count", Type: schema.FieldNumber, Rules: []schema.ValidationRule{{Name: schema.RuleInteger}}},
			{ID: "price", Label: "Price", Type: schema.FieldNumber, Rules: []schema.ValidationRule{{Name: schema.RuleMoney}}},
			{ID: "site", Label: "Site", Type: schema.FieldString, Rules: []schema.ValidationRule{{Name: schema.RuleURI}}},
			{ID: "uid", Label: "UID", Type: schema.FieldString, Rules: []schema.ValidationRule{{Name: schema.RuleUUID}}},
			{ID: "day", Label: "Day", Type: schema.FieldDateTime, Rules: []schema.ValidationRule{{Name: schema.RuleDate}}},
		},
	})
	good, err := BindLocale(record, "en", map[string]any{
		"count": 3.0, "price": 1250.0, "site": "https://x.test", "uid": "123e4567-e89b-12d3-a456-426614174000", "day": "2026-11-01",
	})
	if err != nil || len(good.Issues) != 0 {
		t.Fatalf("good: %v %#v", err, good.Issues)
	}
	bad, err := BindLocale(record, "en", map[string]any{
		"count": 1.5, "price": 12.5, "site": "nope", "uid": "x", "day": "01.11.2026",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(bad.Issues) != 5 {
		t.Fatalf("want 5 codex issues, got %#v", bad.Issues)
	}
	for _, issue := range bad.Issues {
		if issue.Locale != "en" || issue.Code == "" || issue.Message == "" {
			t.Fatalf("issue: %#v", issue)
		}
	}
}

func TestBindRejectsUnknownCodexRule(t *testing.T) {
	t.Parallel()
	record := RecordType{ID: "r", Label: "R", Scope: ScopeWorkspace, Fields: []Field{
		{ID: "a", Label: "A", Type: FieldString, Rules: []ValidationRule{{Name: "acme/unknown"}}},
	}}
	if _, err := BindLocale(record, "en", nil); err == nil {
		t.Fatal("unknown rule must be rejected")
	}
	record.Fields[0].Rules = []ValidationRule{{Name: RuleMoney}}
	if _, err := BindLocale(record, "en", nil); err == nil {
		t.Fatal("money on string must be rejected")
	}
}

func TestJSONSchemaPublishesCodexRules(t *testing.T) {
	t.Parallel()
	record := RecordType{ID: "r", Label: "R", Scope: ScopeWorkspace, Fields: []Field{
		{ID: "n", Label: "N", Type: FieldNumber, Rules: []ValidationRule{{Name: RuleInteger}, {Name: RuleNullable}}},
		{ID: "p", Label: "P", Type: FieldNumber, Rules: []ValidationRule{{Name: RuleMoney, Args: map[string]string{"currency_field": "cur"}}, {Name: RuleReadOnly}}},
		{ID: "d", Label: "D", Type: FieldDateTime, Rules: []ValidationRule{{Name: RuleDate}}},
		{ID: "t", Label: "T", Type: FieldDateTime},
		{ID: "u", Label: "U", Type: FieldString, Rules: []ValidationRule{{Name: RuleURI}}},
		{ID: "id", Label: "ID", Type: FieldString, Rules: []ValidationRule{{Name: RuleUUID}}},
		{ID: "j", Label: "J", Type: FieldJSON, Rules: []ValidationRule{{Name: RuleJSONAny}}},
		{ID: "b", Label: "B", Type: FieldBoolean},
		{ID: "o", Label: "O", Type: FieldObject, Fields: []Field{{ID: "x", Label: "X", Type: FieldString}}},
		{ID: "c", Label: "C", Type: FieldCollection, Items: &Field{ID: "i", Label: "I", Type: FieldString}},
		{ID: "s", Label: "S", Type: FieldSelect, Options: []Option{{Value: "a", Label: "A"}}},
	}}
	props := JSONSchemaFromRecord(record)["properties"].(map[string]any)
	get := func(id string) map[string]any { return props[id].(map[string]any) }
	if !reflect.DeepEqual(get("n")["type"], []string{"integer", "null"}) {
		t.Fatalf("n: %#v", get("n"))
	}
	if get("p")["type"] != "integer" || get("p")["readOnly"] != true || get("p")["x-codex-money"].(map[string]any)["currency_field"] != "cur" {
		t.Fatalf("p: %#v", get("p"))
	}
	if get("d")["format"] != "date" || get("t")["format"] != "date-time" || get("u")["format"] != "uri" || get("id")["format"] != "uuid" {
		t.Fatal("formats")
	}
	if _, typed := get("j")["type"]; typed {
		t.Fatal("json-any must be untyped")
	}
	if get("b")["type"] != "boolean" || get("o")["type"] != "object" || get("c")["type"] != "array" || get("s")["enum"] == nil {
		t.Fatal("base types")
	}
}

func TestReviewAndDiffs(t *testing.T) {
	t.Parallel()
	records := []RecordType{
		{ID: "a", Label: "A", Scope: ScopeWorkspace, Fields: []Field{
			{ID: "x", Label: "X", Type: FieldString}, {ID: "x", Label: "X", Type: FieldText},
			{ID: "secret", Label: "S", Type: FieldEncrypted, Encrypted: true},
		}},
		{ID: "a", Label: "A2", Scope: ScopeWorkspace},
	}
	relations := []Relation{{ID: "r", Source: "missing", Target: "gone", Cardinality: RelationOneToOne, CrossWorkspacePolicy: "legacy"}}
	diff := SchemaDiff{ID: "d", Entries: []DiffEntry{{Operation: DiffRemoveField, RecordID: "a", FieldID: "x", Breaking: true}}}
	report := ReviewSchema(records, relations, diff)
	codes := map[string]bool{}
	for _, finding := range report.Findings {
		codes[finding.Code] = true
	}
	for _, code := range []string{"duplicate-record-type", "duplicate-field", "encrypted-not-sensitive", "missing-relation-source", "missing-relation-target", "legacy-cross-workspace-policy", "breaking-schema-change"} {
		if !codes[code] {
			t.Fatalf("missing %s in %s", code, report.Summary())
		}
	}
	if !report.HasErrors() || (ReviewReport{}).Summary() != "schema review passed" {
		t.Fatal("report helpers")
	}
	if (SchemaDiff{}).Validate() == nil || (SchemaDiff{ID: "d", Entries: []DiffEntry{{}}}).Validate() == nil || (SchemaDiff{ID: "d", Entries: []DiffEntry{{Operation: DiffAddField}}}).Validate() == nil {
		t.Fatal("diff validation")
	}
	if err := DiffForRecordTypes("x", records[0]).Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestBindTypeIssues(t *testing.T) {
	t.Parallel()
	record := RecordType{ID: "r", Label: "R", Description: "d", Scope: ScopeWorkspace, Fields: []Field{
		{ID: "n", Label: "N", Type: FieldNumber},
		{ID: "b", Label: "B", Type: FieldBoolean},
		{ID: "s", Label: "S", Type: FieldString},
		{ID: "sel", Label: "Sel", Type: FieldSelect, Options: []Option{{Value: "a", Label: "A"}}},
		{ID: "sel2", Label: "Sel2", Type: FieldSelect},
		{ID: "o", Label: "O", Type: FieldObject, Fields: []Field{{ID: "req", Label: "Req", Type: FieldString, Required: true}}},
		{ID: "c", Label: "C", Type: FieldCollection, Items: &Field{ID: "i", Label: "I", Type: FieldString}},
		{ID: "need", Label: "Need", Type: FieldString, Required: true},
	}}
	form, err := Bind(record, map[string]map[string]any{"en": {
		"n": "1", "b": "yes", "s": 1, "sel": "z", "sel2": 2, "o": map[string]any{}, "c": "x",
	}, "ru": nil}, "en", "ru")
	if err != nil {
		t.Fatal(err)
	}
	if len(form.Issues) != 9 {
		t.Fatalf("issues: %#v", form.Issues)
	}
	if JSONSchemaFromRecord(record)["description"] != "d" {
		t.Fatal("description")
	}
	if _, err := Bind(RecordType{}, nil); err == nil {
		t.Fatal("invalid record")
	}
	report := ReviewSchema([]RecordType{record, record}, nil)
	if report.Summary() == "" {
		t.Fatal("summary")
	}
	if moneyCurrencyField(ValidationRule{}) != "currency" {
		t.Fatal("default currency field")
	}
}
