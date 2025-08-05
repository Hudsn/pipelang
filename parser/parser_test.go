package parser

import (
	"testing"

	"github.com/hudsn/pipelang/ast"
	"github.com/hudsn/pipelang/lexer"
	"github.com/hudsn/pipelang/utils/testutils"
)

func TestIntegerLiteralExpression(t *testing.T) {
	program := setupTestWithInput(t, "1")
	if len(program.Statements) != 1 {
		t.Fatalf("TestIntegerLiteralExpression: expected len of parsed program to be 1. got=%d", len(program.Statements))
	}

	exp, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("TestIntegerLiteralExpression: program.Statements[0] is not of type *ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	testIntegerLiteral(t, exp.Expression, 1)
}

func testIntegerLiteral(t *testing.T, integerExpression ast.Expression, value int) bool {
	integer, ok := integerExpression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("testIntegerLiteral: passed expression is not an *ast.IntegerLiteral. got type=%T", integerExpression)
	}

	if isEq, failMsg := testutils.Equal(value, integer.Value); !isEq {
		t.Errorf("testIntegerLiteral: value equality check failed: %s", failMsg)
	}

	return true
}

func setupTestWithInput(t *testing.T, input string) *ast.Program {
	lexer := lexer.New([]rune(input))
	parser := New(lexer)
	program, err := parser.ParseProgram()
	if err != nil {
		t.Fatalf("setupTestWithInput: %s", err.Error())
	}
	return program
}
