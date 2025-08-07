package parser

import (
	"testing"

	"github.com/hudsn/pipelang/ast"
	"github.com/hudsn/pipelang/lexer"
	"github.com/hudsn/pipelang/utils/testutils"
)

func TestIfStatement(t *testing.T) {

	input := `if true {
				1
			} else if false {
				2
				3
			} else {
				4
				5
				6
			}`

	program := setupTestWithInput(t, input)
	if len(program.Statements) != 1 {
		t.Fatalf("expected len of parsed program to be 1. got=%d", len(program.Statements))
	}
	ifStmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.IfStatement. got=%T", program.Statements[0])
	}
	testBooleanLiteral(t, ifStmt.Condition, true)
	if len(ifStmt.Consequence.Statements) != 1 {
		t.Fatalf("expected len of consequence block of if-statement to be 1. got=%d", len(ifStmt.Consequence.Statements))
	}
	consequence, ok := ifStmt.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ifStmt.Consequence.Statements[0] is not *ast.ExpressionStatement. got=%T", ifStmt.Consequence.Statements[0])
	}
	testExpression(t, consequence.Expression, 1)
	altIfStmt, ok := ifStmt.Alternative.(*ast.IfStatement)
	if !ok {
		t.Fatalf("ifStmt.Alternative is not *ast.IfStatement. got=%T", ifStmt.Alternative)
	}
	testBooleanLiteral(t, altIfStmt.Condition, false)
	if len(altIfStmt.Consequence.Statements) != 2 {
		t.Fatalf("expected len of consequence block of alt-if-statement to be 2. got=%d", len(altIfStmt.Consequence.Statements))
	}
	altConsequence, ok := altIfStmt.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("altIfStmt.Consequence.Statements[0] is not *ast.ExpressionStatement. got=%T", altIfStmt.Consequence.Statements[0])
	}
	testExpression(t, altConsequence.Expression, 2)
	altConsequence, ok = altIfStmt.Consequence.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("altIfStmt.Consequence.Statements[1] is not *ast.ExpressionStatement. got=%T", altIfStmt.Consequence.Statements[1])
	}
	testExpression(t, altConsequence.Expression, 3)
	altAltBlock, ok := altIfStmt.Alternative.(*ast.BlockStatement)
	if !ok {
		t.Fatalf("altIfStmt.Alternative is not *ast.BlockStatement. got=%T", altIfStmt.Alternative)
	}
	if len(altAltBlock.Statements) != 3 {
		t.Fatalf("expected len of altIf.Alternative block statements to be 3. got=%d", len(altAltBlock.Statements))
	}
	blockEntry, ok := altAltBlock.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("altAltBlock.Statements[0] is not *ast.ExpressionStatement. got=%T", altAltBlock.Statements[0])
	}
	testExpression(t, blockEntry.Expression, 4)
	blockEntry, ok = altAltBlock.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("altAltBlock.Statements[1] is not *ast.ExpressionStatement. got=%T", altAltBlock.Statements[1])
	}
	testExpression(t, blockEntry.Expression, 5)
	blockEntry, ok = altAltBlock.Statements[2].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("altAltBlock.Statements[2] is not *ast.ExpressionStatement. got=%T", altAltBlock.Statements[2])
	}
	testExpression(t, blockEntry.Expression, 6)
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		program := setupTestWithInput(t, tt.input)
		if len(program.Statements) != 1 {
			t.Fatalf("expected len of parsed program to be 1. got=%d", len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("statement is not of type *ast.ExpressionStatement. got=%T", stmt)
		}
		testBooleanLiteral(t, stmt.Expression, tt.want)
	}
}

func testBooleanLiteral(t *testing.T, boolExpression ast.Expression, value bool) bool {
	boolExp, ok := boolExpression.(*ast.Boolean)
	if !ok {
		t.Fatalf("expression is not of type boolean. got=%T", boolExpression)
		return false
	}
	if isEq, failMsg := testutils.Equal(value, boolExp.Value); !isEq {
		t.Errorf("boolean value is incorrect: %s", failMsg)
		return false
	}
	return true
}

func TestIdentifierExpression(t *testing.T) {
	program := setupTestWithInput(t, "myIdent")
	if len(program.Statements) != 1 {
		t.Fatalf("expected len of parsed program to be 1. got=%d", len(program.Statements))
	}
	exp, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("program.Statements[0] is not of type *ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	testIdentifier(t, exp.Expression, "myIdent")
}

func testIdentifier(t *testing.T, identifierExpression ast.Expression, value string) bool {
	ident, ok := identifierExpression.(*ast.Identifier)
	if !ok {
		t.Fatalf("passed identifier expression is not an *ast.Identifier. got type=%T", identifierExpression)
		return false
	}

	if isEq, errMsg := testutils.Equal(value, ident.Value); !isEq {
		t.Errorf("identifiers are not equal: %s", errMsg)
		return false
	}

	return true
}

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
		return false
	}

	if isEq, failMsg := testutils.Equal(value, integer.Value); !isEq {
		t.Errorf("testIntegerLiteral: value equality check failed: %s", failMsg)
		return false
	}

	return true
}

func testExpression(t *testing.T, expression ast.Expression, value any) bool {

	switch v := value.(type) {
	case int:
		return testIntegerLiteral(t, expression, v)
	case int64:
		return testIntegerLiteral(t, expression, int(v))
	case bool:
		return testBooleanLiteral(t, expression, v)
	}
	t.Errorf("type of exp not handled. got=%T", expression)
	return false
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
