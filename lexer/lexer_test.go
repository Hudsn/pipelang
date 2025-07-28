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
	}

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
