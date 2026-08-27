package commit

import "fmt"

// ParseError describes a syntactic problem in a commit message.
type ParseError struct {
	Span    Span
	Message string
}

func (e *ParseError) Error() string {
	if !e.Span.IsValid() {
		return "parse error: " + e.Message
	}

	if e.Span.Start.Line == e.Span.End.Line &&
		e.Span.Start.Column == e.Span.End.Column {
		return fmt.Sprintf(
			"parse error at %d:%d: %s",
			e.Span.Start.Line,
			e.Span.Start.Column,
			e.Message,
		)
	}

	return fmt.Sprintf(
		"parse error at %s: %s",
		e.Span.String(),
		e.Message,
	)
}
