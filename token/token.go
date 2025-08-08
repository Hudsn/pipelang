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

// [start:end)
type Position struct {
	start int
	end   int
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

func (t TokenType) HumanString() string {
	if ret, found := stringTable[t]; found {
		return ret
	}
	return "ILLEGAL"
}

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
	ARROW       // ~>
	LOGIC_OR    // ||
	LOGIC_AND   // &&

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
	TRUE
	FALSE

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
	"true":  TRUE,
	"false": FALSE,
}

var stringTable = map[TokenType]string{
	IDENT:       "identifier",
	STRING:      "string",
	INT:         "integer",
	FLOAT:       "float",
	ARROW:       `arrow ("~>")`,
	ASSIGN:      `assign ("=")`,
	PLUS:        `plus ("+")`,
	MINUS:       `minus ("-")`,
	ASTERISK:    `asterisk ("*")`,
	SLASH:       `slash ("/")`,
	EXCLAMATION: `exclamation ("!")`,
	PIPECHAR:    `pipechar ("|")`,
	EQ:          `equals ("==")`,
	NOT_EQ:      `not equals ("!=")`,
	GT:          `greater than (">")`,
	LT:          `less than ("<")`,
	GTEQ:        `greater than or equal to (">=")`,
	LTEQ:        `less than or equal to ("<=")`,
	DOT:         `dot (".")`,
	COMMA:       `comma (",")`,
	COLON:       `colon (":")`,
	SEMICOLON:   `semicolon (";")`,
	LSQUARE:     `left square bracket ("[")`,
	RSQUARE:     `right square bracket ("]")`,
	LCURLY:      `left curly bracket ("{")`,
	RCURLY:      `right curly bracket ("}")`,
	LPAREN:      `left parenth ("(")`,
	RPAREN:      `right parenth (")")`,
	PIPEDEF:     `pipe definition ("pipe")`,
	IF:          `if statement ("if")`,
	ELSE:        `else statement ("else")`,
	TRUE:        "true",
	FALSE:       "false",
	NULL:        "null token",
	ENV:         "$env",
	VAR:         "$var",
	SRC:         "$src",
	DEST:        "$dest",
	ILLEGAL:     "illegal token",
	EOF:         "end of file token",
}
