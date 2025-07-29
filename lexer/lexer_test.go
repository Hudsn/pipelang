package lexer

import (
	"testing"

	"github.com/hudsn/pipelang/token"
	"github.com/hudsn/pipelang/utils/testutils"
)

func TestLexer(t *testing.T) {
	input := `
	1234
	5.4321;
	myident
	"my double quote"
	'my single quote' 
	"two strings" 'one line'
	math = 1 + 1 - 1 / 1 * 1
	1 = 1 < 2 > 3 <= 4 >= 3 != 0 == 0
	$src $dest $env $madeup
	[1, "a", ident]
	fn(x, y){return x + y}
	pipeline mypipe {
		|func1
		|func2(arg)
	}
	`

	l := New([]rune(input))
	tests := []struct {
		token token.Token
	}{
		{token.Token{
			Type:    token.INT,
			Literal: "1234",
			Line:    2,
			Col:     2,
		}},
		{token.Token{
			Type:    token.SEMICOLON,
			Literal: ";",
			Line:    2,
			Col:     6,
		}},
		{token.Token{
			Type:    token.FLOAT,
			Literal: "5.4321",
			Line:    3,
			Col:     2,
		}},
		{token.Token{
			Type:    token.SEMICOLON,
			Literal: ";",
			Line:    3,
			Col:     8,
		}},
		{token.Token{
			Type:    token.IDENT,
			Literal: "myident",
			Line:    4,
			Col:     2,
		}},
		{token.Token{
			Type:    token.SEMICOLON,
			Literal: ";",
			Line:    4,
			Col:     9,
		}},
		{token.Token{
			Type:    token.STRING,
			Literal: "my double quote",
			Line:    5,
			Col:     2,
		}},
		{token.Token{
			Type:    token.SEMICOLON,
			Literal: ";",
			Line:    5,
			Col:     19,
		}},
		{token.Token{
			Type:    token.STRING,
			Literal: "my single quote",
			Line:    6,
			Col:     2,
		}},
		{token.Token{
			Type:    token.SEMICOLON,
			Literal: ";",
			Line:    6,
			Col:     20,
		}},
		{token.Token{
			Type:    token.STRING,
			Literal: "two strings",
			Line:    7,
			Col:     2,
		}},
		{token.Token{
			Type:    token.STRING,
			Literal: "one line",
			Line:    7,
			Col:     16,
		}},
		{token.Token{
			Type:    token.SEMICOLON,
			Literal: ";",
			Line:    7,
			Col:     26,
		}},
		{token.Token{
			Type:    token.IDENT,
			Literal: "math",
			Line:    8,
			Col:     2,
		}},
		{token.Token{
			Type:    token.ASSIGN,
			Literal: "=",
			Line:    8,
			Col:     7,
		}},
		{token.Token{
			Type:    token.INT,
			Literal: "1",
			Line:    8,
			Col:     9,
		}},
		{token.Token{
			Type:    token.PLUS,
			Literal: "+",
			Line:    8,
			Col:     11,
		}},
		{token.Token{
			Type:    token.INT,
			Literal: "1",
			Line:    8,
			Col:     13,
		}},
		{token.Token{
			Type:    token.MINUS,
			Literal: "-",
			Line:    8,
			Col:     15,
		}},
		{token.Token{
			Type:    token.INT,
			Literal: "1",
			Line:    8,
			Col:     17,
		}},
		{token.Token{
			Type:    token.SLASH,
			Literal: "/",
			Line:    8,
			Col:     19,
		}},
		{token.Token{
			Type:    token.INT,
			Literal: "1",
			Line:    8,
			Col:     21,
		}},
		{token.Token{
			Type:    token.ASTERISK,
			Literal: "*",
			Line:    8,
			Col:     23,
		}},
		{token.Token{
			Type:    token.INT,
			Literal: "1",
			Line:    8,
			Col:     25,
		}},
		{token.Token{
			Type:    token.SEMICOLON,
			Literal: ";",
			Line:    8,
			Col:     26,
		}},
		{token.Token{
			Type:    token.INT,
			Literal: "1",
			Line:    9,
			Col:     2,
		}},
		{token.Token{
			Type:    token.ASSIGN,
			Literal: "=",
			Line:    9,
			Col:     4,
		}},
		{token.Token{
			Type:    token.INT,
			Literal: "1",
			Line:    9,
			Col:     6,
		}},
		{token.Token{
			Type:    token.LT,
			Literal: "<",
			Line:    9,
			Col:     8,
		}},
		{token.Token{
			Type:    token.INT,
			Literal: "2",
			Line:    9,
			Col:     10,
		}},
		{token.Token{
			Type:    token.GT,
			Literal: ">",
			Line:    9,
			Col:     12,
		}},
		{token.Token{
			Type:    token.INT,
			Literal: "3",
			Line:    9,
			Col:     14,
		}},
		{token.Token{
			Type:    token.LTEQ,
			Literal: "<=",
			Line:    9,
			Col:     16,
		}},
		{token.Token{
			Type:    token.INT,
			Literal: "4",
			Line:    9,
			Col:     19,
		}},
		{token.Token{
			Type:    token.GTEQ,
			Literal: ">=",
			Line:    9,
			Col:     21,
		}},
		{token.Token{
			Type:    token.INT,
			Literal: "3",
			Line:    9,
			Col:     24,
		}},
		{token.Token{
			Type:    token.NOT_EQ,
			Literal: "!=",
			Line:    9,
			Col:     26,
		}},
		{token.Token{
			Type:    token.INT,
			Literal: "0",
			Line:    9,
			Col:     29,
		}},
		{token.Token{
			Type:    token.EQ,
			Literal: "==",
			Line:    9,
			Col:     31,
		}},
		{token.Token{
			Type:    token.INT,
			Literal: "0",
			Line:    9,
			Col:     34,
		}},
		{token.Token{
			Type:    token.SEMICOLON,
			Literal: ";",
			Line:    9,
			Col:     35,
		}},
		{token.Token{
			Type:    token.SRC,
			Literal: "$src",
			Line:    10,
			Col:     2,
		}},
		{token.Token{
			Type:    token.DEST,
			Literal: "$dest",
			Line:    10,
			Col:     7,
		}},
		{token.Token{
			Type:    token.ENV,
			Literal: "$env",
			Line:    10,
			Col:     13,
		}},
		{token.Token{
			Type:    token.ILLEGAL,
			Literal: "$madeup",
			Line:    10,
			Col:     18,
		}},
		{token.Token{
			Type:    token.SEMICOLON,
			Literal: ";",
			Line:    10,
			Col:     25,
		}},
		{token.Token{
			Type:    token.LSQUARE,
			Literal: "[",
			Line:    11,
			Col:     2,
		}},
		{token.Token{
			Type:    token.INT,
			Literal: "1",
			Line:    11,
			Col:     3,
		}},
		{token.Token{
			Type:    token.COMMA,
			Literal: ",",
			Line:    11,
			Col:     4,
		}},
		{token.Token{
			Type:    token.STRING,
			Literal: "a",
			Line:    11,
			Col:     6,
		}},
		{token.Token{
			Type:    token.COMMA,
			Literal: ",",
			Line:    11,
			Col:     9,
		}},
		{token.Token{
			Type:    token.IDENT,
			Literal: "ident",
			Line:    11,
			Col:     11,
		}},
		{token.Token{
			Type:    token.RSQUARE,
			Literal: "]",
			Line:    11,
			Col:     16,
		}},
		{token.Token{
			Type:    token.SEMICOLON,
			Literal: ";",
			Line:    11,
			Col:     17,
		}},
		{token.Token{
			Type:    token.FUNCTION,
			Literal: "fn",
			Line:    12,
			Col:     2,
		}},
	}

	// fn(x, y){return x + y}
	// pipeline mypipe {
	// 	|func1
	// 	|func2(arg)
	// }
	// `

	for idx, tt := range tests {
		tok := l.NextToken()
		if isEq, failMsg := testutils.Equal(tt.token, tok); !isEq {
			t.Errorf("TestNumber %d: Wrong token: %s", idx, failMsg)
		}
	}
}

func TestAddSemicolonNewline(t *testing.T) {
	// input := "1234\r\n1.234\n54321\t\n"
	input := "1234\r\n1.234\n54321\t\n"
	l := New([]rune(input))

	tests := []struct {
		token token.Token
	}{
		{
			token.Token{
				Type:    token.INT,
				Literal: "1234",
				Line:    1,
				Col:     1,
			},
		},
		{
			token.Token{
				Type:    token.SEMICOLON,
				Literal: ";",
				Line:    1,
				Col:     5,
			},
		},
		{
			token.Token{
				Type:    token.FLOAT,
				Literal: "1.234",
				Line:    2,
				Col:     1,
			},
		},
		{
			token.Token{
				Type:    token.SEMICOLON,
				Literal: ";",
				Line:    2,
				Col:     6,
			},
		},
		{
			token.Token{
				Type:    token.INT,
				Literal: "54321",
				Line:    3,
				Col:     1,
			},
		},
		{
			token.Token{
				Type:    token.SEMICOLON,
				Literal: ";",
				Line:    3,
				Col:     7,
			},
		},
	}

	for idx, tt := range tests {
		tok := l.NextToken()
		if isEq, failMsg := testutils.Equal(tt.token, tok); !isEq {
			t.Errorf("TestNumber %d: Wrong token: %s", idx, failMsg)
		}
	}
}
