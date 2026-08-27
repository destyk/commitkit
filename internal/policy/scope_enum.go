package policy

import (
	"github.com/destyk/commitkit/internal/commit"
	"github.com/destyk/commitkit/internal/lint"
)

// ScopeEnum restricts the allowed scope values.
// An empty scope is always rejected by this rule; combine with
// projects that omit RequireScope only if you accept commits without
// a scope (this rule simply does not fire when Scope is empty).
func ScopeEnum(scopes ...string) lint.Rule {
	allowed := make(map[string]struct{}, len(scopes))
	for _, s := range scopes {
		allowed[s] = struct{}{}
	}

	return scopeEnumRule{allowed: allowed}
}

type scopeEnumRule struct {
	allowed map[string]struct{}
}

func (r scopeEnumRule) Name() string {
	return "scope-enum"
}

func (r scopeEnumRule) Check(message commit.Message) []lint.Violation {
	if message.Header.Scope == "" {
		return nil
	}

	if _, ok := r.allowed[message.Header.Scope]; ok {
		return nil
	}

	return []lint.Violation{{
		Rule:     r.Name(),
		Message:  "scope is not allowed",
		Severity: lint.SeverityError,
		Span:     message.Header.ScopeSpan,
		Line:     message.Header.ScopeSpan.Start.Line,
		Column:   message.Header.ScopeSpan.Start.Column,
	}}
}
