package lexer

type Keyword string

const (
	KeywordPrint = "print"
	KeywordSet   = "set"
)

var Keywords = []Keyword{
	KeywordPrint,
	KeywordSet,
}

func IsKeyword(s string) bool {
	for _, key := range Keywords {
		if string(key) == s {
			return true
		}
	}
	return false
}
