// Package commit is the public adapter for the commit domain.
//
// The implementation lives in internal/commit. This package re-exports
// types and functions so external consumers depend only on stable
// domain paths, while the core remains free of outward dependencies.
package commit

import (
	intcommit "github.com/destyk/commitkit/internal/commit"
)

// Core types.
type (
	Message  = intcommit.Message
	Header   = intcommit.Header
	Body     = intcommit.Body
	Footer   = intcommit.Footer
	Position = intcommit.Position
	Span     = intcommit.Span
	ParseError = intcommit.ParseError
)

// Parse turns a raw commit message into a structured Message.
func Parse(input string) (Message, error) {
	return intcommit.Parse(input)
}

// IsBreakingChangeToken reports whether token is a BREAKING CHANGE trailer name.
func IsBreakingChangeToken(token string) bool {
	return intcommit.IsBreakingChangeToken(token)
}

// ToUpperASCII uppercases ASCII letters only.
func ToUpperASCII(s string) string {
	return intcommit.ToUpperASCII(s)
}
