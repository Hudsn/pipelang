package parser

import (
	"testing"

	"github.com/hudsn/pipelang/ast"
	"github.com/hudsn/pipelang/lexer"
	"github.com/hudsn/pipelang/utils/testutils"
)

func TestArrowFunctionExpression(t *testing.T) {

}

func TestFunctionCallExpression(t *testing.T) {

}

func TestPrefixExpression(t *testing.T) {

}

func TestInfixExpression(t *testing.T) {

}

func TestPipedefStatement(t *testing.T) {

}

func TestPipeCallStatement(t *testing.T) {

}

func TestDotAccessExpression(t *testing.T) {

}

// TODO LINE -- move items below line after finished

func TestAssignStatement(t *testing.T) {
	input := "a = 2"
	program := setupTestWithInput(t, input)
	if len(program.Statements) != 1 {
		t.Fatalf("expected len of parsed program to be 1. got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.AssignStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.AssignStatement. got=%T", program.Statements[0])
	}
	testIdentifier(t, stmt.Name, "a")

}

func TestFloatLiteral(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{".1234", 0.1234},
		{"5.4321", 5.4321},
	}

	for _, tt := range tests {
		program := setupTestWithInput(t, tt.input)
		if len(program.Statements) != 1 {
			t.Fatalf("expected len of parsed program to be 1. got=%d", len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
		}
		testFloatLiteral(t, stmt.Expression, tt.want)
	}
}

func testFloatLiteral(t *testing.T, floatExpression ast.Expression, want float64) bool {
	fl, ok := floatExpression.(*ast.FloatLiteral)
	if !ok {
		t.Fatalf("expression is not *ast.FloatLiteral. got=%T", floatExpression)
		return false
	}
	if isEQ, failMsg := testutils.Equal(want, fl.Value); !isEQ {
		t.Errorf("wrong float value: %s", failMsg)
		return false
	}
	return true
}

func TestStringLiteral(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{`"mystring"`, "mystring"},
		{"'mystring'", "mystring"},
	}

	for _, tt := range tests {
		program := setupTestWithInput(t, tt.input)
		if len(program.Statements) != 1 {
			t.Fatalf("expected len of parsed program to be 1. got=%d", len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
		}
		testStringLiteral(t, stmt.Expression, tt.want)
	}
}

func testStringLiteral(t *testing.T, stringExpression ast.Expression, want string) bool {
	str, ok := stringExpression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.StringLiteral. got=%T", stringExpression)
		return false
	}
	if isEq, failMsg := testutils.Equal(want, str.Value); !isEq {
		t.Errorf("wrong string value: %s", failMsg)
		return false
	}
	return true
}

func TestIfStatement(t *testing.T) {
	input := `if true {
				1
			}`

	program := setupTestWithInput(t, input)
	if len(program.Statements) != 1 {
		t.Fatalf("expected len of parsed program to be 1. got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.IfStatement. got=%T", program.Statements[0])
	}
	ifExp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.IfExpression. got=%T", stmt.Expression)
	}
	testBooleanLiteral(t, ifExp.Condition, true)
	if len(ifExp.Consequence.Statements) != 1 {
		t.Fatalf("expected len of consequence block of if-statement to be 1. got=%d", len(ifExp.Consequence.Statements))
	}
	consequence, ok := ifExp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ifStmt.Consequence.Statements[0] is not *ast.ExpressionStatement. got=%T", ifExp.Consequence.Statements[0])
	}
	testLiteralExpression(t, consequence.Expression, 1)

	if isEq, failMsg := testutils.Equal(nil, ifExp.Alternative); !isEq {
		t.Errorf("wanted empty/nil value for if-statement's alternative: %s", failMsg)
	}
}

func TestIfElseStatement(t *testing.T) {

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
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.IfStatement. got=%T", program.Statements[0])
	}
	ifExp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.IfExpression. got=%T", stmt.Expression)
	}

	testBooleanLiteral(t, ifExp.Condition, true)

	if len(ifExp.Consequence.Statements) != 1 {
		t.Fatalf("expected len of consequence block of if-statement to be 1. got=%d", len(ifExp.Consequence.Statements))
	}
	consequence, ok := ifExp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ifStmt.Consequence.Statements[0] is not *ast.ExpressionStatement. got=%T", ifExp.Consequence.Statements[0])
	}
	testLiteralExpression(t, consequence.Expression, 1)
	altIfStmt, ok := ifExp.Alternative.(*ast.IfExpression)
	if !ok {
		t.Fatalf("ifStmt.Alternative is not *ast.IfStatement. got=%T", ifExp.Alternative)
	}
	testBooleanLiteral(t, altIfStmt.Condition, false)
	if len(altIfStmt.Consequence.Statements) != 2 {
		t.Fatalf("expected len of consequence block of alt-if-statement to be 2. got=%d", len(altIfStmt.Consequence.Statements))
	}
	altConsequence, ok := altIfStmt.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("altIfStmt.Consequence.Statements[0] is not *ast.ExpressionStatement. got=%T", altIfStmt.Consequence.Statements[0])
	}
	testLiteralExpression(t, altConsequence.Expression, 2)
	altConsequence, ok = altIfStmt.Consequence.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("altIfStmt.Consequence.Statements[1] is not *ast.ExpressionStatement. got=%T", altIfStmt.Consequence.Statements[1])
	}
	testLiteralExpression(t, altConsequence.Expression, 3)
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
	testLiteralExpression(t, blockEntry.Expression, 4)
	blockEntry, ok = altAltBlock.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("altAltBlock.Statements[1] is not *ast.ExpressionStatement. got=%T", altAltBlock.Statements[1])
	}
	testLiteralExpression(t, blockEntry.Expression, 5)
	blockEntry, ok = altAltBlock.Statements[2].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("altAltBlock.Statements[2] is not *ast.ExpressionStatement. got=%T", altAltBlock.Statements[2])
	}
	testLiteralExpression(t, blockEntry.Expression, 6)
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

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"1 + (2 + 3) * 4",
			"(1 + ((2 + 3) * 4))",
		},
		{
			"a = true || (1 == 2) && false",
			"a = ((true || (1 == 2)) && false)",
		},
		{
			"!a <= -b != c < d > e >= f",
			"(((!a) <= (-b)) != (((c < d) > e) >= f))",
		},
		{
			"a = b == c && d",
			"a = ((b == c) && d)",
		},
	}

	for _, tt := range tests {
		program := setupTestWithInput(t, tt.input)

		got := program.String()
		if isEq, failMsg := testutils.Equal(tt.expected, got); !isEq {
			t.Errorf("wrong precedence result: %s", failMsg)
		}
	}
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

func testInfixExpression(t *testing.T, expr ast.Expression, left any, operator string, right any) bool {
	infix, ok := expr.(*ast.InfixExpression)
	if !ok {
		t.Errorf("expression is not *ast.InfixExpression")
		return false
	}

	if !testLiteralExpression(t, infix.Left, left) {
		return false
	}

	if infix.Operator != operator {
		t.Errorf("expression.Operator is not %s. got=%q", operator, infix.Operator)
	}

	if !testLiteralExpression(t, infix.Right, right) {
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, expression ast.Expression, value any) bool {

	switch v := value.(type) {
	case int:
		return testIntegerLiteral(t, expression, v)
	case int64:
		return testIntegerLiteral(t, expression, int(v))
	case float64:
		return testFloatLiteral(t, expression, v)
	case bool:
		return testBooleanLiteral(t, expression, v)
	case string:
		return testStringLiteral(t, expression, v)
	default:
		t.Errorf("type of exp not handled. got=%T", expression)
		return false
	}
}

func setupTestWithInput(t *testing.T, input string) *ast.Program {
	lexer := lexer.New([]rune(input))
	parser := New(lexer)
	program, err := parser.ParseProgram()
	if err != nil {
		t.Fatalf("setupTestWithInput: %s", err.Error())
	}
	if err := parser.CheckParserErrors(); err != nil {
		t.Fatal(err)
	}
	return program
}
