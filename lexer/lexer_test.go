package lexer

import (
	"testing"
)

func TestLexer(t *testing.T) {
	input := `1234
	
	`

	l := New([]rune(input))

	l.NextToken()

}
