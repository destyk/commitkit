// Package lint is the public adapter for the lint domain.
//
// Implementation lives in internal/lint.
package lint

import (
	"github.com/destyk/commitkit/commit"
	intlint "github.com/destyk/commitkit/internal/lint"
)

type (
	Severity  = intlint.Severity
	Violation = intlint.Violation
	Rule      = intlint.Rule
	Rules     = intlint.Rules
	Engine    = intlint.Engine
	Result    = intlint.Result
)

const (
	SeverityError   = intlint.SeverityError
	SeverityWarning = intlint.SeverityWarning
	SeverityInfo    = intlint.SeverityInfo
)

// NewEngine creates an engine with the given rules.
func NewEngine(rules ...Rule) *Engine {
	return intlint.NewEngine(rules...)
}

// Check applies the supplied rules to the message.
func Check(message commit.Message, rules Rules) Result {
	return intlint.Check(message, rules)
}
