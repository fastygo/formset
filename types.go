package formset

import "github.com/fastygo/codex/schema"

// The record vocabulary is owned by github.com/fastygo/codex/schema.
// formset re-exports it as type aliases so callers can use either name;
// there is no formset-owned copy of these structs (see ADR-0001).
type (
	RecordTypeID        = schema.RecordTypeID
	FieldID             = schema.FieldID
	RelationID          = schema.RelationID
	CapabilityID        = schema.CapabilityID
	SchemaVersion       = schema.SchemaVersion
	Scope               = schema.Scope
	FieldType           = schema.FieldType
	RelationCardinality = schema.RelationCardinality
	DeleteBehavior      = schema.DeleteBehavior
	CrossWorkspaceMode  = schema.CrossWorkspaceMode
	RelationPolicy      = schema.RelationPolicy
	Option              = schema.Option
	ValidationRule      = schema.ValidationRule
	Field               = schema.Field
	Relation            = schema.Relation
	RecordType          = schema.RecordType
)

const (
	ScopeGlobal    = schema.ScopeGlobal
	ScopeWorkspace = schema.ScopeWorkspace
	ScopeTenant    = schema.ScopeTenant
	ScopeUser      = schema.ScopeUser
)

const (
	FieldText       = schema.FieldText
	FieldString     = schema.FieldString
	FieldTextarea   = schema.FieldTextarea
	FieldRichText   = schema.FieldRichText
	FieldMarkdown   = schema.FieldMarkdown
	FieldNumber     = schema.FieldNumber
	FieldBoolean    = schema.FieldBoolean
	FieldSelect     = schema.FieldSelect
	FieldDateTime   = schema.FieldDateTime
	FieldJSON       = schema.FieldJSON
	FieldObject     = schema.FieldObject
	FieldRelation   = schema.FieldRelation
	FieldCollection = schema.FieldCollection
	FieldComputed   = schema.FieldComputed
	FieldEncrypted  = schema.FieldEncrypted
)

const (
	RelationOneToOne   = schema.RelationOneToOne
	RelationOneToMany  = schema.RelationOneToMany
	RelationManyToMany = schema.RelationManyToMany
)

const (
	DeleteRestrict = schema.DeleteRestrict
	DeleteCascade  = schema.DeleteCascade
	DeleteNullify  = schema.DeleteNullify
)

// Codex scalar profile rules (fastygo.codex/*). Money, integer, decimal,
// date, uri and uuid are rules on a base Codex field type, not types.
const (
	RuleInteger  = schema.RuleInteger
	RuleDecimal  = schema.RuleDecimal
	RuleMoney    = schema.RuleMoney
	RuleDate     = schema.RuleDate
	RuleURI      = schema.RuleURI
	RuleUUID     = schema.RuleUUID
	RuleNullable = schema.RuleNullable
	RuleReadOnly = schema.RuleReadOnly
	RuleJSONAny  = schema.RuleJSONAny
)

// Slot UI hints. Codex owns UIHintMedia; the editor hints are formset slots.
// Renderers must not infer them from Type alone.
const (
	UIHintMedia    = schema.UIHintMedia
	UIHintTipTap   = "tiptap"
	UIHintMarkdown = "markdown"
	UIHintTextarea = "textarea"
)

// Cross-workspace modes. Codex v0.4.2 declares the CrossWorkspaceMode type
// but no constants; these values are typed as the Codex type and kept here
// until Codex exports its own (documented gap, see CHANGELOG).
const (
	CrossWorkspaceForbidden          CrossWorkspaceMode = "forbidden"
	CrossWorkspaceSameProfile        CrossWorkspaceMode = "same-profile"
	CrossWorkspaceExplicitTargets    CrossWorkspaceMode = "explicit-targets"
	CrossWorkspaceRequiresCapability CrossWorkspaceMode = "requires-capability"
)
