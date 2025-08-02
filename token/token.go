package token

type Token struct {
	Type     TokenType
	Value    string
	Position Position
}

func (t *Token) SetPosition(start int, end int) {
	t.Position = Position{
		start: start,
		end:   end,
	}
}

type Position struct {
	start int // inclusive: [start:end)
	end   int // exclusive: [start:end)
}

var NullPosition Position = Position{
	start: -1,
	end:   -1,
}

func (p *Position) GetPosition() (int, int) {
	return p.start, p.end
}

func (p *Position) SetPosition(start int, end int) {
	p.start = start
	p.end = end
}

type TokenType int

const (
	_ TokenType = iota

	//identifiers and literal types
	IDENT  // identifier; ex: myVar
	STRING // "mystring"
	INT
	FLOAT

	//operators
	ASSIGN      // "="
	PLUS        // "+"
	MINUS       // "-"
	ASTERISK    // "*"
	SLASH       // "/"
	EXCLAMATION // "!"
	PIPECHAR    // "|"

	//comparisons
	EQ     // "=="
	NOT_EQ //"!="
	GT     // >
	LT     // <
	GTEQ   //">="
	LTEQ   // "<="

	//delimiters
	DOT       // "."
	COMMA     // ","
	COLON     // ":"
	SEMICOLON // ";"
	LSQUARE   // "["
	RSQUARE   // "]"
	LCURLY    // "{"
	RCURLY    // "}"
	LPAREN    // "("
	RPAREN    // ")"

	//keywords
	PIPEDEF // "pipe"
	IF
	ELSE
	NULL

	//mem accessors
	ENV  // "$env"
	VAR  // "$var"
	SRC  // $src
	DEST // $dest

	// meta
	ILLEGAL // "ILLEGAL"
	EOF     // "EOF"
)

func LookupKeyword(ident string) TokenType {
	if keyword, ok := keywordTable[ident]; ok {
		return keyword
	}
	return IDENT
}

var keywordTable = map[string]TokenType{
	"pipe":  PIPEDEF,
	"$env":  ENV,
	"$var":  VAR,
	"$src":  SRC,
	"$dest": DEST,
	"if":    IF,
	"else":  ELSE,
	"null":  NULL,
}
