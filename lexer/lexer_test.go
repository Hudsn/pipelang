package lexer

import (
	"fmt"
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
	1 < 2 > 3 <= 4 >= 3 != 0 == 0
	$src $dest $env $var $madeup
	[1, "a", ident]
	{"key1" : 1, "key2" : 1.234}
	pipe mypipe {
		| otherPipe
		$dest.field1 = func1($src.field)
		if $src.field2 == "match" {
			$dest.field2 = "blah"
		} else if $src.field2 == null {
			emit()
		} else {
			drop()
		}
	}
`

	l := New([]rune(input))
	tests := []struct {
		value     string
		tokenType token.TokenType
		start     int
		end       int
	}{
		{
			value:     "1234",
			tokenType: token.INT,
			start:     2, //entire input starts wtih a newline and tab.
			end:       6,
		},
		{
			value:     ";",
			tokenType: token.SEMICOLON,
			start:     6, //entire input starts wtih a newline and tab.
			end:       7,
		},
		{
			tokenType: token.FLOAT,
			value:     "5.4321",
			start:     9,
			end:       15,
		},
		{
			tokenType: token.SEMICOLON,
			value:     ";",
			start:     15, //entire input starts wtih a newline and tab.
			end:       16,
		},
		{
			tokenType: token.IDENT,
			value:     "myident",
			start:     18,
			end:       25,
		},
		{
			tokenType: token.SEMICOLON,
			value:     ";",
			start:     25,
			end:       26,
		},
		{
			tokenType: token.STRING,
			value:     "my double quote",
			start:     28,
			end:       45,
		},
		{
			tokenType: token.SEMICOLON,
			value:     ";",
			start:     45,
			end:       46,
		},
		{
			tokenType: token.STRING,
			value:     "my single quote",
			start:     48,
			end:       65,
		},
		{ //note that this semicolon comes after a whitespace in the source
			tokenType: token.SEMICOLON,
			value:     ";",
			start:     66,
			end:       67,
		},
		{
			tokenType: token.STRING,
			value:     "two strings",
			start:     69,
			end:       82,
		},
		{
			tokenType: token.STRING,
			value:     "one line",
			start:     83,
			end:       93,
		},
		{
			tokenType: token.SEMICOLON,
			value:     ";",
			start:     93,
			end:       94,
		},
		{
			tokenType: token.IDENT,
			value:     "math",
			start:     96,
			end:       100,
		},
		{
			tokenType: token.ASSIGN,
			value:     "=",
			start:     101,
			end:       102,
		},
		{
			tokenType: token.INT,
			value:     "1",
			start:     103,
			end:       104,
		},
		{
			tokenType: token.PLUS,
			value:     "+",
			start:     105,
			end:       106,
		},
		{
			tokenType: token.INT,
			value:     "1",
			start:     107,
			end:       108,
		},
		{
			tokenType: token.MINUS,
			value:     "-",
			start:     109,
			end:       110,
		},
		{
			tokenType: token.INT,
			value:     "1",
			start:     111,
			end:       112,
		},
		{
			tokenType: token.SLASH,
			value:     "/",
			start:     113,
			end:       114,
		},
		{
			tokenType: token.INT,
			value:     "1",
			start:     115,
			end:       116,
		},
		{
			tokenType: token.ASTERISK,
			value:     "*",
			start:     117,
			end:       118,
		},
		{
			tokenType: token.INT,
			value:     "1",
			start:     119,
			end:       120,
		},
		{
			tokenType: token.SEMICOLON,
			value:     ";",
			start:     120,
			end:       121,
		},
		{
			tokenType: token.INT,
			value:     "1",
			start:     123,
			end:       124,
		},
		{
			tokenType: token.LT,
			value:     "<",
			start:     125,
			end:       126,
		},
		{
			tokenType: token.INT,
			value:     "2",
			start:     127,
			end:       128,
		},
		{
			tokenType: token.GT,
			value:     ">",
			start:     129,
			end:       130,
		},
		{
			tokenType: token.INT,
			value:     "3",
			start:     131,
			end:       132,
		},
		{
			tokenType: token.LTEQ,
			value:     "<=",
			start:     133,
			end:       135,
		},
		{
			tokenType: token.INT,
			value:     "4",
			start:     136,
			end:       137,
		},
		{
			tokenType: token.GTEQ,
			value:     ">=",
			start:     138,
			end:       140,
		},
		{
			tokenType: token.INT,
			value:     "3",
			start:     141,
			end:       142,
		},
		{
			tokenType: token.NOT_EQ,
			value:     "!=",
			start:     143,
			end:       145,
		},
		{
			tokenType: token.INT,
			value:     "0",
			start:     146,
			end:       147,
		},
		{
			tokenType: token.EQ,
			value:     "==",
			start:     148,
			end:       150,
		},
		{
			tokenType: token.INT,
			value:     "0",
			start:     151,
			end:       152,
		},
		{
			tokenType: token.SEMICOLON,
			value:     ";",
			start:     152,
			end:       153,
		},
		{
			tokenType: token.SRC,
			value:     "$src",
			start:     155,
			end:       159,
		},
		{
			tokenType: token.DEST,
			value:     "$dest",
			start:     160,
			end:       165,
		},
		{
			tokenType: token.ENV,
			value:     "$env",
			start:     166,
			end:       170,
		},
		{
			tokenType: token.VAR,
			value:     "$var",
			start:     171,
			end:       175,
		},
		{
			tokenType: token.ILLEGAL,
			value:     "$madeup",
			start:     176,
			end:       183,
		},
		{
			tokenType: token.SEMICOLON,
			value:     ";",
			start:     183,
			end:       184,
		},
		{
			tokenType: token.LSQUARE,
			value:     "[",
			start:     186,
			end:       187,
		},
		{
			tokenType: token.INT,
			value:     "1",
			start:     187,
			end:       188,
		},
		{
			tokenType: token.COMMA,
			value:     ",",
			start:     188,
			end:       189,
		},
		{
			tokenType: token.STRING,
			value:     "a",
			start:     190,
			end:       193,
		},
		{
			tokenType: token.COMMA,
			value:     ",",
			start:     193,
			end:       194,
		},
		{
			tokenType: token.IDENT,
			value:     "ident",
			start:     195,
			end:       200,
		},
		{
			tokenType: token.RSQUARE,
			value:     "]",
			start:     200,
			end:       201,
		},
		{
			tokenType: token.SEMICOLON,
			value:     ";",
			start:     201,
			end:       202,
		},
		{
			tokenType: token.LCURLY,
			value:     "{",
			start:     204,
			end:       205,
		},
		{
			tokenType: token.STRING,
			value:     "key1",
			start:     205,
			end:       211,
		},
		{
			tokenType: token.COLON,
			value:     ":",
			start:     212,
			end:       213,
		},
		{
			tokenType: token.INT,
			value:     "1",
			start:     214,
			end:       215,
		},
		{
			tokenType: token.COMMA,
			value:     ",",
			start:     215,
			end:       216,
		},
		{
			tokenType: token.STRING,
			value:     "key2",
			start:     217,
			end:       223,
		},
		{
			tokenType: token.COLON,
			value:     ":",
			start:     224,
			end:       225,
		},
		{
			tokenType: token.FLOAT,
			value:     "1.234",
			start:     226,
			end:       231,
		},
		{
			tokenType: token.RCURLY,
			value:     "}",
			start:     231,
			end:       232,
		},
		{
			tokenType: token.SEMICOLON,
			value:     ";",
			start:     232,
			end:       233,
		},
		{
			tokenType: token.PIPEDEF,
			value:     "pipe",
			start:     235,
			end:       239,
		},
		{
			tokenType: token.IDENT,
			value:     "mypipe",
			start:     240,
			end:       246,
		},
		{
			tokenType: token.LCURLY,
			value:     "{",
			start:     247,
			end:       248,
		},
		{
			tokenType: token.PIPECHAR,
			value:     "|",
			start:     251,
			end:       252,
		},
		{
			tokenType: token.IDENT,
			value:     "otherPipe",
			start:     253,
			end:       262,
		},
		{
			tokenType: token.SEMICOLON,
			value:     ";",
			start:     262,
			end:       263,
		},
		{
			tokenType: token.DEST,
			value:     "$dest",
			start:     266,
			end:       271,
		},
		{
			tokenType: token.DOT,
			value:     ".",
			start:     271,
			end:       272,
		},
		{
			tokenType: token.IDENT,
			value:     "field1",
			start:     272,
			end:       278,
		},
		{
			tokenType: token.ASSIGN,
			value:     "=",
			start:     279,
			end:       280,
		},
		{
			tokenType: token.IDENT,
			value:     "func1",
			start:     281,
			end:       286,
		},
		{
			tokenType: token.LPAREN,
			value:     "(",
			start:     286,
			end:       287,
		},
		{
			tokenType: token.SRC,
			value:     "$src",
			start:     287,
			end:       291,
		},
		{
			tokenType: token.DOT,
			value:     ".",
			start:     291,
			end:       292,
		},
		{
			tokenType: token.IDENT,
			value:     "field",
			start:     292,
			end:       297,
		},
		{
			tokenType: token.RPAREN,
			value:     ")",
			start:     297,
			end:       298,
		},
		{
			tokenType: token.SEMICOLON,
			value:     ";",
			start:     298,
			end:       299,
		},
		{
			tokenType: token.IF,
			value:     "if",
			start:     302,
			end:       304,
		},
		{
			tokenType: token.SRC,
			value:     "$src",
			start:     305,
			end:       309,
		},
		{
			tokenType: token.DOT,
			value:     ".",
			start:     309,
			end:       310,
		},
		{
			tokenType: token.IDENT,
			value:     "field2",
			start:     310,
			end:       316,
		},
		{
			tokenType: token.EQ,
			value:     "==",
			start:     317,
			end:       319,
		},
		{
			tokenType: token.STRING,
			value:     "match",
			start:     320,
			end:       327,
		},
		{
			tokenType: token.LCURLY,
			value:     "{",
			start:     328,
			end:       329,
		},
		{
			tokenType: token.DEST,
			value:     "$dest",
			start:     333,
			end:       338,
		},
		{
			tokenType: token.DOT,
			value:     ".",
			start:     338,
			end:       339,
		},
		{
			tokenType: token.IDENT,
			value:     "field2",
			start:     339,
			end:       345,
		},
		{
			tokenType: token.ASSIGN,
			value:     "=",
			start:     346,
			end:       347,
		},
		{
			tokenType: token.STRING,
			value:     "blah",
			start:     348,
			end:       354,
		},
		{
			tokenType: token.SEMICOLON,
			value:     ";",
			start:     354,
			end:       355,
		},
		{
			tokenType: token.RCURLY,
			value:     "}",
			start:     358,
			end:       359,
		},
		{
			tokenType: token.ELSE,
			value:     "else",
			start:     360,
			end:       364,
		},
		{
			tokenType: token.IF,
			value:     "if",
			start:     365,
			end:       367,
		},
		{
			tokenType: token.SRC,
			value:     "$src",
			start:     368,
			end:       372,
		},
		{
			tokenType: token.DOT,
			value:     ".",
			start:     372,
			end:       373,
		},
		{
			tokenType: token.IDENT,
			value:     "field2",
			start:     373,
			end:       379,
		},
		{
			tokenType: token.EQ,
			value:     "==",
			start:     380,
			end:       382,
		},
		{
			tokenType: token.NULL,
			value:     "null",
			start:     383,
			end:       387,
		},
		{
			tokenType: token.LCURLY,
			value:     "{",
			start:     388,
			end:       389,
		},
		{
			tokenType: token.IDENT,
			value:     "emit",
			start:     393,
			end:       397,
		},
		{
			tokenType: token.LPAREN,
			value:     "(",
			start:     397,
			end:       398,
		},
		{
			tokenType: token.RPAREN,
			value:     ")",
			start:     398,
			end:       399,
		},
		{
			tokenType: token.SEMICOLON,
			value:     ";",
			start:     399,
			end:       400,
		},
		{
			tokenType: token.RCURLY,
			value:     "}",
			start:     403,
			end:       404,
		},
		{
			tokenType: token.ELSE,
			value:     "else",
			start:     405,
			end:       409,
		},
		{
			tokenType: token.LCURLY,
			value:     "{",
			start:     410,
			end:       411,
		},
		{
			tokenType: token.IDENT,
			value:     "drop",
			start:     415,
			end:       419,
		},
		{
			tokenType: token.LPAREN,
			value:     "(",
			start:     419,
			end:       420,
		},
		{
			tokenType: token.RPAREN,
			value:     ")",
			start:     420,
			end:       421,
		},
		{
			tokenType: token.SEMICOLON,
			value:     ";",
			start:     421,
			end:       422,
		},
		{
			tokenType: token.RCURLY,
			value:     "}",
			start:     425,
			end:       426,
		},
		{
			tokenType: token.SEMICOLON,
			value:     ";",
			start:     426,
			end:       427,
		},
		{
			tokenType: token.RCURLY,
			value:     "}",
			start:     429,
			end:       430,
		},
		{
			tokenType: token.SEMICOLON,
			value:     ";",
			start:     430,
			end:       431,
		},
		{
			tokenType: token.EOF,
			value:     string(rune(0)),
			start:     432,
			end:       433,
		},
	}

	type checkRanges struct {
		want  string
		start int
		end   int
	}
	checkRangesList := []checkRanges{}

	for idx, tt := range tests {
		tok := l.NextToken()
		if isEq, failMsg := testutils.Equal(tt.tokenType, tok.Type); !isEq {
			t.Errorf("Test #%d: Wrong token type: %s", idx, failMsg)
		}
		if isEq, failMsg := testutils.Equal(tt.value, tok.Value); !isEq {
			t.Errorf("Test #%d: Wrong token value: %s", idx, failMsg)
		}

		start, end := tok.Position.GetPosition()
		if isEq, failMsg := testutils.Equal(tt.start, start); !isEq {
			t.Errorf("Test #%d: Wrong start index: %s", idx, failMsg)
		}
		if isEq, failMsg := testutils.Equal(tt.end, end); !isEq {
			t.Errorf("Test #%d: Wrong end index: %s", idx, failMsg)
		}

		if tok.Type == token.STRING {
			quoteChar := l.input[start]
			if quoteChar != l.input[end-1] {
				t.Errorf("Test #%d: mismatched quote char. want=%s. got=%s", idx, string(quoteChar), string(l.input[end-1]))
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
