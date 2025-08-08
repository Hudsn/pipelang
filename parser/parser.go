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

	errors []error
}

const (
	_ int = iota
	LOWEST
	ARROW          // ~>
	ASSIGN         // =
	LOGIC_OP       // || &&
	EQUALITY       // == !=
	COMPARISON     // < > <= >=
	SUM            // + -
	PRODUCT        // * /
	PREFIX         // -x !x
	CHAIN_CALL_IDX // asdf.fdsa ; func(); arr[]
)

var precedenceMap = map[token.TokenType]int{
	token.ASSIGN:    ASSIGN,
	token.EQ:        EQUALITY,
	token.NOT_EQ:    EQUALITY,
	token.LT:        COMPARISON,
	token.LTEQ:      COMPARISON,
	token.GT:        COMPARISON,
	token.GTEQ:      COMPARISON,
	token.PLUS:      SUM,
	token.MINUS:     SUM,
	token.ASTERISK:  PRODUCT,
	token.SLASH:     PRODUCT,
	token.LPAREN:    CHAIN_CALL_IDX,
	token.DOT:       CHAIN_CALL_IDX,
	token.LSQUARE:   CHAIN_CALL_IDX,
	token.ARROW:     ARROW,
	token.LOGIC_AND: LOGIC_OP,
	token.LOGIC_OR:  LOGIC_OP,
}

type prefixFunc func() ast.Expression
type infixFunc func(left ast.Expression) ast.Expression

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:           l,
		prefixFunctions: make(map[token.TokenType]prefixFunc),
		infixFunctions:  make(map[token.TokenType]infixFunc),
		errors:          []error{},
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
	p.registerPrefixFunc(token.IF, p.parseIfExpression)
	p.registerPrefixFunc(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefixFunc(token.MINUS, p.parsePrefixExpression)
	p.registerPrefixFunc(token.EXCLAMATION, p.parsePrefixExpression)

	p.registerInfixFunc(token.PLUS, p.parseInfixExpression)
	p.registerInfixFunc(token.MINUS, p.parseInfixExpression)
	p.registerInfixFunc(token.ASTERISK, p.parseInfixExpression)
	p.registerInfixFunc(token.SLASH, p.parseInfixExpression)
	p.registerInfixFunc(token.LOGIC_OR, p.parseInfixExpression)
	p.registerInfixFunc(token.LOGIC_AND, p.parseInfixExpression)
	p.registerInfixFunc(token.EQ, p.parseInfixExpression)
	p.registerInfixFunc(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfixFunc(token.LT, p.parseInfixExpression)
	p.registerInfixFunc(token.LTEQ, p.parseInfixExpression)
	p.registerInfixFunc(token.GT, p.parseInfixExpression)
	p.registerInfixFunc(token.GTEQ, p.parseInfixExpression)
	p.registerInfixFunc(token.ARROW, p.parseArrowFunctionExpression)
	p.registerInfixFunc(token.LPAREN, p.parseCallExpression)
	// TODO
	// p.registerInfixFunc(token.DOT, p.parseDotAccessExpression)
	// p.registerInfixFunc(token.LSQUARE, p.parseIndexExpression)

}

func (p *Parser) progressTokens() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
	if p.isCurrentToken(token.ILLEGAL) {
		err := newParsingError(fmt.Errorf("illegal token"), p.lexer.InputRunes(), p.currentToken)
		p.errors = append(p.errors, err)
	}
}
func (p *Parser) ParseProgram() (*ast.Program, error) {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.isCurrentToken(token.EOF) && !p.isCurrentToken(token.ILLEGAL) {
		statement := p.parseStatement()
		program.Statements = append(program.Statements, statement)
		p.progressTokens()
	}
	if p.isCurrentToken(token.ILLEGAL) {
		return nil, newParsingError(fmt.Errorf("illegal token"), p.lexer.InputRunes(), p.currentToken)
	}

	return program, nil
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.IDENT:
		if p.isPeekToken(token.ASSIGN) {
			ident := p.parseIdentifier()
			p.progressTokens()
			return p.parseAssignStatement(ident)
		}
		return p.parseExpressionStatement()
	case token.PIPEDEF:
		// TODO handle pipedef
		return nil
	case token.PIPECHAR:
		// TODO handle pipe invocation
		return nil
	default:
		// handle rest of expressions
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	var leftExpression ast.Expression
	prefixFn := p.prefixFunctions[p.currentToken.Type]
	if prefixFn == nil {
		p.errUnexpected()
		return nil
	}

	leftExpression = prefixFn()

	for precedence < p.peekPrecedence() {
		infixFn := p.infixFunctions[p.peekToken.Type]
		if infixFn == nil {
			return leftExpression
		}
		p.progressTokens()
		leftExpression = infixFn(leftExpression)
	}
	return leftExpression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{Token: p.currentToken, Left: left, Operator: p.currentToken.Value}

	precedence, ok := precedenceMap[p.currentToken.Type]
	if !ok {
		precedence = LOWEST
	}

	p.progressTokens()

	exp.Right = p.parseExpression(precedence)

	return exp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{Token: p.currentToken, Operator: p.currentToken.Value}
	p.progressTokens()
	exp.Right = p.parseExpression(PREFIX)
	return exp
}

// literals and specific parsers

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currentToken}
	expr := p.parseExpression(LOWEST)

	stmt.Expression = expr

	if p.isPeekToken(token.SEMICOLON) {
		p.progressTokens()
	}

	return stmt
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	ret := &ast.IntegerLiteral{Token: p.currentToken}
	val, err := strconv.ParseInt(p.currentToken.Value, 0, 64)
	if err != nil {
		parseErr := fmt.Errorf("parse integer: %q is not an integer", p.currentToken.Value)
		err := newParsingError(parseErr, p.lexer.InputRunes(), p.currentToken)
		p.errors = append(p.errors, err)
	}
	ret.Value = int(val)
	return ret
}

func (p *Parser) ParseFloatLiteral() ast.Expression {
	ret := &ast.FloatLiteral{Token: p.currentToken}
	val, err := strconv.ParseFloat(p.currentToken.Value, 64)
	if err != nil {
		parseErr := fmt.Errorf("parse float: %q is not a float", p.currentToken.Value)
		err := newParsingError(parseErr, p.lexer.InputRunes(), p.currentToken)
		p.errors = append(p.errors, err)
	}
	ret.Value = val
	return ret
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}
}

func (p *Parser) parseBoolean() ast.Expression {
	var val bool
	switch p.currentToken.Value {
	case "true":
		val = true
	case "false":
		val = false
	default:
		err := fmt.Errorf("parse boolean: %q is not a boolean", p.currentToken.Value)
		err = newParsingError(err, p.lexer.InputRunes(), p.currentToken)
		p.errors = append(p.errors, err)
	}
	return &ast.Boolean{Token: p.currentToken, Value: val}
}

func (p *Parser) parseString() ast.Expression {
	return &ast.StringLiteral{Token: p.currentToken, Value: p.currentToken.Value}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.progressTokens()
	exp := p.parseExpression(LOWEST)
	if !p.mustNextToken(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseCallExpression(left ast.Expression) ast.Expression {
	ret := &ast.CallExpression{Token: p.currentToken}
	ident, ok := left.(*ast.Identifier)
	if !ok {
		p.errUnexpectedToken(left.GetToken())
		return nil
	}
	ret.Name = ident

	ret.Arguments = p.parseExpressionList(token.RPAREN)
	return ret
}

func (p *Parser) parseExpressionList(endType token.TokenType) []ast.Expression {
	// enter function still on opening character
	// for example we are still on the '(' in: (arg1, arg2, arg3)

	// end early if empty enclosure chars like ()
	// make sure to end on closing char
	if p.isPeekToken(endType) {
		p.progressTokens()
		return []ast.Expression{}
	}

	p.progressTokens() // now on first substantive expr entry

	ret := []ast.Expression{p.parseExpression(LOWEST)}

	for p.isPeekToken(token.COMMA) {
		p.progressTokens() // now at comma
		p.progressTokens() // skip comma to next entry
		ret = append(ret, p.parseExpression(LOWEST))
	}

	if !p.mustNextToken(endType) {
		return nil
	}

	return ret
}

func (p *Parser) parseArrowFunctionExpression(left ast.Expression) ast.Expression {
	ret := &ast.ArrowFunctionExpression{Token: p.currentToken}

	ident, ok := left.(*ast.Identifier)
	if !ok {
		p.errUnexpectedToken(left.GetToken())
		return nil
	}
	ret.Param = ident

	p.progressTokens()

	ret.QueryExpression = p.parseExpression(LOWEST)

	return ret
}

func (p *Parser) parseIfExpression() ast.Expression {
	ret := &ast.IfExpression{Token: p.currentToken}

	p.progressTokens()
	condition := p.parseExpression(LOWEST)
	ret.Condition = condition

	if !p.mustNextToken(token.LCURLY) {
		return nil
	}

	consequence := p.parseBlockStatement()

	ret.Consequence = consequence

	if p.isPeekToken(token.ELSE) {
		p.progressTokens()
		p.progressTokens()

		switch p.currentToken.Type {
		case token.IF:
			alt := p.parseIfExpression()
			ret.Alternative = alt
		case token.LCURLY:
			alt := p.parseBlockStatement()
			ret.Alternative = alt
		default:
			p.errUnexpected()
		}
	}

	if p.isPeekToken(token.SEMICOLON) {
		p.progressTokens()
	}
	return ret
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	ret := &ast.BlockStatement{OpenToken: p.currentToken}
	ret.Statements = []ast.Statement{}
	p.progressTokens()
	for !p.isCurrentToken(token.RCURLY) && !p.isCurrentToken(token.EOF) {
		statement := p.parseStatement()
		if statement != nil {
			ret.Statements = append(ret.Statements, statement)
		}
		p.progressTokens()
	}
	if !p.isCurrentToken(token.RCURLY) {
		p.errUnexpected()
		return nil
	}
	ret.CloseToken = p.currentToken
	return ret
}

func (p *Parser) parseAssignStatement(expr ast.Expression) *ast.AssignStatement {
	ident, ok := expr.(*ast.Identifier)
	if !ok {
		err := fmt.Errorf("expect a valid identifier on the left side of assign statement")
		p.errors = append(p.errors, newParsingError(err, p.lexer.InputRunes(), p.currentToken))
	}
	ret := &ast.AssignStatement{
		Token: p.currentToken,
		Name:  ident,
	}

	p.progressTokens()

	ret.Value = p.parseExpression(LOWEST)

	if p.isPeekToken(token.SEMICOLON) {
		p.progressTokens()
	}
	return ret
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

func (p *Parser) mustNextToken(tokenType token.TokenType) bool {
	if p.isPeekToken(tokenType) {
		p.progressTokens()
		return true
	}
	p.errUnexpectedToken(p.peekToken)
	return false
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

// generic error for unexpected sequences (missing operator funcs, parsing statements where the order is incorrect, etc...)
func (p *Parser) errUnexpected() {
	start, end := p.currentToken.Position.GetPosition()
	str := string(p.lexer.InputRunes()[start:end])
	if p.currentToken.Type == token.EOF {
		str = "EOF"
	}
	e := fmt.Errorf("unexpected sequence: %s", str)
	p.errors = append(p.errors, newParsingError(e, p.lexer.InputRunes(), p.currentToken))
}

// same as errUnexpected, but for a specific token in our lexed output.
// useful for flagging tokens we've alredy progressed past, like in left sides of infix expressions
func (p *Parser) errUnexpectedToken(t token.Token) {
	start, end := t.Position.GetPosition()
	str := string(p.lexer.InputRunes()[start:end])
	if t.Type == token.EOF {
		str = "EOF"
	}
	e := fmt.Errorf("unexpected sequence: %s", str)
	p.errors = append(p.errors, newParsingError(e, p.lexer.InputRunes(), t))
}

func newParsingError(innerErr error, inputRunes []rune, token token.Token) error {
	start, _ := token.Position.GetPosition()
	line, col := getLineAndColumn(inputRunes, start)
	return fmt.Errorf("parse error at %d:%d:\n\t%w", line, col, innerErr)
}

func (p *Parser) CheckParserErrors() error {
	if len(p.errors) > 0 {
		return p.errors[0]
	}
	return nil
}
