package policy

import (
	"github.com/destyk/commitkit/internal/commit"
	"github.com/destyk/commitkit/internal/lint"
)

// Custom creates a policy from arbitrary rules.
func Custom(rules ...lint.Rule) lint.Rules {
	return lint.Rules(rules)
}

// RequireScope creates a rule that requires a non-empty scope.
func RequireScope() lint.Rule {
	return requireScopeRule{}
}

type requireScopeRule struct{}

func (requireScopeRule) Name() string {
	return "scope-required"
}

func (requireScopeRule) Check(message commit.Message) []lint.Violation {
	if message.Header.Scope != "" {
		return nil
	}

	return []lint.Violation{{
		Rule:     "scope-required",
		Message:  "scope is required",
		Severity: lint.SeverityError,
		Span:     message.Header.Span,
		Line:     1,
		Column:   1,
	}}
}

// HeaderMaxLength limits the total length of the header line.
func HeaderMaxLength(max int) lint.Rule {
	return headerMaxLengthRule{max: max}
}

type headerMaxLengthRule struct {
	max int
}

func (r headerMaxLengthRule) Name() string {
	return "header-max-length"
}

func (r headerMaxLengthRule) Check(message commit.Message) []lint.Violation {
	// Reconstruct approximate header length from the original raw
	// when available; otherwise fall back to type+scope+desc.
	length := 0
	if message.Raw != "" {
		// First line length.
		for i, c := range message.Raw {
			if c == '\n' {
				length = i
				break
			}
			length = i + 1
		}
	} else {
		length = len(message.Header.Type)
		if message.Header.Scope != "" {
			length += 2 + len(message.Header.Scope) // (scope)
		}

		if message.Header.Breaking {
			length++
		}

		length += 2 + len(message.Header.Description) // ": "
	}

	if length > r.max {
		return []lint.Violation{{
			Rule:     r.Name(),
			Message:  "header exceeds maximum length",
			Severity: lint.SeverityError,
			Span:     message.Header.Span,
			Line:     1,
		}}
	}

	return nil
}

// NoTrailingPeriod forbids a trailing '.' in the description.
func NoTrailingPeriod() lint.Rule {
	return noTrailingPeriodRule{}
}

type noTrailingPeriodRule struct{}

func (noTrailingPeriodRule) Name() string {
	return "no-trailing-period"
}

func (noTrailingPeriodRule) Check(message commit.Message) []lint.Violation {
	desc := message.Header.Description
	if len(desc) > 0 && desc[len(desc)-1] == '.' {
		return []lint.Violation{{
			Rule:     "no-trailing-period",
			Message:  "description must not end with a period",
			Severity: lint.SeverityError,
			Span:     message.Header.DescSpan,
			Line:     message.Header.DescSpan.Start.Line,
		}}
	}

	return nil
}
