package parser

type TokenKind int

const (
	TokenDash TokenKind = iota
	TokenIdentifier
	TokenColon
	TokenScalar
	TokenIndent
	TokenDedent
	TokenEOL
	TokenEOF
)

type Token struct {
	Kind  TokenKind
	Value string
	Line  int
	Col   int
}

func NewToken(kind TokenKind, value string, line, col int) *Token {
	return &Token{
		Kind:  kind,
		Value: value,
		Line:  line,
		Col:   col,
	}
}
