package lint

import "github.com/destyk/commitkit/internal/commit"

// Rule validates a parsed commit message and returns zero or more
// violations. Rules are pure: they do not perform I/O and do not
// mutate the message.
type Rule interface {
	Name() string
	Check(commit.Message) []Violation
}

// Rules is an ordered collection of rules. Order determines the order
// of emitted diagnostics.
type Rules []Rule

// Engine is a small rule engine that applies a set of rules to a
// message and aggregates the results. It is intentionally simple –
// no dependency graph, no parallel execution – so it stays easy to
// reason about and free of external packages.
type Engine struct {
	rules Rules
}

// NewEngine creates an engine with the given rules.
func NewEngine(rules ...Rule) *Engine {
	return &Engine{rules: Rules(rules)}
}

// Rules returns a copy of the configured rules.
func (e *Engine) Rules() Rules {
	out := make(Rules, len(e.rules))
	copy(out, e.rules)
	return out
}

// Add appends rules to the engine (useful for building policies
// incrementally).
func (e *Engine) Add(rules ...Rule) {
	e.rules = append(e.rules, rules...)
}

// Result contains all violations found by the rules.
type Result struct {
	Violations []Violation
}

// Check applies all rules of the engine to the message.
func (e *Engine) Check(message commit.Message) Result {
	return Check(message, e.rules)
}

// Check applies the supplied rules to the message.
// This is the functional entry point used by the policy package.
func Check(message commit.Message, rules Rules) Result {
	var violations []Violation
	for _, rule := range rules {
		violations = append(violations, rule.Check(message)...)
	}

	return Result{Violations: violations}
}

// Valid reports whether no violations were found.
func (r Result) Valid() bool {
	return len(r.Violations) == 0
}

// HasErrors reports whether at least one error-level violation exists.
func (r Result) HasErrors() bool {
	for _, v := range r.Violations {
		if v.Severity == SeverityError {
			return true
		}
	}

	return false
}

// Errors returns only error-level violations.
func (r Result) Errors() []Violation {
	var out []Violation
	for _, v := range r.Violations {
		if v.Severity == SeverityError {
			out = append(out, v)
		}
	}

	return out
}

// Warnings returns only warning-level violations.
func (r Result) Warnings() []Violation {
	var out []Violation
	for _, v := range r.Violations {
		if v.Severity == SeverityWarning {
			out = append(out, v)
		}
	}

	return out
}

// Filter returns a new Result containing only violations that satisfy
// the predicate.
func (r Result) Filter(pred func(Violation) bool) Result {
	var out []Violation
	for _, v := range r.Violations {
		if pred(v) {
			out = append(out, v)
		}
	}

	return Result{Violations: out}
}
