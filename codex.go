package formset

import "github.com/fastygo/codex/schema"

// FromCodex projects a Codex resource record into an editor form record.
// Codex owns the content fields. This package only binds them for a slot.
func FromCodex(record schema.RecordType) RecordType {
	return RecordType{
		ID:            RecordTypeID(record.ID),
		Label:         record.Label,
		Description:   record.Description,
		SchemaVersion: SchemaVersion(record.SchemaVersion),
		OwnerModule:   record.OwnerModule,
		Scope:         Scope(record.Scope),
		Fields:        projectFields(record.Fields),
		Relations:     projectRelations(record.Relations),
		Capabilities:  projectCapabilities(record.Capabilities),
		Visibility:    record.Visibility,
	}
}

func projectFields(fields []schema.Field) []Field {
	if fields == nil {
		return nil
	}
	projected := make([]Field, len(fields))
	for index, field := range fields {
		projected[index] = projectField(field)
	}
	return projected
}

func projectField(field schema.Field) Field {
	projected := Field{
		ID:           FieldID(field.ID),
		Label:        field.Label,
		Type:         FieldType(field.Type),
		Namespace:    field.Namespace,
		OwnerModule:  field.OwnerModule,
		Description:  field.Description,
		Required:     field.Required,
		Localized:    field.Localized,
		DefaultValue: field.DefaultValue,
		Options:      projectOptions(field.Options),
		Rules:        projectRules(field.Rules),
		Fields:       projectFields(field.Fields),
		Searchable:   field.Searchable,
		Indexed:      field.Indexed,
		Unique:       field.Unique,
		Sensitive:    field.Sensitive,
		Encrypted:    field.Encrypted,
		UIHint:       field.UIHint,
		StorageHint:  field.StorageHint,
	}
	if field.Items != nil {
		item := projectField(*field.Items)
		projected.Items = &item
	}
	return projected
}

func projectOptions(options []schema.Option) []Option {
	if options == nil {
		return nil
	}
	projected := make([]Option, len(options))
	for index, option := range options {
		projected[index] = Option{Value: option.Value, Label: option.Label}
	}
	return projected
}

func projectRules(rules []schema.ValidationRule) []ValidationRule {
	if rules == nil {
		return nil
	}
	projected := make([]ValidationRule, len(rules))
	for index, rule := range rules {
		projected[index] = ValidationRule{
			Name: rule.Name, Message: rule.Message, Severity: rule.Severity,
		}
		if rule.Args != nil {
			projected[index].Args = make(map[string]string, len(rule.Args))
			for key, value := range rule.Args {
				projected[index].Args[key] = value
			}
		}
	}
	return projected
}

func projectRelations(relations []schema.Relation) []Relation {
	if relations == nil {
		return nil
	}
	projected := make([]Relation, len(relations))
	for index, relation := range relations {
		projected[index] = Relation{
			ID:                   RelationID(relation.ID),
			Label:                relation.Label,
			Source:               RecordTypeID(relation.Source),
			Target:               RecordTypeID(relation.Target),
			Cardinality:          RelationCardinality(relation.Cardinality),
			InverseName:          relation.InverseName,
			CrossWorkspacePolicy: relation.CrossWorkspacePolicy,
			Policy: RelationPolicy{
				CrossWorkspaceMode: CrossWorkspaceMode(relation.Policy.CrossWorkspaceMode),
				AllowedTargets:     append([]string(nil), relation.Policy.AllowedTargets...),
				Capability:         CapabilityID(relation.Policy.Capability),
				ReadOnly:           relation.Policy.ReadOnly,
			},
			DeleteBehavior: DeleteBehavior(relation.DeleteBehavior),
		}
	}
	return projected
}

func projectCapabilities(capabilities []schema.CapabilityID) []CapabilityID {
	if capabilities == nil {
		return nil
	}
	projected := make([]CapabilityID, len(capabilities))
	for index, capability := range capabilities {
		projected[index] = CapabilityID(capability)
	}
	return projected
}
