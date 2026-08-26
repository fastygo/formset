package formset

import (
	"fmt"
	"strings"
)

type ReviewSeverity string

const (
	ReviewInfo    ReviewSeverity = "info"
	ReviewWarning ReviewSeverity = "warning"
	ReviewError   ReviewSeverity = "error"
)

type ReviewFinding struct {
	Severity ReviewSeverity `json:"severity"`
	Code     string         `json:"code"`
	Message  string         `json:"message"`
	RecordID RecordTypeID   `json:"record_id,omitempty"`
	FieldID  FieldID        `json:"field_id,omitempty"`
	Relation RelationID     `json:"relation,omitempty"`
}

type ReviewReport struct {
	Findings []ReviewFinding `json:"findings"`
}

func (report ReviewReport) HasErrors() bool {
	for _, finding := range report.Findings {
		if finding.Severity == ReviewError {
			return true
		}
	}
	return false
}

func (report ReviewReport) Summary() string {
	if len(report.Findings) == 0 {
		return "schema review passed"
	}
	parts := make([]string, 0, len(report.Findings))
	for _, finding := range report.Findings {
		parts = append(parts, fmt.Sprintf("%s:%s:%s", finding.Severity, finding.Code, finding.Message))
	}
	return strings.Join(parts, "\n")
}

func ReviewSchema(records []RecordType, relations []Relation, diffs ...SchemaDiff) ReviewReport {
	report := ReviewReport{}
	recordIDs := map[RecordTypeID]struct{}{}
	for _, record := range records {
		if _, exists := recordIDs[record.ID]; exists {
			report.Findings = append(report.Findings, ReviewFinding{
				Severity: ReviewError,
				Code:     "duplicate-record-type",
				Message:  "record type is declared more than once",
				RecordID: record.ID,
			})
		}
		recordIDs[record.ID] = struct{}{}
		report.Findings = append(report.Findings, reviewRecord(record)...)
		relations = append(relations, record.Relations...)
	}
	for _, relation := range relations {
		if _, ok := recordIDs[relation.Source]; !ok {
			report.Findings = append(report.Findings, ReviewFinding{
				Severity: ReviewError, Code: "missing-relation-source",
				Message: "relation source record is missing", Relation: relation.ID,
			})
		}
		if _, ok := recordIDs[relation.Target]; !ok {
			report.Findings = append(report.Findings, ReviewFinding{
				Severity: ReviewError, Code: "missing-relation-target",
				Message: "relation target record is missing", Relation: relation.ID,
			})
		}
		if relation.Policy.CrossWorkspaceMode == "" && relation.CrossWorkspacePolicy != "" {
			report.Findings = append(report.Findings, ReviewFinding{
				Severity: ReviewWarning, Code: "legacy-cross-workspace-policy",
				Message: "relation uses legacy string cross-workspace policy", Relation: relation.ID,
			})
		}
	}
	for _, diff := range diffs {
		report.Findings = append(report.Findings, reviewDiff(diff)...)
	}
	return report
}

func reviewRecord(record RecordType) []ReviewFinding {
	var findings []ReviewFinding
	fields := map[FieldID]Field{}
	for _, field := range record.Fields {
		if existing, exists := fields[field.ID]; exists {
			findings = append(findings, ReviewFinding{
				Severity: ReviewError,
				Code:     "duplicate-field",
				Message:  fmt.Sprintf("field is declared more than once; first type %q, duplicate type %q", existing.Type, field.Type),
				RecordID: record.ID,
				FieldID:  field.ID,
			})
		}
		fields[field.ID] = field
		if field.Encrypted && !field.Sensitive {
			findings = append(findings, ReviewFinding{
				Severity: ReviewWarning,
				Code:     "encrypted-not-sensitive",
				Message:  "encrypted field should usually also be marked sensitive",
				RecordID: record.ID,
				FieldID:  field.ID,
			})
		}
	}
	return findings
}

func reviewDiff(diff SchemaDiff) []ReviewFinding {
	var findings []ReviewFinding
	for _, entry := range diff.Entries {
		switch entry.Operation {
		case DiffRemoveField, DiffChangeField:
			if entry.Breaking {
				findings = append(findings, ReviewFinding{
					Severity: ReviewWarning,
					Code:     "breaking-schema-change",
					Message:  "diff entry is marked breaking",
					RecordID: entry.RecordID,
					FieldID:  entry.FieldID,
				})
			}
		}
	}
	return findings
}
