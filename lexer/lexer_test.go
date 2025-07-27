package lexer

import (
	"fmt"
	"testing"

	"github.com/hudsn/pipelang/token"
	"github.com/hudsn/pipelang/utils/testutils"
)

func TestAddSemicolonNewline(t *testing.T) {
	input := "1234\r\n1.234\n\n54321"
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
				Col:     6,
			},
		},
	}

	for idx, tt := range tests {
		tok := l.NextToken()
		fmt.Println(tt)
		if isEq, failMsg := testutils.Equal(tt.token, tok); !isEq {
			t.Errorf("TestNumber %d: Wrong token: %s", idx, failMsg)
		}
	}
}
