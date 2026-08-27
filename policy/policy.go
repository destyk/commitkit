// Package policy is the public adapter for commit policies.
//
// Implementation lives in internal/policy.
package policy

import (
	"github.com/destyk/commitkit/lint"
	intpolicy "github.com/destyk/commitkit/internal/policy"
)

// ConventionalCommits returns the default Conventional Commits policy.
func ConventionalCommits() lint.Rules {
	return intpolicy.ConventionalCommits()
}

// TypeEnum restricts the allowed commit types.
func TypeEnum(types ...string) lint.Rule {
	return intpolicy.TypeEnum(types...)
}

// DescriptionLength limits the description length in runes.
func DescriptionLength(min, max int) lint.Rule {
	return intpolicy.DescriptionLength(min, max)
}

// DescriptionLowercase requires the first letter of the description to be lowercase.
func DescriptionLowercase() lint.Rule {
	return intpolicy.DescriptionLowercase()
}

// BreakingChangeFooter requires a BREAKING CHANGE footer when the header has '!'.
func BreakingChangeFooter() lint.Rule {
	return intpolicy.BreakingChangeFooter()
}

// Custom creates a policy from arbitrary rules.
func Custom(rules ...lint.Rule) lint.Rules {
	return intpolicy.Custom(rules...)
}

// RequireScope requires a non-empty scope.
func RequireScope() lint.Rule {
	return intpolicy.RequireScope()
}

// HeaderMaxLength limits the total length of the header line.
func HeaderMaxLength(max int) lint.Rule {
	return intpolicy.HeaderMaxLength(max)
}

// NoTrailingPeriod forbids a trailing '.' in the description.
func NoTrailingPeriod() lint.Rule {
	return intpolicy.NoTrailingPeriod()
}

// ScopeEnum restricts the allowed scope values.
func ScopeEnum(scopes ...string) lint.Rule {
	return intpolicy.ScopeEnum(scopes...)
}
