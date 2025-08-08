package ast

import (
	"fmt"
	"strings"

	"github.com/hudsn/pipelang/token"
)

//TODO:
// add function call

// add index expr

// add dot access expression
// token, object (left) expression, property (right) identifier

type Node interface {
	Position() token.Position
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

//
// top-level program
//

type Program struct {
	Statements []Statement
}

func (p *Program) Position() token.Position {
	if len(p.Statements) > 0 { //return the whole scannable program input
		first, _ := getNodePositions(p.Statements[0])
		_, last := getNodePositions(p.Statements[len(p.Statements)])
		return newPosition(first, last)
	}
	return token.NullPosition
}
func (p *Program) String() string {
	stmts := []string{}
	for _, s := range p.Statements {
		stmts = append(stmts, s.String())
	}

	return strings.Join(stmts, "\n")
}

//
//statements: expressionstatement(meta), if, pipedef, pipechar
//

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) Position() token.Position {
	first, _ := es.Token.Position.GetPosition()
	_, last := getNodePositions(es.Expression)
	return newPosition(first, last)
}
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

//

type BlockStatement struct {
	OpenToken  token.Token
	CloseToken token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) Position() token.Position {
	first, _ := bs.OpenToken.Position.GetPosition()
	_, last := bs.CloseToken.Position.GetPosition()
	return newPosition(first, last)
}
func (bs *BlockStatement) String() string {
	stmts := []string{}
	for _, s := range bs.Statements {
		stmts = append(stmts, s.String())
	}

	return strings.Join(stmts, "\n")
}

type AssignStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (as *AssignStatement) statementNode() {}
func (as *AssignStatement) Position() token.Position {
	startPos := as.Name.Position()
	start, _ := startPos.GetPosition()
	endPos := as.Value.Position()
	_, end := endPos.GetPosition()
	pos := &token.Position{}
	pos.SetPosition(start, end)
	return *pos
}
func (as *AssignStatement) String() string {
	return fmt.Sprintf("%s = %s", as.Name.String(), as.Value.String())
}

//
// expressions
//

type IntegerLiteral struct {
	Token token.Token
	Value int
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) Position() token.Position {
	return il.Token.Position
}
func (il *IntegerLiteral) String() string { return il.Token.Value }

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (f *FloatLiteral) expressionNode() {}
func (f *FloatLiteral) Position() token.Position {
	return f.Token.Position
}
func (f *FloatLiteral) String() string {
	return f.Token.Value
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) Position() token.Position {
	return i.Token.Position
}
func (i *Identifier) String() string {
	return i.Token.Value
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}
func (b *Boolean) Position() token.Position {
	return b.Token.Position
}
func (b *Boolean) String() string {
	return b.Token.Value
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (s *StringLiteral) expressionNode() {}
func (s *StringLiteral) Position() token.Position {
	return s.Token.Position
}
func (s *StringLiteral) String() string {
	return fmt.Sprintf(`"%s"`, s.Value)
}

type ArrowFunctionExpression struct {
	Token           token.Token
	Params          []Identifier
	QueryExpression Expression
	start           int
}

func (f *ArrowFunctionExpression) expressionNode() {}
func (f *ArrowFunctionExpression) Position() token.Position {
	bpos := f.QueryExpression.Position()
	_, end := bpos.GetPosition()
	pos := &token.Position{}
	pos.SetPosition(f.start, end)
	return *pos
}
func (f *ArrowFunctionExpression) String() string {
	params := ""
	switch len(f.Params) {
	case 0:
		return ""
	case 1:
		params = f.Params[0].String()
	default:
		paramList := []string{}
		for _, p := range f.Params {
			paramList = append(paramList, p.String())
		}
		params = fmt.Sprintf("( %s )", strings.Join(paramList, ", "))
	}
	return fmt.Sprintf("%s -> %s", params, f.QueryExpression.String())
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative Node // block statement or ifExpression
}

func (i *IfExpression) expressionNode() {}
func (i *IfExpression) Position() token.Position {
	first, _ := i.Token.Position.GetPosition()
	var last int
	if i.Alternative == nil {
		_, last = getNodePositions(i.Consequence)
		return newPosition(first, last)
	}
	return newPosition(getNodePositions(i.Alternative))
}
func (i *IfExpression) String() string {
	ret := fmt.Sprintf("if %s { %s }", i.Condition.String(), i.Consequence.String())
	if i.Alternative != nil {
		switch alt := i.Alternative.(type) {
		case *BlockStatement:
			ret += fmt.Sprintf(" else { %s }", alt.String())
		case *IfExpression:
			ret += fmt.Sprintf(" else %s", alt.String())
		}
	}
	return ret
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) Position() token.Position {
	start, _ := pe.Token.Position.GetPosition()
	rightPos := pe.Right.Position()
	_, end := rightPos.GetPosition()
	pos := &token.Position{}
	pos.SetPosition(start, end)
	return *pos
}
func (pe *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", pe.Operator, pe.Right.String())
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) Position() token.Position {
	leftPos := ie.Left.Position()
	rightPos := ie.Right.Position()
	start, _ := leftPos.GetPosition()
	_, end := rightPos.GetPosition()
	pos := &token.Position{}
	pos.SetPosition(start, end)
	return *pos
}
func (ie *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Left.String(), ie.Operator, ie.Right.String())
}

//
// helpers
//

func newPosition(start, end int) token.Position {
	pos := token.Position{}
	pos.SetPosition(start, end)
	return pos
}

func getNodePositions(node Node) (int, int) {
	pos := node.Position()
	return pos.GetPosition()
}
