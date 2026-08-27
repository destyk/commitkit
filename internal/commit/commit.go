package commit

// Message is the structured representation of a commit message.
// Every significant part carries a Span so rules and diagnostics
// can point at exact source locations.
type Message struct {
	// Raw is the normalized source text that was parsed.
	Raw string

	Header  Header
	Body    Body
	Footers []Footer
}

// Header is the first non-empty line of a commit message.
type Header struct {
	Type        string
	Scope       string
	Breaking    bool
	Description string

	// Spans
	Span         Span // whole header line
	TypeSpan     Span
	ScopeSpan    Span // empty if no scope
	BreakingSpan Span // the '!' if present
	DescSpan     Span
}

// Body holds the free-form body text that sits between the header
// and the trailer block.
type Body struct {
	Text string
	Span Span
}

// Footer represents a single Git trailer (token: value).
// Value may contain newlines when continuation lines are present.
type Footer struct {
	Token string
	Value string

	// Spans
	Span      Span // whole trailer including continuations
	TokenSpan Span
	ValueSpan Span
}

// HasBreakingFooter reports whether any footer is a BREAKING CHANGE
// trailer (case-insensitive, both space and hyphen forms).
func (m Message) HasBreakingFooter() bool {
	for _, f := range m.Footers {
		if IsBreakingChangeToken(f.Token) {
			return true
		}
	}

	return false
}

// IsBreakingChangeToken reports whether token is a BREAKING CHANGE
// trailer name (space or hyphen form, case-insensitive).
func IsBreakingChangeToken(token string) bool {
	switch token {
	case "BREAKING CHANGE", "BREAKING-CHANGE":
		return true
	}

	upper := ToUpperASCII(token)

	return upper == "BREAKING CHANGE" || upper == "BREAKING-CHANGE"
}

// ToUpperASCII uppercases ASCII letters only.
func ToUpperASCII(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		b[i] = c
	}

	return string(b)
}
