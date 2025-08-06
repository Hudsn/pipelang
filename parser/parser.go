package parser

import (
	"fmt"
	"strconv"

	"github.com/hudsn/pipelang/ast"
	"github.com/hudsn/pipelang/lexer"
	"github.com/hudsn/pipelang/token"
)

type Parser struct {
	lexer *lexer.Lexer

	currentToken token.Token
	peekToken    token.Token

	prefixFunctions map[token.TokenType]prefixFunc
	infixFunctions  map[token.TokenType]infixFunc
}

const (
	_ int = iota
	LOWEST
	EQUALITY // == != < > <= >=
	SUM      // + -
	PRODUCT  // * /
	PREFIX   // -x !x
	CHAIN    // asdf.fdsa
	CALL     // func()
	INDEX    // []
)

var precedenceMap = map[token.TokenType]int{
	token.EQ:       EQUALITY,
	token.LT:       EQUALITY,
	token.LTEQ:     EQUALITY,
	token.GT:       EQUALITY,
	token.GTEQ:     EQUALITY,
	token.NOT_EQ:   EQUALITY,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
	token.LPAREN:   CALL,
	token.DOT:      CHAIN,
	token.LSQUARE:  INDEX,
}

type prefixFunc func() (ast.Expression, error)
type infixFunc func(left ast.Expression) (ast.Expression, error)

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:           l,
		prefixFunctions: make(map[token.TokenType]prefixFunc),
		infixFunctions:  make(map[token.TokenType]infixFunc),
	}

	p.registerFuncs()

	p.progressTokens()
	p.progressTokens()
	return p
}

func (p *Parser) registerFuncs() {
	p.registerPrefixFunc(token.INT, p.parseIntegerLiteral)
	p.registerPrefixFunc(token.IDENT, p.parseIdentifier)
}

func (p *Parser) progressTokens() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}
func (p *Parser) ParseProgram() (*ast.Program, error) {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.isCurrentToken(token.EOF) {
		statement, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		program.Statements = append(program.Statements, statement)
		p.progressTokens()
	}

	return program, nil
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.currentToken.Type {
	case token.PIPEDEF:
		// handle pipedef
		return nil, nil
	case token.PIPECHAR:
		// handle pipe invocation
		return nil, nil
	case token.IF:
		// handle if statement
		return nil, nil
	default:
		// handle rest of expressions
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() (*ast.ExpressionStatement, error) {
	stmt := &ast.ExpressionStatement{Token: p.currentToken}
	expr, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	stmt.Expression = expr

	if p.isPeekToken(token.SEMICOLON) {
		p.progressTokens()
	}

	return stmt, nil
}

func (p *Parser) parseExpression(precedence int) (ast.Expression, error) {
	var leftExpression ast.Expression
	var err error
	prefixFn := p.prefixFunctions[p.currentToken.Type]
	if prefixFn == nil {
		err := fmt.Errorf("prefix function not found for token type %q", p.currentToken.Type.HumanString())
		return nil, newParsingError(err, p.lexer.InputRunes(), p.currentToken)
	}

	if leftExpression, err = prefixFn(); err != nil {
		return nil, err
	}

	for precedence < p.peekPrecedence() {
		infixFn := p.infixFunctions[p.peekToken.Type]
		if infixFn == nil {
			err := fmt.Errorf("infix function not found for token type %s", p.peekToken.Type.HumanString())
			return nil, newParsingError(err, p.lexer.InputRunes(), p.currentToken)
		}
		p.progressTokens()
		if leftExpression, err = infixFn(leftExpression); err != nil {
			return nil, err
		}
	}

	return leftExpression, nil
}

func (p *Parser) parseIntegerLiteral() (ast.Expression, error) {
	ret := &ast.IntegerLiteral{Token: p.currentToken}
	val, err := strconv.ParseInt(p.currentToken.Value, 0, 64)
	if err != nil {
		parseErr := fmt.Errorf("failed to parse %q as an integer", p.currentToken.Value)
		return nil, newParsingError(parseErr, p.lexer.InputRunes(), p.currentToken)
	}
	ret.Value = int(val)
	return ret, nil
}

func (p *Parser) parseIdentifier() (ast.Expression, error) {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}, nil
}

func (p *Parser) registerPrefixFunc(tokenType token.TokenType, fn prefixFunc) {
	p.prefixFunctions[tokenType] = fn
}
func (p *Parser) registerInfixFunc(tokenType token.TokenType, fn infixFunc) {
	p.infixFunctions[tokenType] = fn
}

//HELPERS

func (p *Parser) peekPrecedence() int {
	if precedence, ok := precedenceMap[p.peekToken.Type]; ok {
		return precedence
	}
	return LOWEST
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
