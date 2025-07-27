package lexer

import (
	"fmt"
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
		input:   input,
		lineNum: 1,
		colNum:  1,
	}
	l.readNext()

	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.handleWhitespace()

	switch l.currentChar {
	case rune(0):
	case ';':
		tok = newToken(token.SEMICOLON, l.currentChar)
		tok.SetPosition(l.lineNum, l.colNum)
		l.readNext()
	case '.':
		tok = newToken(token.DOT, l.currentChar)
		tok.SetPosition(l.lineNum, l.colNum)
	case '"':
	case '\'':

	default:
		if isDigit(l.currentChar) {
			tok = l.readNumber()
		} else if isLetter(l.currentChar) {
			// tok = l.readIdentifier()
			tok = newToken(token.ILLEGAL, l.currentChar)
		} else {
			fmt.Println(string(l.currentChar))
			tok = newToken(token.ILLEGAL, l.currentChar)
		}
	}
	return tok
}

func (l *Lexer) readNext() {
	if l.currentChar == '\n' {
		// any newline should reset col to 1 and increment the line count
		l.colNum = 1
		l.lineNum++
	} else if l.currentIdx == l.nextIdx {
		// only true if we just initialized our lexer
		l.colNum = 1
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
	return *tok
}

func (l *Lexer) handleWhitespace() {
	for slices.Contains([]rune{'\r', '\n', '\t', ' '}, l.currentChar) {
		l.readNext()
	}
}

func (l *Lexer) maybeAddSemicolon() {
	shouldAddSemicolon := false

	for slices.Contains([]rune{'\r', '\n', '\t', ' '}, l.currentChar) {
		switch l.currentChar {
		case '\r':
			if l.peekNext() == '\n' {
				shouldAddSemicolon = true
			}
		case '\n':
			shouldAddSemicolon = true
		}
	}
	if !shouldAddSemicolon {
		return
	}

	//TODO: add semi in current char spot and decrement col counter to compensate for disparity between input and added semis

}

func isRecognizedLineChar(char rune) bool {
	// dashes and underscores
	return char == '_' || char == '-'
}

func isLetter(char rune) bool {
	return ('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z')
}

func isDigit(char rune) bool {
	return '0' <= char && char <= '9'
}
