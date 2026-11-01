package formset

import (
	"maps"

	"github.com/fastygo/codex/schema"
)

// FromCodex is the single entry point from a Codex record into formset.
// The vocabulary is shared (type aliases), so the projection is a deep copy:
// mutating the result never mutates the caller's Codex manifest.
func FromCodex(record schema.RecordType) RecordType {
	projected := record
	projected.Fields = cloneFields(record.Fields)
	projected.Relations = cloneRelations(record.Relations)
	if record.Capabilities != nil {
		projected.Capabilities = append([]schema.CapabilityID(nil), record.Capabilities...)
	}
	return projected
}

// FromCodexResource projects a Codex manifest resource.
func FromCodexResource(resource schema.Resource) RecordType {
	return FromCodex(resource.Record)
}

// FromCodexManifest projects every resource of a Codex manifest, in order.
func FromCodexManifest(manifest schema.Manifest) []RecordType {
	records := make([]RecordType, 0, len(manifest.Resources))
	for _, resource := range manifest.Resources {
		records = append(records, FromCodexResource(resource))
	}
	return records
}

func cloneFields(fields []schema.Field) []schema.Field {
	if fields == nil {
		return nil
	}
	cloned := make([]schema.Field, len(fields))
	for index, field := range fields {
		cloned[index] = cloneField(field)
	}
	return cloned
}

func cloneField(field schema.Field) schema.Field {
	cloned := field
	if field.Options != nil {
		cloned.Options = append([]schema.Option(nil), field.Options...)
	}
	if field.Rules != nil {
		cloned.Rules = make([]schema.ValidationRule, len(field.Rules))
		for index, rule := range field.Rules {
			cloned.Rules[index] = rule
			cloned.Rules[index].Args = maps.Clone(rule.Args)
		}
	}
	cloned.Fields = cloneFields(field.Fields)
	if field.Items != nil {
		item := cloneField(*field.Items)
		cloned.Items = &item
	}
	return cloned
}

func cloneRelations(relations []schema.Relation) []schema.Relation {
	if relations == nil {
		return nil
	}
	cloned := make([]schema.Relation, len(relations))
	for index, relation := range relations {
		cloned[index] = relation
		if relation.Policy.AllowedTargets != nil {
			cloned[index].Policy.AllowedTargets = append([]string(nil), relation.Policy.AllowedTargets...)
		}
	}
	return cloned
}
