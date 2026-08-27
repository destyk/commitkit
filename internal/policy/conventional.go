package policy

import (
	"strings"
	"unicode/utf8"

	"github.com/destyk/commitkit/internal/commit"
	"github.com/destyk/commitkit/internal/lint"
)

// ConventionalCommits returns the default Conventional Commits policy.
func ConventionalCommits() lint.Rules {
	return lint.Rules{
		TypeEnum(
			"feat",
			"fix",
			"docs",
			"style",
			"refactor",
			"perf",
			"test",
			"build",
			"ci",
			"chore",
			"revert",
		),
		DescriptionLength(1, 72),
		DescriptionLowercase(),
		BreakingChangeFooter(),
	}
}

// ---------------------------------------------------------------------------
// type-enum
// ---------------------------------------------------------------------------

type typeRule struct {
	allowed map[string]struct{}
}

// TypeEnum restricts the allowed commit types.
func TypeEnum(types ...string) lint.Rule {
	allowed := make(map[string]struct{}, len(types))
	for _, typ := range types {
		allowed[typ] = struct{}{}
	}

	return typeRule{allowed: allowed}
}

func (r typeRule) Name() string {
	return "type-enum"
}

func (r typeRule) Check(message commit.Message) []lint.Violation {
	if _, ok := r.allowed[message.Header.Type]; ok {
		return nil
	}

	return []lint.Violation{{
		Rule:     r.Name(),
		Message:  "type is not allowed",
		Severity: lint.SeverityError,
		Span:     message.Header.TypeSpan,
		Line:     message.Header.TypeSpan.Start.Line,
		Column:   message.Header.TypeSpan.Start.Column,
	}}
}

// ---------------------------------------------------------------------------
// description-length
// ---------------------------------------------------------------------------

type descriptionLengthRule struct {
	min int
	max int
}

// DescriptionLength limits the description length in runes.
// A zero value disables the corresponding limit.
func DescriptionLength(min, max int) lint.Rule {
	return descriptionLengthRule{min: min, max: max}
}

func (r descriptionLengthRule) Name() string {
	return "description-length"
}

func (r descriptionLengthRule) Check(message commit.Message) []lint.Violation {
	length := utf8.RuneCountInString(message.Header.Description)

	if r.min > 0 && length < r.min {
		return []lint.Violation{{
			Rule:     r.Name(),
			Message:  "description is too short",
			Severity: lint.SeverityError,
			Span:     message.Header.DescSpan,
			Line:     message.Header.DescSpan.Start.Line,
			Column:   message.Header.DescSpan.Start.Column,
		}}
	}

	if r.max > 0 && length > r.max {
		return []lint.Violation{{
			Rule:     r.Name(),
			Message:  "description exceeds maximum length",
			Severity: lint.SeverityError,
			Span:     message.Header.DescSpan,
			Line:     message.Header.DescSpan.Start.Line,
			Column:   message.Header.DescSpan.Start.Column,
		}}
	}

	return nil
}

// ---------------------------------------------------------------------------
// description-lowercase
// ---------------------------------------------------------------------------

type descriptionLowercaseRule struct{}

// DescriptionLowercase requires the first letter of the description
// to be lowercase (ASCII).
func DescriptionLowercase() lint.Rule {
	return descriptionLowercaseRule{}
}

func (descriptionLowercaseRule) Name() string {
	return "description-lowercase"
}

func (descriptionLowercaseRule) Check(message commit.Message) []lint.Violation {
	desc := message.Header.Description
	if desc == "" {
		return nil
	}

	r, _ := utf8.DecodeRuneInString(desc)
	if r >= 'A' && r <= 'Z' {
		return []lint.Violation{{
			Rule:     "description-lowercase",
			Message:  "description must start with a lowercase letter",
			Severity: lint.SeverityError,
			Span:     message.Header.DescSpan,
			Line:     message.Header.DescSpan.Start.Line,
			Column:   message.Header.DescSpan.Start.Column,
		}}
	}

	return nil
}

// ---------------------------------------------------------------------------
// breaking-change-footer
// ---------------------------------------------------------------------------

type breakingChangeFooterRule struct{}

// BreakingChangeFooter requires a BREAKING CHANGE (or BREAKING-CHANGE)
// footer when the header is marked with '!'.
func BreakingChangeFooter() lint.Rule {
	return breakingChangeFooterRule{}
}

func (breakingChangeFooterRule) Name() string {
	return "breaking-change-footer"
}

func (breakingChangeFooterRule) Check(message commit.Message) []lint.Violation {
	if !message.Header.Breaking {
		return nil
	}

	if message.HasBreakingFooter() {
		return nil
	}

	return []lint.Violation{{
		Rule:     "breaking-change-footer",
		Message:  "breaking change requires BREAKING CHANGE footer",
		Severity: lint.SeverityError,
		Span:     message.Header.BreakingSpan,
		Line:     message.Header.BreakingSpan.Start.Line,
		Column:   message.Header.BreakingSpan.Start.Column,
	}}
}

// ---------------------------------------------------------------------------
// helpers used by other rules
// ---------------------------------------------------------------------------

func hasFooterToken(message commit.Message, tokens ...string) bool {
	set := make(map[string]struct{}, len(tokens))
	for _, t := range tokens {
		set[strings.ToUpper(t)] = struct{}{}
	}
	for _, f := range message.Footers {
		if _, ok := set[strings.ToUpper(f.Token)]; ok {
			return true
		}
	}
	return false
}
