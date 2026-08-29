package formset

import (
	"fmt"
	"strings"
)

func (field Field) Validate() error {
	if strings.TrimSpace(string(field.ID)) == "" {
		return fmt.Errorf("field id is required")
	}
	if strings.TrimSpace(field.Label) == "" {
		return fmt.Errorf("field %q label is required", field.ID)
	}
	if strings.TrimSpace(string(field.Type)) == "" {
		return fmt.Errorf("field %q type is required", field.ID)
	}
	if field.Type == FieldCollection && field.Items == nil {
		return fmt.Errorf("field %q collection requires items", field.ID)
	}
	if field.Type == FieldObject && len(field.Fields) == 0 {
		return fmt.Errorf("field %q object requires fields", field.ID)
	}
	if field.Type != FieldCollection && field.Items != nil {
		return fmt.Errorf("field %q items require collection type", field.ID)
	}
	if field.Type != FieldObject && len(field.Fields) > 0 {
		return fmt.Errorf("field %q nested fields require object type", field.ID)
	}
	if field.Items != nil {
		if err := field.Items.Validate(); err != nil {
			return err
		}
	}
	seen := map[FieldID]struct{}{}
	for _, nested := range field.Fields {
		if _, exists := seen[nested.ID]; exists {
			return fmt.Errorf("field %q nested field %q is duplicated", field.ID, nested.ID)
		}
		seen[nested.ID] = struct{}{}
		if err := nested.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (relation Relation) Validate() error {
	if strings.TrimSpace(string(relation.ID)) == "" {
		return fmt.Errorf("relation id is required")
	}
	if relation.Source == "" || relation.Target == "" {
		return fmt.Errorf("relation %q source and target are required", relation.ID)
	}
	if strings.TrimSpace(string(relation.Cardinality)) == "" {
		return fmt.Errorf("relation %q cardinality is required", relation.ID)
	}
	return nil
}

func (record RecordType) Validate() error {
	if strings.TrimSpace(string(record.ID)) == "" {
		return fmt.Errorf("record type id is required")
	}
	if strings.TrimSpace(record.Label) == "" {
		return fmt.Errorf("record type %q label is required", record.ID)
	}
	if strings.TrimSpace(string(record.Scope)) == "" {
		return fmt.Errorf("record type %q scope is required", record.ID)
	}
	for _, field := range record.Fields {
		if err := field.Validate(); err != nil {
			return err
		}
	}
	for _, relation := range record.Relations {
		if err := relation.Validate(); err != nil {
			return err
		}
	}
	return nil
}
