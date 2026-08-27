package lexer

// TokenKind identifies the kind of a token produced by the lexer.
type TokenKind int

const (
	TokenEOF TokenKind = iota
	TokenNewline
	TokenText   // any non-newline run of characters
	TokenColon  // ':'
	TokenBang   // '!'
	TokenLParen // '('
	TokenRParen // ')'
	TokenSpace  // one or more spaces / tabs
)

func (k TokenKind) String() string {
	switch k {
	case TokenEOF:
		return "EOF"
	case TokenNewline:
		return "NEWLINE"
	case TokenText:
		return "TEXT"
	case TokenColon:
		return "COLON"
	case TokenBang:
		return "BANG"
	case TokenLParen:
		return "LPAREN"
	case TokenRParen:
		return "RPAREN"
	case TokenSpace:
		return "SPACE"
	default:
		return "UNKNOWN"
	}
}

// Position describes a location in the source text.
type Position struct {
	Offset int
	Line   int
	Column int
}

// Span covers a continuous region of the source.
type Span struct {
	Start Position
	End   Position
}

// Token is a single lexical unit with source location.
type Token struct {
	Kind  TokenKind
	Value string
	Span  Span
}

// Lexer turns a normalized commit message into a stream of tokens.
type Lexer struct {
	src    string
	offset int
	line   int
	column int
}

// New creates a lexer over already-normalized input.
func New(src string) *Lexer {
	return &Lexer{
		src:    src,
		offset: 0,
		line:   1,
		column: 1,
	}
}

// Next returns the next token. At end of input it returns TokenEOF.
func (l *Lexer) Next() Token {
	if l.offset >= len(l.src) {
		p := Position{Offset: l.offset, Line: l.line, Column: l.column}
		return Token{Kind: TokenEOF, Span: Span{Start: p, End: p}}
	}

	startOff := l.offset
	startLine := l.line
	startCol := l.column
	ch := l.src[l.offset]

	switch ch {
	case '\n':
		l.advance()
		return Token{
			Kind: TokenNewline,
			Span: Span{
				Start: Position{startOff, startLine, startCol},
				End:   Position{l.offset, l.line, l.column},
			},
		}

	case ':':
		l.advance()
		return Token{
			Kind:  TokenColon,
			Value: ":",
			Span: Span{
				Start: Position{startOff, startLine, startCol},
				End:   Position{l.offset, l.line, l.column},
			},
		}

	case '!':
		l.advance()
		return Token{
			Kind:  TokenBang,
			Value: "!",
			Span: Span{
				Start: Position{startOff, startLine, startCol},
				End:   Position{l.offset, l.line, l.column},
			},
		}

	case '(':
		l.advance()
		return Token{
			Kind:  TokenLParen,
			Value: "(",
			Span: Span{
				Start: Position{startOff, startLine, startCol},
				End:   Position{l.offset, l.line, l.column},
			},
		}

	case ')':
		l.advance()
		return Token{
			Kind:  TokenRParen,
			Value: ")",
			Span: Span{
				Start: Position{startOff, startLine, startCol},
				End:   Position{l.offset, l.line, l.column},
			},
		}

	case ' ', '\t':
		for l.offset < len(l.src) {
			c := l.src[l.offset]
			if c != ' ' && c != '\t' {
				break
			}
			l.advance()
		}
		return Token{
			Kind:  TokenSpace,
			Value: l.src[startOff:l.offset],
			Span: Span{
				Start: Position{startOff, startLine, startCol},
				End:   Position{l.offset, l.line, l.column},
			},
		}

	default:
		for l.offset < len(l.src) {
			c := l.src[l.offset]
			if c == '\n' || c == ':' || c == '!' ||
				c == '(' || c == ')' || c == ' ' || c == '\t' {
				break
			}
			l.advance()
		}
		return Token{
			Kind:  TokenText,
			Value: l.src[startOff:l.offset],
			Span: Span{
				Start: Position{startOff, startLine, startCol},
				End:   Position{l.offset, l.line, l.column},
			},
		}
	}
}

// Peek returns the next token without advancing.
func (l *Lexer) Peek() Token {
	savedOff := l.offset
	savedLine := l.line
	savedCol := l.column
	tok := l.Next()
	l.offset = savedOff
	l.line = savedLine
	l.column = savedCol
	return tok
}

func (l *Lexer) advance() {
	if l.offset >= len(l.src) {
		return
	}

	if l.src[l.offset] == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}

	l.offset++
}

// Snapshot captures the current lexer state.
type Snapshot struct {
	offset int
	line   int
	column int
}

// Snapshot returns the current state for backtracking.
func (l *Lexer) Snapshot() Snapshot {
	return Snapshot{offset: l.offset, line: l.line, column: l.column}
}

// Restore restores a previously taken snapshot.
func (l *Lexer) Restore(s Snapshot) {
	l.offset = s.offset
	l.line = s.line
	l.column = s.column
}
