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
		lineNum: 0,
	}
	l.readNext()

	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.handleWhitespace(false)

	switch l.currentChar {
	case rune(0):
	case ';':
		tok = newToken(token.SEMICOLON, l.currentChar)
		tok.SetPositions(l.lineNum, l.colNum)
		l.readNext()
	case '.':
		tok = newToken(token.DOT, l.currentChar)
		tok.SetPositions(l.lineNum, l.colNum)
	case '"':
	case '\'':

	default:
		if isDigit(l.currentChar) {
			tok = l.readNumber()
			l.handleWhitespace(true)
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
		// any newline should reset col to 0 and increment the line count
		l.colNum = 0
		l.lineNum++
	} else if l.currentIdx == l.nextIdx {
		// only true if we just initialized our lexer
		l.colNum = 0
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

// func (l *Lexer) readIdentifier() string {

// }

func (l *Lexer) readNumber() token.Token {
	tok := &token.Token{Type: token.INT}
	tok.SetStart(l.lineNum, l.colNum)
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
	tok.SetEnd(l.lineNum, l.colNum-1)
	return *tok
}

func (l *Lexer) handleWhitespace(addMissingSemicolons bool) {
	for slices.Contains([]rune{'\r', '\n', '\t', ' '}, l.currentChar) {

		if addMissingSemicolons && (l.currentChar != '\n' || l.currentChar != '\r') {
			l.maybeAddSemicolon()
		}
		l.readNext()
	}
}

func (l *Lexer) maybeAddSemicolon() {
	shouldAddSemi := false
	if l.currentChar == '\r' && l.peekNext() == '\n' {
		shouldAddSemi = true
	}
	if l.currentChar == '\n' {
		shouldAddSemi = true
	}
	if !shouldAddSemi {
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
