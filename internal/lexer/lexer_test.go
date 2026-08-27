package lexer

import "testing"

func TestLexerBasic(t *testing.T) {
	lex := New("feat(api)!: hello")
	var kinds []TokenKind
	for {
		tok := lex.Next()
		kinds = append(kinds, tok.Kind)
		if tok.Kind == TokenEOF {
			break
		}
	}

	expected := []TokenKind{
		TokenText,   // feat
		TokenLParen, // (
		TokenText,   // api
		TokenRParen, // )
		TokenBang,   // !
		TokenColon,  // :
		TokenSpace,  //
		TokenText,   // hello
		TokenEOF,
	}

	if len(kinds) != len(expected) {
		t.Fatalf("got %v, want %v", kinds, expected)
	}

	for i := range expected {
		if kinds[i] != expected[i] {
			t.Fatalf("token %d: got %v, want %v", i, kinds[i], expected[i])
		}
	}
}

func TestLexerNewlines(t *testing.T) {
	lex := New("a\nb")
	tok := lex.Next()

	if tok.Kind != TokenText || tok.Value != "a" {
		t.Fatalf("first = %v", tok)
	}

	tok = lex.Next()
	if tok.Kind != TokenNewline {
		t.Fatalf("expected newline, got %v", tok)
	}

	tok = lex.Next()
	if tok.Kind != TokenText || tok.Value != "b" {
		t.Fatalf("third = %v", tok)
	}
}
