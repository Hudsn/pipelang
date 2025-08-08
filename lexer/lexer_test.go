package lexer

import (
	"fmt"
	"testing"

	"github.com/hudsn/pipelang/token"
	"github.com/hudsn/pipelang/utils/testutils"
)

func TestLogicAnd(t *testing.T) {
	input := "a && b"
	cases := []testCase{
		{
			value:     "a",
			tokenType: token.IDENT,
			start:     0,
			end:       1,
		},
		{
			value:     "&&",
			tokenType: token.LOGIC_AND,
			start:     2,
			end:       4,
		},
		{
			value:     "b",
			tokenType: token.IDENT,
			start:     5,
			end:       6,
		},
	}
	checkTestCase(t, input, cases)
}
func TestLogicOr(t *testing.T) {
	input := "a || b"
	cases := []testCase{
		{
			value:     "a",
			tokenType: token.IDENT,
			start:     0,
			end:       1,
		},
		{
			value:     "||",
			tokenType: token.LOGIC_OR,
			start:     2,
			end:       4,
		},
		{
			value:     "b",
			tokenType: token.IDENT,
			start:     5,
			end:       6,
		},
	}
	checkTestCase(t, input, cases)
}

func TestLexArrow(t *testing.T) {
	input := "abc ~> 123"

	cases := []testCase{
		{
			value:     "abc",
			tokenType: token.IDENT,
			start:     0,
			end:       3,
		},
		{
			value:     "~>",
			tokenType: token.ARROW,
			start:     4,
			end:       6,
		},
		{
			value:     "123",
			tokenType: token.INT,
			start:     7,
			end:       10,
		},
	}
	checkTestCase(t, input, cases)
}

func TestLexInt(t *testing.T) {
	input := "1234"

	cases := []testCase{
		{
			value:     "1234",
			tokenType: token.INT,
			start:     0,
			end:       4,
		},
	}
	checkTestCase(t, input, cases)
}

func TestLexFloat(t *testing.T) {
	input := "1.234"
	cases := []testCase{
		{
			value:     "1.234",
			tokenType: token.FLOAT,
			start:     0,
			end:       5,
		},
	}
	checkTestCase(t, input, cases)
}

func TestLexIdent(t *testing.T) {
	input := "myidentifier"
	cases := []testCase{
		{
			value:     "myidentifier",
			tokenType: token.IDENT,
			start:     0,
			end:       12,
		},
	}
	checkTestCase(t, input, cases)
}

func TestLexDoubleQuote(t *testing.T) {
	input := `"quoted text"`
	cases := []testCase{
		{
			value:     "quoted text",
			tokenType: token.STRING,
			start:     0,
			end:       13,
		},
	}
	checkTestCase(t, input, cases)
}

func TestLexSingleQuote(t *testing.T) {
	input := `'quoted text'`
	cases := []testCase{
		{
			value:     "quoted text",
			tokenType: token.STRING,
			start:     0,
			end:       13,
		},
	}
	checkTestCase(t, input, cases)
}

func TestLexMathOps(t *testing.T) {
	input := "1 + 2 - 3 * 4 / 5"
	cases := []testCase{
		{
			value:     "1",
			tokenType: token.INT,
			start:     0,
			end:       1,
		},
		{
			value:     "+",
			tokenType: token.PLUS,
			start:     2,
			end:       3,
		},
		{
			value:     "2",
			tokenType: token.INT,
			start:     4,
			end:       5,
		},
		{
			value:     "-",
			tokenType: token.MINUS,
			start:     6,
			end:       7,
		},
		{
			value:     "3",
			tokenType: token.INT,
			start:     8,
			end:       9,
		},
		{
			value:     "*",
			tokenType: token.ASTERISK,
			start:     10,
			end:       11,
		},
		{
			value:     "4",
			tokenType: token.INT,
			start:     12,
			end:       13,
		},
		{
			value:     "/",
			tokenType: token.SLASH,
			start:     14,
			end:       15,
		},
		{
			value:     "5",
			tokenType: token.INT,
			start:     16,
			end:       17,
		},
	}
	checkTestCase(t, input, cases)
}

func TestLexInequality(t *testing.T) {
	input := "< <= > >= != =="
	cases := []testCase{
		{
			value:     "<",
			tokenType: token.LT,
			start:     0,
			end:       1,
		},
		{
			value:     "<=",
			tokenType: token.LTEQ,
			start:     2,
			end:       4,
		},
		{
			value:     ">",
			tokenType: token.GT,
			start:     5,
			end:       6,
		},
		{
			value:     ">=",
			tokenType: token.GTEQ,
			start:     7,
			end:       9,
		},
		{
			value:     "!=",
			tokenType: token.NOT_EQ,
			start:     10,
			end:       12,
		},
		{
			value:     "==",
			tokenType: token.EQ,
			start:     13,
			end:       15,
		},
	}
	checkTestCase(t, input, cases)
}

func TestLexSpecialKeywords(t *testing.T) {
	input := `
	$src
	$dest $env $var $madeup
	`
	cases := []testCase{
		{
			value:     "$src",
			tokenType: token.SRC,
			start:     2,
			end:       6,
		},
		{
			value:     ";",
			tokenType: token.SEMICOLON,
			start:     6,
			end:       7,
		},
		{
			value:     "$dest",
			tokenType: token.DEST,
			start:     9,
			end:       14,
		},
		{
			value:     "$env",
			tokenType: token.ENV,
			start:     15,
			end:       19,
		},
		{
			value:     "$var",
			tokenType: token.VAR,
			start:     20,
			end:       24,
		},
		{
			value:     "$madeup",
			tokenType: token.ILLEGAL,
			start:     25,
			end:       32,
		},
		{
			value:     ";",
			tokenType: token.SEMICOLON,
			start:     32,
			end:       33,
		},
		{
			value:     string(rune(0)),
			tokenType: token.EOF,
			start:     35,
			end:       36,
		},
	}
	checkTestCase(t, input, cases)
}

func TestLexDelimiters(t *testing.T) {
	input := "|()[]{}.,:;"
	cases := []testCase{
		{
			value:     "|",
			tokenType: token.PIPECHAR,
			start:     0,
			end:       1,
		},
		{
			value:     "(",
			tokenType: token.LPAREN,
			start:     1,
			end:       2,
		},
		{
			value:     ")",
			tokenType: token.RPAREN,
			start:     2,
			end:       3,
		},
		{
			value:     "[",
			tokenType: token.LSQUARE,
			start:     3,
			end:       4,
		},
		{
			value:     "]",
			tokenType: token.RSQUARE,
			start:     4,
			end:       5,
		},
		{
			value:     "{",
			tokenType: token.LCURLY,
			start:     5,
			end:       6,
		},
		{
			value:     "}",
			tokenType: token.RCURLY,
			start:     6,
			end:       7,
		},
		{
			value:     ".",
			tokenType: token.DOT,
			start:     7,
			end:       8,
		},
		{
			value:     ",",
			tokenType: token.COMMA,
			start:     8,
			end:       9,
		},
		{
			value:     ":",
			tokenType: token.COLON,
			start:     9,
			end:       10,
		},
		{
			value:     ";",
			tokenType: token.SEMICOLON,
			start:     10,
			end:       11,
		},
	}

	checkTestCase(t, input, cases)
}
func TestLexKeywords(t *testing.T) {
	input := "true false pipe if else null"
	cases := []testCase{
		{
			value:     "true",
			tokenType: token.TRUE,
			start:     0,
			end:       4,
		},
		{
			value:     "false",
			tokenType: token.FALSE,
			start:     5,
			end:       10,
		},
		{
			value:     "pipe",
			tokenType: token.PIPEDEF,
			start:     11,
			end:       15,
		},
		{
			value:     "if",
			tokenType: token.IF,
			start:     16,
			end:       18,
		},
		{
			value:     "else",
			tokenType: token.ELSE,
			start:     19,
			end:       23,
		},
		{
			value:     "null",
			tokenType: token.NULL,
			start:     24,
			end:       28,
		},
	}

	checkTestCase(t, input, cases)
}

type testCase struct {
	value     string
	tokenType token.TokenType
	start     int
	end       int
}
type checkRanges struct {
	want  string
	start int
	end   int
}

func checkTestCase(t *testing.T, input string, tests []testCase) {
	l := New([]rune(input))
	checkRangesList := []checkRanges{}

	for idx, tt := range tests {
		tok := l.NextToken()
		if isEq, failMsg := testutils.Equal(tt.tokenType, tok.Type); !isEq {
			t.Errorf("Test case #%d: Wrong token type: %s", idx, failMsg)
		}
		if isEq, failMsg := testutils.Equal(tt.value, tok.Value); !isEq {
			t.Errorf("Test case #%d: Wrong token value: %s", idx, failMsg)
		}

		start, end := tok.Position.GetPosition()
		if isEq, failMsg := testutils.Equal(tt.start, start); !isEq {
			t.Errorf("Test case #%d: Wrong start index: %s", idx, failMsg)
		}
		if isEq, failMsg := testutils.Equal(tt.end, end); !isEq {
			t.Errorf("Test case #%d: Wrong end index: %s", idx, failMsg)
		}

		if tok.Type == token.STRING {
			quoteChar := l.input[start]
			if quoteChar != l.input[end-1] {
				t.Errorf("Test case #%d: mismatched quote char. want=%s. got=%s", idx, string(quoteChar), string(l.input[end-1]))
			}
			toAdd := checkRanges{
				want:  fmt.Sprintf(`%s%s%s`, string(quoteChar), tok.Value, string(quoteChar)),
				start: start,
				end:   end,
			}
			checkRangesList = append(checkRangesList, toAdd)
		}
		if tok.Type == token.IDENT {
			toAdd := checkRanges{
				want:  tok.Value,
				start: start,
				end:   end,
			}
			checkRangesList = append(checkRangesList, toAdd)
		}
	}

	for _, tt := range checkRangesList {
		got := string(l.input[tt.start:tt.end])
		if isEq, failMsg := testutils.Equal(tt.want, got); !isEq {
			t.Errorf("Incorrect range produced unwanted string: %s", failMsg)
		}
	}
}
