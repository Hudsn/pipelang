package token

type TokenType string

type Token struct {
	Literal string
	Type    TokenType
	Start   Position
	End     Position
}

func (t *Token) SetPositions(line int, col int) {
	t.Start.Set(line, col)
	t.End.Set(line, col)
}

func (t *Token) SetStart(line int, col int) {
	t.Start.Set(line, col)
}

func (t *Token) SetEnd(line int, col int) {
	t.End.Set(line, col)
}

type Position struct {
	Line int
	Col  int
}

func (p *Position) Set(line int, col int) {
	p.Line = line
	p.Col = col
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
