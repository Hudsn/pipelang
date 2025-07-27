package token

type TokenType string

type Token struct {
	Literal string
	Type    TokenType
	Line    int
	Col     int
}

func (t *Token) SetPosition(line int, col int) {
	t.Line = line
	t.Col = col
}

const (
	//identifiers and literals
	IDENT  TokenType = "IDENT"
	STRING TokenType = "STRING"
	INT    TokenType = "INT"
	FLOAT  TokenType = "FLOAT"

	// operators
	ASSIGN TokenType = "="

	// comparisons
	EQ TokenType = "=="

	// delimiters
	DOT       TokenType = "."
	SEMICOLON TokenType = ";"

	ILLEGAL TokenType = "ILLEGAL"
)
