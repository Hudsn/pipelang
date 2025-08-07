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
	p.registerPrefixFunc(token.FLOAT, p.ParseFloatLiteral)
	p.registerPrefixFunc(token.TRUE, p.parseBoolean)
	p.registerPrefixFunc(token.FALSE, p.parseBoolean)
	p.registerPrefixFunc(token.IDENT, p.parseIdentifier)
	p.registerPrefixFunc(token.STRING, p.parseString)
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
		return p.parseIfStatement()
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
		err := fmt.Errorf("prefix function not found for token type: %s", p.currentToken.Type.HumanString())
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

func (p *Parser) ParseFloatLiteral() (ast.Expression, error) {
	ret := &ast.FloatLiteral{Token: p.currentToken}
	val, err := strconv.ParseFloat(p.currentToken.Value, 64)
	if err != nil {
		parseErr := fmt.Errorf("failed to parse %q as a float", p.currentToken.Value)
		return nil, newParsingError(parseErr, p.lexer.InputRunes(), p.currentToken)
	}
	ret.Value = val
	return ret, nil
}

func (p *Parser) parseIdentifier() (ast.Expression, error) {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}, nil
}

func (p *Parser) parseBoolean() (ast.Expression, error) {
	var val bool
	switch p.currentToken.Value {
	case "true":
		val = true
	case "false":
		val = false
	default:
		err := fmt.Errorf("failed to parse %q as a boolean", p.currentToken.Value)
		return nil, newParsingError(err, p.lexer.InputRunes(), p.currentToken)
	}
	return &ast.Boolean{Token: p.currentToken, Value: val}, nil
}

func (p *Parser) parseString() (ast.Expression, error) {
	return &ast.StringLiteral{Token: p.currentToken, Value: p.currentToken.Value}, nil
}

func (p *Parser) parseIfStatement() (ast.Statement, error) {
	ret := &ast.IfStatement{Token: p.currentToken}

	p.progressTokens()
	condition, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	ret.Condition = condition

	if err := p.mustNextToken(token.LCURLY); err != nil {
		return nil, err
	}

	consequence, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	ret.Consequence = consequence

	if p.isPeekToken(token.ELSE) {
		p.progressTokens()
		p.progressTokens()

		switch p.currentToken.Type {
		case token.IF:
			alt, err := p.parseIfStatement()
			if err != nil {
				return nil, err
			}
			ret.Alternative = alt
		case token.LCURLY:
			alt, err := p.parseBlockStatement()
			if err != nil {
				return nil, err
			}
			ret.Alternative = alt
		default:
			err := fmt.Errorf("unexpected token type. wanted %s or %s. got %s", token.IF.HumanString(), token.ELSE.HumanString(), p.currentToken.Type.HumanString())
			return nil, newParsingError(err, p.lexer.InputRunes(), p.currentToken)
		}
	}

	if p.isPeekToken(token.SEMICOLON) {
		p.progressTokens()
	}
	return ret, nil
}

func (p *Parser) parseBlockStatement() (*ast.BlockStatement, error) {
	ret := &ast.BlockStatement{OpenToken: p.currentToken}
	ret.Statements = []ast.Statement{}
	p.progressTokens()
	for !p.isCurrentToken(token.RCURLY) && !p.isCurrentToken(token.EOF) {
		statement, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if statement != nil {
			ret.Statements = append(ret.Statements, statement)
		}
		p.progressTokens()
	}
	if p.isCurrentToken(token.EOF) {
		e := fmt.Errorf("unexpected end-of-file (EOF)")
		return nil, newParsingError(e, p.lexer.InputRunes(), p.currentToken)
	}
	if p.isCurrentToken(token.RCURLY) {
		ret.CloseToken = p.currentToken
	}
	return ret, nil
}

//HELPERS

func (p *Parser) registerPrefixFunc(tokenType token.TokenType, fn prefixFunc) {
	p.prefixFunctions[tokenType] = fn
}
func (p *Parser) registerInfixFunc(tokenType token.TokenType, fn infixFunc) {
	p.infixFunctions[tokenType] = fn
}

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
	err := errUnexpectedTokenType(tokenType, p.peekToken.Type)
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

func errUnexpectedTokenType(want, got token.TokenType) error {
	return fmt.Errorf("expected next token to be %s. instead got %s", want.HumanString(), got.HumanString())
}

func newParsingError(innerErr error, inputRunes []rune, token token.Token) error {
	start, _ := token.Position.GetPosition()
	line, col := getLineAndColumn(inputRunes, start)
	return fmt.Errorf("parse error at %d:%d:\n\t%w", line, col, innerErr)
}
