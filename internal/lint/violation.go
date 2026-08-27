package lint

import "github.com/destyk/commitkit/internal/commit"

// Severity describes the severity of a violation.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
)

// Violation is a single lint diagnostic.
type Violation struct {
	Rule     string
	Message  string
	Severity Severity
	// Prefer Span when available; Line/Column kept for backward compatibility.
	Span   commit.Span
	Line   int
	Column int
}

// Position returns a human-readable location string.
func (v Violation) Position() string {
	if v.Span.IsValid() {
		return v.Span.String()
	}

	if v.Line > 0 && v.Column > 0 {
		return sprintf("%d:%d", v.Line, v.Column)
	}

	if v.Line > 0 {
		return sprintf("%d", v.Line)
	}

	return ""
}

func sprintf(format string, a ...any) string {
	return formatString(format, a...)
}
