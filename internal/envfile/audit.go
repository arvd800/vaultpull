package envfile

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// AuditEntry records a single change event for an env key.
type AuditEntry struct {
	Timestamp time.Time
	Key       string
	Action    string // "added", "removed", "changed"
	Source    string // e.g. vault path
}

func (e AuditEntry) String() string {
	return fmt.Sprintf("%s  %-10s  %-20s  %s",
		e.Timestamp.Format(time.RFC3339), e.Action, e.Key, e.Source)
}

// AuditLog holds a list of audit entries.
type AuditLog []AuditEntry

// BuildAuditLog creates an AuditLog from a Diff result.
func BuildAuditLog(d DiffResult, source string) AuditLog {
	now := time.Now().UTC()
	var log AuditLog

	for _, k := range sorted(d.Added) {
		log = append(log, AuditEntry{Timestamp: now, Key: k, Action: "added", Source: source})
	}
	for _, k := range sorted(d.Removed) {
		log = append(log, AuditEntry{Timestamp: now, Key: k, Action: "removed", Source: source})
	}
	for _, k := range sorted(d.Changed) {
		log = append(log, AuditEntry{Timestamp: now, Key: k, Action: "changed", Source: source})
	}
	return log
}

// Format returns a human-readable multi-line string of the audit log.
func (al AuditLog) Format() string {
	if len(al) == 0 {
		return "(no changes)"
	}
	lines := make([]string, len(al))
	for i, e := range al {
		lines[i] = e.String()
	}
	return strings.Join(lines, "\n")
}

func sorted(keys []string) []string {
	out := make([]string, len(keys))
	copy(out, keys)
	sort.Strings(out)
	return out
}
