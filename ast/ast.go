package ast

import (
	"fmt"

	"github.com/hudsn/pipelang/token"
)

type Node interface {
	Position() token.Position
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type ExpressionDebugger interface {
	Expression
	String() string
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

//
//statements: expressionstatement(meta), if, pipedef, pipechar
//

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (ex *ExpressionStatement) Position() token.Position {
	first, _ := ex.Token.Position.GetPosition()
	_, last := getNodePositions(ex.Expression)
	return newPosition(first, last)
}

//

type IfStatement struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative Statement // block statement or ifstatement
}

func (i *IfStatement) statementNode() {}
func (i *IfStatement) Position() token.Position {
	first, _ := i.Token.Position.GetPosition()
	var last int
	if i.Alternative == nil {
		_, last = getNodePositions(i.Consequence)
		return newPosition(first, last)
	}
	return newPosition(getNodePositions(i.Alternative))
}

type BlockStatement struct {
	OpenToken  token.Token
	CloseToken token.Token // we'll need to populate this based on expectToken in parser func for block statements
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) Position() token.Position {
	first, _ := bs.OpenToken.Position.GetPosition()
	_, last := bs.CloseToken.Position.GetPosition()
	return newPosition(first, last)
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
