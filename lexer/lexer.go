package lexer

import (
	"slices"

	"github.com/hudsn/pipelang/token"
)

type Lexer struct {
	input []rune

	currentChar rune

	currentIdx int
	nextIdx    int

	lineNum int
	colNum  int
}

func New(input []rune) *Lexer {
	l := &Lexer{
		input: input,
	}
	l.readNext()
	l.lineNum = 1
	l.colNum = 1
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.handleWhitespace()

	switch l.currentChar {
	case ';':
		tok = newToken(token.SEMICOLON, l.currentChar)
		tok.SetPosition(l.lineNum, l.colNum)
	case '.':
		tok = newToken(token.DOT, l.currentChar)
		tok.SetPosition(l.lineNum, l.colNum)
	case '"':
	case '\'':

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
		}
	}

	l.readNext()

	return tok
}

func (l *Lexer) readNext() {
	if l.currentChar == '\n' {
		// any newline should reset col to 1 and increment the line count
		l.colNum = 1
		l.lineNum += 1
	} else {
		// otherwise we just prog the char count of the current line
		l.colNum += 1
	}
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
	return token.Token{Type: tokenType, Literal: string(char)}
}

// multi-char reader helpers

func (l *Lexer) readNumber() token.Token {
	tok := &token.Token{Type: token.INT}
	tok.SetPosition(l.lineNum, l.colNum)
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
	literal := string(l.input[startIdx:l.currentIdx])
	tok.Literal = literal

	l.maybeAddSemicolon()
	return *tok
}

func (l *Lexer) readIdentifier() token.Token {
	tok := &token.Token{Type: token.IDENT}
	tok.SetPosition(l.lineNum, l.colNum)
	start := l.currentIdx

	for isLetter(l.currentChar) || isRecognizedLineChar(l.currentChar) {
		l.readNext()
	}

	tok.Literal = string(l.input[start:l.currentIdx])
	return *tok
}

func (l *Lexer) readString() token.Token //TODO

func (l *Lexer) handleWhitespace() {
	for slices.Contains([]rune{'\r', '\n', '\t', ' '}, l.currentChar) {
		l.readNext()
	}
}

// Other helpers

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
