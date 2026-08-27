package commit

import "fmt"

// Position describes a location in the source text.
// Line and Column are 1-based. Offset is 0-based byte offset.
type Position struct {
	Offset int // byte offset from start of input
	Line   int // 1-based
	Column int // 1-based
}

// Span covers a continuous region of the source.
type Span struct {
	Start Position
	End   Position
}

// IsValid reports whether the span has a non-zero start line.
func (s Span) IsValid() bool {
	return s.Start.Line > 0
}

// String returns a human-readable representation, e.g. "1:5-1:12".
func (s Span) String() string {
	if !s.IsValid() {
		return ""
	}

	if s.Start.Line == s.End.Line && s.Start.Column == s.End.Column {
		return fmt.Sprintf("%d:%d", s.Start.Line, s.Start.Column)
	}

	if s.Start.Line == s.End.Line {
		return fmt.Sprintf("%d:%d-%d", s.Start.Line, s.Start.Column, s.End.Column)
	}

	return fmt.Sprintf(
		"%d:%d-%d:%d",
		s.Start.Line,
		s.Start.Column,
		s.End.Line,
		s.End.Column,
	)
}
