package commit

import (
	"strings"
	"unicode/utf8"

	"github.com/destyk/commitkit/internal/lexer"
)

// Parse turns a raw commit message into a structured Message.
func Parse(input string) (Message, error) {
	src := normalize(input)

	if strings.TrimSpace(src) == "" {
		return Message{}, &ParseError{
			Message: "commit message is empty",
		}
	}

	lines := strings.Split(src, "\n")

	header, err := parseHeaderLine(lines[0], 0)
	if err != nil {
		return Message{}, err
	}

	rest := lines[1:]
	footerStart := findFooterStart(rest)

	var bodyLines []string
	var footerLines []string

	if footerStart >= 0 {
		bodyLines = rest[:footerStart]
		footerLines = rest[footerStart:]
	} else {
		bodyLines = rest
	}

	bodyText := strings.TrimSpace(strings.Join(bodyLines, "\n"))
	bodySpan := computeBodySpan(bodyLines, footerStart)

	footers := parseFooters(footerLines, footerStart)

	return Message{
		Raw:     src,
		Header:  header,
		Body:    Body{Text: bodyText, Span: bodySpan},
		Footers: footers,
	}, nil
}

func normalize(input string) string {
	input = strings.ReplaceAll(input, "\r\n", "\n")
	input = strings.ReplaceAll(input, "\r", "\n")
	return strings.TrimRight(input, "\n")
}

func parseHeaderLine(line string, lineOffset int) (Header, error) {
	lex := lexer.New(line)

	var prefixTokens []lexer.Token
	var colonTok lexer.Token
	foundColon := false

	for {
		tok := lex.Next()
		if tok.Kind == lexer.TokenEOF {
			break
		}

		if tok.Kind == lexer.TokenColon {
			colonTok = tok
			foundColon = true
			break
		}

		prefixTokens = append(prefixTokens, tok)
	}

	if !foundColon {
		endCol := utf8.RuneCountInString(line) + 1
		return Header{}, &ParseError{
			Span:    singlePos(pos(lineOffset+len(line), 1, endCol)),
			Message: "header must contain ':'",
		}
	}

	descStartOff := colonTok.Span.End.Offset
	descRaw := line[descStartOff:]
	trimmedDesc := strings.TrimSpace(descRaw)
	if trimmedDesc == "" {
		return Header{}, &ParseError{
			Span: singlePos(pos(
				lineOffset+colonTok.Span.End.Offset,
				1,
				colonTok.Span.End.Column,
			)),
			Message: "description is empty",
		}
	}

	leadingSpaces := len(descRaw) - len(strings.TrimLeft(descRaw, " \t"))
	descStartCol := colonTok.Span.End.Column + leadingSpaces
	descEndCol := descStartCol + utf8.RuneCountInString(trimmedDesc)

	descSpan := span(
		pos(lineOffset+descStartOff+leadingSpaces, 1, descStartCol),
		pos(lineOffset+descStartOff+leadingSpaces+len(trimmedDesc), 1, descEndCol),
	)

	header, err := interpretPrefix(prefixTokens, lineOffset)
	if err != nil {
		return Header{}, err
	}

	header.Description = trimmedDesc
	header.DescSpan = descSpan
	header.Span = span(
		pos(lineOffset, 1, 1),
		pos(lineOffset+len(line), 1, utf8.RuneCountInString(line)+1),
	)

	return header, nil
}

func interpretPrefix(tokens []lexer.Token, lineOffset int) (Header, error) {
	var h Header

	i := 0
	var typeParts []string
	var typeStart, typeEnd lexer.Position
	typeSet := false

	for i < len(tokens) {
		tok := tokens[i]
		if tok.Kind == lexer.TokenLParen || tok.Kind == lexer.TokenBang {
			break
		}
		if !typeSet {
			typeStart = tok.Span.Start
			typeSet = true
		}
		typeEnd = tok.Span.End
		typeParts = append(typeParts, tok.Value)
		i++
	}

	if !typeSet || strings.TrimSpace(strings.Join(typeParts, "")) == "" {
		return Header{}, &ParseError{
			Span:    singlePos(pos(lineOffset, 1, 1)),
			Message: "type is empty",
		}
	}

	h.Type = strings.TrimSpace(strings.Join(typeParts, ""))
	h.TypeSpan = span(toCommitPos(typeStart), toCommitPos(typeEnd))

	if i < len(tokens) && tokens[i].Kind == lexer.TokenLParen {
		openTok := tokens[i]
		i++

		var scopeParts []string
		var scopeStart, scopeEnd lexer.Position
		scopeSet := false

		for i < len(tokens) && tokens[i].Kind != lexer.TokenRParen {
			tok := tokens[i]
			if !scopeSet {
				scopeStart = tok.Span.Start
				scopeSet = true
			}
			scopeEnd = tok.Span.End
			scopeParts = append(scopeParts, tok.Value)
			i++
		}

		if i >= len(tokens) || tokens[i].Kind != lexer.TokenRParen {
			return Header{}, &ParseError{
				Span:    toCommitSpan(openTok.Span),
				Message: "malformed scope: missing ')'",
			}
		}
		closeTok := tokens[i]
		i++

		scopeText := strings.TrimSpace(strings.Join(scopeParts, ""))
		if scopeText == "" {
			return Header{}, &ParseError{
				Span: span(
					toCommitPos(openTok.Span.Start),
					toCommitPos(closeTok.Span.End),
				),
				Message: "scope is empty",
			}
		}

		h.Scope = scopeText
		if scopeSet {
			h.ScopeSpan = span(toCommitPos(scopeStart), toCommitPos(scopeEnd))
		} else {
			h.ScopeSpan = span(toCommitPos(openTok.Span.End), toCommitPos(closeTok.Span.Start))
		}
	}

	if i < len(tokens) && tokens[i].Kind == lexer.TokenBang {
		h.Breaking = true
		h.BreakingSpan = toCommitSpan(tokens[i].Span)
		i++
	}

	if i < len(tokens) {
		return Header{}, &ParseError{
			Span:    toCommitSpan(tokens[i].Span),
			Message: "unexpected token in header prefix",
		}
	}

	return h, nil
}

func findFooterStart(lines []string) int {
	if len(lines) == 0 {
		return -1
	}

	end := len(lines) - 1
	for end >= 0 && lines[end] == "" {
		end--
	}
	if end < 0 {
		return -1
	}

	start := end
	for start >= 0 {
		line := lines[start]
		if line == "" {
			break
		}
		if isPossibleTrailerLine(line) || isContinuationLine(line) {
			start--
			continue
		}
		break
	}
	start++

	if start > end {
		return -1
	}

	if !isTrailerBlock(lines[start : end+1]) {
		return -1
	}
	if start > 0 && lines[start-1] != "" {
		return -1
	}

	return start
}

func isContinuationLine(line string) bool {
	return line != "" && (line[0] == ' ' || line[0] == '\t')
}

func isPossibleTrailerLine(line string) bool {
	if line == "" {
		return false
	}
	colon := strings.IndexByte(line, ':')
	if colon <= 0 {
		return false
	}
	token := strings.TrimSpace(line[:colon])
	if token == "" {
		return false
	}
	if IsBreakingChangeToken(token) {
		return true
	}
	if strings.ContainsAny(token, " \t") {
		return false
	}

	return true
}

func isTrailerBlock(lines []string) bool {
	if len(lines) == 0 {
		return false
	}

	if !isPossibleTrailerLine(lines[0]) {
		return false
	}

	for i := 1; i < len(lines); i++ {
		line := lines[i]
		if line == "" {
			return false
		}
		if isPossibleTrailerLine(line) {
			continue
		}
		if line[0] == ' ' || line[0] == '\t' {
			continue
		}

		return false
	}

	return true
}

func parseFooters(lines []string, footerStartInRest int) []Footer {
	if len(lines) == 0 {
		return nil
	}

	baseLine := 2 + footerStartInRest

	var result []Footer
	var current *Footer
	var currentStartLine int

	flush := func(endLine int) {
		if current == nil {
			return
		}
		startPos := pos(0, currentStartLine, 1)
		endPos := pos(0, endLine, 1)
		current.Span = span(startPos, endPos)
		result = append(result, *current)
		current = nil
	}

	for idx, line := range lines {
		absLine := baseLine + idx

		if isPossibleTrailerLine(line) {
			flush(absLine - 1)

			colon := strings.IndexByte(line, ':')
			token := strings.TrimSpace(line[:colon])
			value := strings.TrimSpace(line[colon+1:])

			current = &Footer{
				Token: token,
				Value: value,
			}
			currentStartLine = absLine

			tokenStartCol := 1 + strings.Index(line, token)
			current.TokenSpan = span(
				pos(0, absLine, tokenStartCol),
				pos(0, absLine, tokenStartCol+len(token)),
			)
			valueStartInLine := colon + 1
			for valueStartInLine < len(line) && (line[valueStartInLine] == ' ' || line[valueStartInLine] == '\t') {
				valueStartInLine++
			}
			current.ValueSpan = span(
				pos(0, absLine, valueStartInLine+1),
				pos(0, absLine, len(line)+1),
			)

			continue
		}

		if current != nil && (line[0] == ' ' || line[0] == '\t') {
			current.Value += "\n" + strings.TrimLeft(line, " \t")
			current.ValueSpan.End = pos(0, absLine, len(line)+1)
			continue
		}
	}

	if current != nil {
		flush(baseLine + len(lines) - 1)
	}

	return result
}

func computeBodySpan(bodyLines []string, footerStart int) Span {
	if len(bodyLines) == 0 {
		return Span{}
	}
	startLine := 2
	endLine := 1 + len(bodyLines)
	if footerStart >= 0 {
		endLine = 1 + footerStart
	}

	return span(pos(0, startLine, 1), pos(0, endLine, 1))
}

// --- position helpers ---

func pos(offset, line, column int) Position {
	return Position{Offset: offset, Line: line, Column: column}
}

func span(start, end Position) Span {
	return Span{Start: start, End: end}
}

func singlePos(p Position) Span {
	return Span{Start: p, End: p}
}

func toCommitPos(p lexer.Position) Position {
	return Position{Offset: p.Offset, Line: p.Line, Column: p.Column}
}

func toCommitSpan(s lexer.Span) Span {
	return Span{
		Start: toCommitPos(s.Start),
		End:   toCommitPos(s.End),
	}
}
