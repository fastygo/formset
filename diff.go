package formset

import "fmt"

type DiffOperation string

const (
	DiffAddRecordType    DiffOperation = "add-record-type"
	DiffExtendRecordType DiffOperation = "extend-record-type"
	DiffAddField         DiffOperation = "add-field"
	DiffChangeField      DiffOperation = "change-field"
	DiffRemoveField      DiffOperation = "remove-field"
	DiffDeprecateField   DiffOperation = "deprecate-field"
	DiffAddRelation      DiffOperation = "add-relation"
	DiffChangeRelation   DiffOperation = "change-relation"
)

type DiffEntry struct {
	Operation   DiffOperation `json:"operation"`
	RecordID    RecordTypeID  `json:"record_id"`
	FieldID     FieldID       `json:"field_id,omitempty"`
	Relation    RelationID    `json:"relation,omitempty"`
	Summary     string        `json:"summary,omitempty"`
	OwnerModule string        `json:"owner_module,omitempty"`
	Breaking    bool          `json:"breaking,omitempty"`
}

type SchemaDiff struct {
	ID      string      `json:"id"`
	Entries []DiffEntry `json:"entries"`
}

func (diff SchemaDiff) Validate() error {
	if diff.ID == "" {
		return fmt.Errorf("schema diff id is required")
	}
	for _, entry := range diff.Entries {
		if entry.Operation == "" {
			return fmt.Errorf("schema diff %q contains entry without operation", diff.ID)
		}
		if entry.RecordID == "" {
			return fmt.Errorf("schema diff %q contains entry without record id", diff.ID)
		}
	}
	return nil
}

func DiffForRecordTypes(id string, records ...RecordType) SchemaDiff {
	entries := make([]DiffEntry, 0, len(records))
	for _, record := range records {
		entries = append(entries, DiffEntry{
			Operation: DiffAddRecordType,
			RecordID:  record.ID,
			Summary:   "Add record type " + string(record.ID),
		})
	}
	return SchemaDiff{ID: id, Entries: entries}
}
