package ast

import (
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

//
//

type Program struct {
	Statements []Statement
}

func (p *Program) Position() token.Position {
	if len(p.Statements) > 0 {
		return p.Statements[0].Position()
	}
	return token.NullPosition
}

//
//

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (ex *ExpressionStatement) Position() token.Position
