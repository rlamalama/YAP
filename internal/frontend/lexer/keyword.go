package lexer

type Keyword string

const (
	KeywordPrint = "print"
	KeywordSet   = "set"
	KeywordTrue  = "True"
	KeywordFalse = "False"
	KeywordIf    = "if"
	KeywordThen  = "then"
	KeywordElse  = "else"
)

var Keywords = []Keyword{
	KeywordPrint,
	KeywordSet,
	KeywordTrue,
	KeywordFalse,
	KeywordIf,
	KeywordThen,
	KeywordElse,
}

func IsKeyword(s string) bool {
	for _, key := range Keywords {
		if string(key) == s {
			return true
		}
	}
	return false
}
