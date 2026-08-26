package formset

type RecordTypeID string
type FieldID string
type RelationID string
type CapabilityID string
type SchemaVersion string
type Scope string

const (
	ScopeGlobal    Scope = "global"
	ScopeWorkspace Scope = "workspace"
	ScopeTenant    Scope = "tenant"
	ScopeUser      Scope = "user"
)

type FieldType string

const (
	FieldText       FieldType = "text"
	FieldTextarea   FieldType = "textarea"
	FieldRichText   FieldType = "richtext"
	FieldMarkdown   FieldType = "markdown"
	FieldNumber     FieldType = "number"
	FieldBoolean    FieldType = "boolean"
	FieldSelect     FieldType = "select"
	FieldDateTime   FieldType = "datetime"
	FieldJSON       FieldType = "json"
	FieldRelation   FieldType = "relation"
	FieldCollection FieldType = "collection"
	FieldComputed   FieldType = "computed"
	FieldEncrypted  FieldType = "encrypted"
)

type RelationCardinality string

const (
	RelationOneToOne   RelationCardinality = "one-to-one"
	RelationOneToMany  RelationCardinality = "one-to-many"
	RelationManyToMany RelationCardinality = "many-to-many"
)

type DeleteBehavior string

const (
	DeleteRestrict DeleteBehavior = "restrict"
	DeleteCascade  DeleteBehavior = "cascade"
	DeleteNullify  DeleteBehavior = "nullify"
)

type CrossWorkspaceMode string

const (
	CrossWorkspaceForbidden          CrossWorkspaceMode = "forbidden"
	CrossWorkspaceSameProfile        CrossWorkspaceMode = "same-profile"
	CrossWorkspaceExplicitTargets    CrossWorkspaceMode = "explicit-targets"
	CrossWorkspaceRequiresCapability CrossWorkspaceMode = "requires-capability"
)

type RelationPolicy struct {
	CrossWorkspaceMode CrossWorkspaceMode `json:"cross_workspace_mode,omitempty"`
	AllowedTargets     []string           `json:"allowed_targets,omitempty"`
	Capability         CapabilityID       `json:"capability,omitempty"`
	ReadOnly           bool               `json:"read_only,omitempty"`
}

type Option struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type ValidationRule struct {
	Name     string            `json:"name"`
	Message  string            `json:"message,omitempty"`
	Args     map[string]string `json:"args,omitempty"`
	Severity string            `json:"severity,omitempty"`
}

type Field struct {
	ID           FieldID          `json:"id"`
	Label        string           `json:"label"`
	Type         FieldType        `json:"type"`
	Namespace    string           `json:"namespace,omitempty"`
	OwnerModule  string           `json:"owner_module,omitempty"`
	Description  string           `json:"description,omitempty"`
	Required     bool             `json:"required,omitempty"`
	Localized    bool             `json:"localized,omitempty"`
	DefaultValue string           `json:"default_value,omitempty"`
	Options      []Option         `json:"options,omitempty"`
	Rules        []ValidationRule `json:"rules,omitempty"`
	Items        *Field           `json:"items,omitempty"`
	Searchable   bool             `json:"searchable,omitempty"`
	Indexed      bool             `json:"indexed,omitempty"`
	Unique       bool             `json:"unique,omitempty"`
	Sensitive    bool             `json:"sensitive,omitempty"`
	Encrypted    bool             `json:"encrypted,omitempty"`
	UIHint       string           `json:"ui_hint,omitempty"`
	StorageHint  string           `json:"storage_hint,omitempty"`
}

type Relation struct {
	ID                   RelationID          `json:"id"`
	Label                string              `json:"label,omitempty"`
	Source               RecordTypeID        `json:"source"`
	Target               RecordTypeID        `json:"target"`
	Cardinality          RelationCardinality `json:"cardinality"`
	InverseName          string              `json:"inverse_name,omitempty"`
	CrossWorkspacePolicy string              `json:"cross_workspace_policy,omitempty"`
	Policy               RelationPolicy      `json:"policy,omitempty"`
	DeleteBehavior       DeleteBehavior      `json:"delete_behavior,omitempty"`
}

type RecordType struct {
	ID            RecordTypeID   `json:"id"`
	Label         string         `json:"label"`
	Description   string         `json:"description,omitempty"`
	SchemaVersion SchemaVersion  `json:"schema_version,omitempty"`
	OwnerModule   string         `json:"owner_module,omitempty"`
	Scope         Scope          `json:"scope"`
	Fields        []Field        `json:"fields,omitempty"`
	Relations     []Relation     `json:"relations,omitempty"`
	Capabilities  []CapabilityID `json:"capabilities,omitempty"`
	Visibility    string         `json:"visibility,omitempty"`
}
