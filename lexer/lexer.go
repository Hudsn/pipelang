package lexer

import (
	"slices"
	"strings"

	"github.com/hudsn/pipelang/token"
)

type Lexer struct {
	input []rune

	currentChar rune

	currentIdx int
	nextIdx    int
}

func (l *Lexer) InputRunes() []rune {
	return l.input
}

func New(input []rune) *Lexer {
	l := &Lexer{
		input: input,
	}
	l.readNext()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.handleWhitespace()

	switch l.currentChar {
	case rune(0):
		tok = newToken(token.EOF, l.currentChar)
		tok.SetPosition(l.currentIdx, l.nextIdx)
	case ';':
		tok = newToken(token.SEMICOLON, l.currentChar)
		tok.SetPosition(l.currentIdx, l.nextIdx)
	case ':':
		tok = newToken(token.COLON, l.currentChar)
		tok.SetPosition(l.currentIdx, l.nextIdx)
	case '.':
		tok = newToken(token.DOT, l.currentChar)
		tok.SetPosition(l.currentIdx, l.nextIdx)
	case ',':
		tok = newToken(token.COMMA, l.currentChar)
		tok.SetPosition(l.currentIdx, l.nextIdx)
	case '|':
		tok = newToken(token.PIPECHAR, l.currentChar)
		tok.SetPosition(l.currentIdx, l.nextIdx)
	case '"':
		tok = l.readString()
		return tok
	case '\'':
		tok = l.readString()
		return tok
	case '=':
		start := l.currentIdx
		tok = l.handleEquals()
		tok.SetPosition(start, l.nextIdx)
	case '>':
		start := l.currentIdx
		tok = l.handleGT()
		tok.SetPosition(start, l.nextIdx)
	case '<':
		start := l.currentIdx
		tok = l.handleLT()
		tok.SetPosition(start, l.nextIdx)
	case '!':
		start := l.currentIdx
		tok = l.handleExclamation()
		tok.SetPosition(start, l.nextIdx)
	case '+':
		tok = newToken(token.PLUS, l.currentChar)
		tok.SetPosition(l.currentIdx, l.nextIdx)
	case '*':
		tok = newToken(token.ASTERISK, l.currentChar)
		tok.SetPosition(l.currentIdx, l.nextIdx)
	case '/':
		tok = newToken(token.SLASH, l.currentChar)
		tok.SetPosition(l.currentIdx, l.nextIdx)
	case '-':
		tok = newToken(token.MINUS, l.currentChar)
		tok.SetPosition(l.currentIdx, l.nextIdx)
	case '$':
		tok = l.readIdentifier()
		return tok
	case '(':
		tok = newToken(token.LPAREN, l.currentChar)
		tok.SetPosition(l.currentIdx, l.nextIdx)
	case ')':
		tok = newToken(token.RPAREN, l.currentChar)
		tok.SetPosition(l.currentIdx, l.nextIdx)
		l.readNext()
		l.maybeAddSemicolon()
		return tok
	case '[':
		tok = newToken(token.LSQUARE, l.currentChar)
		tok.SetPosition(l.currentIdx, l.nextIdx)
	case ']':
		tok = newToken(token.RSQUARE, l.currentChar)
		tok.SetPosition(l.currentIdx, l.nextIdx)
		l.readNext()
		l.maybeAddSemicolon()
		return tok
	case '{':
		tok = newToken(token.LCURLY, l.currentChar)
		tok.SetPosition(l.currentIdx, l.nextIdx)
	case '}':
		tok = newToken(token.RCURLY, l.currentChar)
		tok.SetPosition(l.currentIdx, l.nextIdx)
		l.readNext()
		l.maybeAddSemicolon()
		return tok
	default:
		if isDigit(l.currentChar) {
			tok = l.readNumber()
			return tok
			// NOTE: identifiers can't start with a digit.
		} else if isLetter(l.currentChar) {
			tok = l.readIdentifier()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.currentChar)
			tok.SetPosition(l.currentIdx, l.nextIdx)
		}
	}

	l.readNext()

	return tok
}

func (l *Lexer) readNext() {
	if l.nextIdx >= len(l.input) {
		l.currentChar = rune(0)
	} else {
		l.currentChar = l.input[l.nextIdx]
	}
	l.currentIdx = l.nextIdx
	l.nextIdx += 1
}

func (l *Lexer) peekNext() rune {
	if l.nextIdx >= len(l.input) {
		return rune(0)
	}
	return l.input[l.nextIdx]
}

func newToken(tokenType token.TokenType, char rune) token.Token {
	return token.Token{Type: tokenType, Value: string(char)}
}

// multi-char reader helpers

func (l *Lexer) readNumber() token.Token {
	tok := &token.Token{Type: token.INT}
	startIdx := l.currentIdx
	encounteredDot := false
	for isDigit(l.currentChar) || l.currentChar == '.' {
		if l.currentChar == '.' {
			if !isDigit(l.peekNext()) || encounteredDot {
				break
			}
			tok.Type = token.FLOAT
			encounteredDot = true
		}
		l.readNext()
	}
	tok.SetPosition(startIdx, l.currentIdx)
	tok.Value = string(l.input[startIdx:l.currentIdx])

	// something like 123abc should be illegal.
	if isLetter(l.currentChar) {
		tok.Type = token.ILLEGAL
		return *tok
	}

	l.maybeAddSemicolon()
	return *tok
}

func (l *Lexer) readIdentifier() token.Token {
	tok := &token.Token{}
	startIdx := l.currentIdx
	if l.currentChar == '$' {
		l.readNext()
	}
	for isLetter(l.currentChar) || isRecognizedLineChar(l.currentChar) || isDigit(l.currentChar) {
		l.readNext()
	}

	tok.SetPosition(startIdx, l.currentIdx)
	tok.Value = string(l.input[startIdx:l.currentIdx])
	tok.Type = token.LookupKeyword(tok.Value)

	// if we don't find any $keywords, we return an illegal token since the only valid $ words should be predefined
	if tok.Type == token.IDENT && strings.HasPrefix(tok.Value, "$") {
		tok.Type = token.ILLEGAL
	}

	l.maybeAddSemicolon()
	return *tok
}

func (l *Lexer) readString() token.Token {
	tok := &token.Token{Type: token.STRING}
	endChar := l.currentChar
	startIdx := l.currentIdx
	l.readNext()
	for l.currentChar != endChar && l.currentChar != rune(0) {
		l.readNext()
	}
	tok.SetPosition(startIdx, l.nextIdx)                   // want to capture quotes and contents.
	tok.Value = string(l.input[startIdx+1 : l.currentIdx]) // want to capture only inside quotes here
	l.readNext()                                           // go from end quote to next char
	l.maybeAddSemicolon()
	return *tok
}

func (l *Lexer) handleWhitespace() {
	for slices.Contains([]rune{'\r', '\n', '\t', ' '}, l.currentChar) {
		l.readNext()
	}
}

// Other helpers

// returns the target index unless it is out of bounds. Then it just returns a safe max idx (length of input)
func (l *Lexer) safeIdx(idx int) int {
	if idx >= len(l.input) {
		return len(l.input)
	}
	return idx
}

func (l *Lexer) handleEquals() token.Token {
	tok := &token.Token{
		Type:  token.ASSIGN,
		Value: string(l.currentChar),
	}
	if l.peekNext() == '=' {
		start := l.currentIdx
		l.readNext()
		tok.Type = token.EQ
		tok.Value = string(l.input[start:l.nextIdx])
		return *tok
	}
	return *tok
}
func (l *Lexer) handleExclamation() token.Token {
	tok := &token.Token{
		Type:  token.EXCLAMATION,
		Value: string(l.currentChar),
	}
	if l.peekNext() == '=' {
		start := l.currentIdx
		l.readNext()
		tok.Type = token.NOT_EQ
		tok.Value = string(l.input[start:l.nextIdx])
		return *tok
	}
	return *tok
}

func (l *Lexer) handleGT() token.Token {
	tok := &token.Token{
		Type:  token.GT,
		Value: string(l.currentChar),
	}
	if l.peekNext() == '=' {
		start := l.currentIdx
		l.readNext()
		tok.Type = token.GTEQ
		tok.Value = string(l.input[start:l.nextIdx])
		return *tok
	}
	return *tok
}
func (l *Lexer) handleLT() token.Token {
	tok := &token.Token{
		Type:  token.LT,
		Value: string(l.currentChar),
	}
	if l.peekNext() == '=' {
		start := l.currentIdx
		l.readNext()
		tok.Type = token.LTEQ
		tok.Value = string(l.input[start:l.nextIdx])
		return *tok
	}
	return *tok
}

func (l *Lexer) maybeAddSemicolon() {
	shouldAddSemicolon := false

	for slices.Contains([]rune{'\r', '\n', '\t', ' '}, l.currentChar) {
		if l.currentChar == '\r' && l.peekNext() == '\n' {
			shouldAddSemicolon = true
			break
		}
		if l.currentChar == '\n' {
			shouldAddSemicolon = true
			break
		}
		l.readNext()
	}

	if shouldAddSemicolon { // need to do a bunch of allocations and copying because doing direct array modification was a buggy mess
		newPrefix := make([]rune, len(l.input[:l.currentIdx]))
		rest := make([]rune, len(l.input[l.currentIdx:]))
		copy(newPrefix, l.input[:l.currentIdx])
		newPrefix = append(newPrefix, ';')
		copy(rest, l.input[l.currentIdx:])
		l.input = append(newPrefix, rest...)
		l.currentChar = ';'
		return
	}

	if l.currentChar == rune(0) {
		l.currentChar = ';'
		l.input = append(l.input[:l.currentIdx], ';')
	}
}

func isRecognizedLineChar(char rune) bool {
	return char == '_' || char == '-'
}

func isLetter(char rune) bool {
	return ('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z')
}

func isDigit(char rune) bool {
	return '0' <= char && char <= '9'
}
