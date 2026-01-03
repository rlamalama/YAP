package parser

import (
	"os"
)

func ParseFile(pathname string) (string, error) {
	file, err := os.Open(pathname)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// TODO
	// lexer := NewLexer(file)
	//
	// lexer.Lex()
	return "", nil
}
