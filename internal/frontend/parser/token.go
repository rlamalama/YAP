package parser

type TokenKind int

// Order must remain the same as below
const (
	TokenDash TokenKind = iota
	TokenIdentifier
	TokenColon
	TokenScalar
	TokenIndent
	TokenDedent
	TokenNewline
	TokenEOF
)

func (tk TokenKind) String() string {
	// Order must remain the same as above
	names := [...]string{
		"Dash",
		"Identifer",
		"Colon",
		"Scalar",
		"Indent",
		"TokenDedent",
		"TokenNewline",
		"TokenEOF",
	}

	if tk < TokenDash || tk > TokenEOF {
		return "Unknown"
	}
	return names[tk]
}

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
