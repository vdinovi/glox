package lox

//go:generate stringer -type TokenType
type TokenType int

const (
	None TokenType = iota
	LeftParen
	RightParen
	LeftBrace
	RightBrace
	Comma
	Dot
	Minus
	Plus
	Semicolon
	Slash
	Star
	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual
	Identifier
	String
	Number
	And
	Class
	Else
	False
	Fun
	For
	If
	Nil
	Or
	Print
	Return
	Super
	This
	True
	Var
	While
	Comment
	EOF
)

var tokenTypeMap = map[string]TokenType{
	"(":      LeftParen,
	")":      RightParen,
	"{":      LeftBrace,
	"}":      RightBrace,
	",":      Comma,
	".":      Dot,
	"-":      Minus,
	"+":      Plus,
	";":      Semicolon,
	"*":      Star,
	"!":      Bang,
	"!=":     BangEqual,
	"=":      Equal,
	"==":     EqualEqual,
	"<":      Less,
	"<=":     LessEqual,
	">":      Greater,
	">=":     GreaterEqual,
	"/":      Slash,
	"and":    And,
	"class":  Class,
	"else":   Else,
	"false":  False,
	"fun":    Fun,
	"for":    For,
	"if":     If,
	"nil":    Nil,
	"or":     Or,
	"print":  Print,
	"return": Return,
	"super":  Super,
	"this":   This,
	"true":   True,
	"var":    Var,
	"while":  While,
}

func TokenTypeFor(lexem string) TokenType {
	if tt, ok := tokenTypeMap[lexem]; ok {
		return tt
	} else {
		return None
	}
}

type Token struct {
	Type   TokenType
	Lexem  string
	Line   int
	Column int
}

var EofToken = Token{Type: EOF}

func (t Token) Eof() bool {
	return t == EofToken
}
