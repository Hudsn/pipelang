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
	ASSIGN      TokenType = "="
	PLUS        TokenType = "+"
	MINUS       TokenType = "-"
	EXCLAMATION TokenType = "!"
	ASTERISK    TokenType = "*"
	SLASH       TokenType = "/"
	PIPE        TokenType = "|"

	// comparisons
	EQ     TokenType = "=="
	NOT_EQ TokenType = "!="
	GT     TokenType = ">"
	LT     TokenType = "<"
	GTEQ   TokenType = ">="
	LTEQ   TokenType = "<="

	// delimiters
	DOT       TokenType = "."
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"
	LSQUARE   TokenType = "["
	RSQUARE   TokenType = "]"
	LCURLY    TokenType = "{"
	RCURLY    TokenType = "}"
	LPAREN    TokenType = "("
	RPAREN    TokenType = ")"

	//keywords
	FUNCTION TokenType = "FUNCTION"
	PIPELINE TokenType = "PIPELINE"
	RETURN   TokenType = "RETURN"

	// global var accessors
	ENV  TokenType = "ENV"
	SRC  TokenType = "SRC"
	DEST TokenType = "DEST"

	//meta
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"
)

func LookupKeyword(ident string) TokenType {
	if keyword, ok := keywordTable[ident]; ok {
		return keyword
	}
	return IDENT
}

var keywordTable = map[string]TokenType{
	"fn":      FUNCTION,
	"pipeine": PIPELINE,
	"return":  RETURN,

	"$env":  ENV,
	"$src":  SRC,
	"$dest": DEST,
}
