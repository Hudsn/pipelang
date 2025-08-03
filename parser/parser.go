package parser

import (
	"fmt"

	"github.com/hudsn/pipelang/ast"
	"github.com/hudsn/pipelang/lexer"
	"github.com/hudsn/pipelang/token"
)

type Parser struct {
	lexer *lexer.Lexer

	currentToken token.Token
	peekToken    token.Token
}

func (p *Parser) progressTokens() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}
func (p *Parser) ParseProgram() (*ast.Program, error) {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.isCurrentToken(token.EOF) {
		statement, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		program.Statements = append(program.Statements, statement)
	}

	return program, nil
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.currentToken.Type {
	case token.PIPEDEF:
		// handle pipedef
	case token.PIPECHAR:
		// handle pipe invocation
	case token.IF:
		// handle if statement
	default:
		// handle rest of expressions
	}

	// TODO: remove and just return in the default case
	return nil, nil
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer: l,
	}

	p.progressTokens()
	p.progressTokens()

	return p
}

func (p *Parser) isCurrentToken(tokenType token.TokenType) bool {
	return p.currentToken.Type == tokenType
}

func (p *Parser) isPeekToken(tokenType token.TokenType) bool {
	return p.peekToken.Type == tokenType
}

func (p *Parser) mustNextToken(tokenType token.TokenType) error {
	if p.isPeekToken(tokenType) {
		p.progressTokens()
		return nil
	}
	err := fmt.Errorf("expected next token to be %s. instead got %s", tokenType.HumanString(), p.peekToken.Type.HumanString())
	return newParsingError(err, p.lexer.InputRunes(), p.currentToken)
}

func getLineAndColumn(inputRunes []rune, targetIdx int) (int, int) {
	// 1 indexing b/c its more useful for people-facing errors in editors
	line := 1
	col := 1
	for _, r := range inputRunes[:targetIdx] {
		switch r {
		case '\n': // reset if newline
			line++
			col = 1
		default:
			col++
		}
	}
	return line, col
}

// error

func newParsingError(innerErr error, inputRunes []rune, token token.Token) error {
	start, _ := token.Position.GetPosition()
	line, col := getLineAndColumn(inputRunes, start)
	return fmt.Errorf("parse error at %d:%d:\n\t%w", line, col, innerErr)
}
