package parser_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/rlamalama/YAP/internal/frontend/parser"
	test_util "github.com/rlamalama/YAP/test/test-util"
	"github.com/stretchr/testify/assert"
)

const testFileDirPrefix = "../../.."

func TestLexEmptyFile(t *testing.T) {
	file := test_util.OpenTestFile(t, test_util.EmptyFileYAP, testFileDirPrefix)
	defer file.Close()

	lex := parser.NewLexer(file, test_util.EmptyFileYAP)
	tokens, err := lex.Lex()

	// Expecting 0 tokens in an empty file
	assert.Nil(t, err)
	assert.Equal(t, 0, len(tokens))
}

// - print: "hello world"
func TestLexOneLinePrintStatement(t *testing.T) {
	file := test_util.OpenTestFile(t, test_util.OneLinePrintYAP, testFileDirPrefix)
	defer file.Close()

	lex := parser.NewLexer(file, test_util.OneLinePrintYAP)
	tokens, err := lex.Lex()

	// Expecting 5 tokens in simple print file
	assert.Nil(t, err)
	assert.Equal(t, 5, len(tokens))
	assert.Equal(t, parser.TokenDash, tokens[0].Kind)
	assert.Equal(t, "-", tokens[0].Value)
	assert.Equal(t, 1, tokens[0].Col)
	assert.Equal(t, 1, tokens[0].Line)

	assert.Equal(t, parser.TokenIdentifier, tokens[1].Kind)
	assert.Equal(t, "print", tokens[1].Value)
	assert.Equal(t, 3, tokens[1].Col)
	assert.Equal(t, 1, tokens[1].Line)

	assert.Equal(t, parser.TokenColon, tokens[2].Kind)
	assert.Equal(t, ":", tokens[2].Value)
	assert.Equal(t, 8, tokens[2].Col)
	assert.Equal(t, 1, tokens[2].Line)

	assert.Equal(t, parser.TokenScalar, tokens[3].Kind)
	assert.Equal(t, "hello world", tokens[3].Value)
	assert.Equal(t, 10, tokens[3].Col)
	assert.Equal(t, 1, tokens[3].Line)

	assert.Equal(t, parser.TokenNewline, tokens[4].Kind)
	assert.Equal(t, "", tokens[4].Value)
	assert.Equal(t, 23, tokens[4].Col)
	assert.Equal(t, 1, tokens[4].Line)
}

// Multiline should have 5 tokens per line, with 5 lines = 25 tokens
func TestLexMultiLinePrintStatement(t *testing.T) {
	file := test_util.OpenTestFile(t, test_util.MultiLinePrintYAP, testFileDirPrefix)
	defer file.Close()

	lex := parser.NewLexer(file, test_util.MultiLinePrintYAP)
	tokens, err := lex.Lex()

	// Expecting 5 tokens in print file
	assert.Nil(t, err)
	assert.Equal(t, 25, len(tokens))
}

func TestLexBasicSetIfStatement(t *testing.T) {
	_ = `
- set:
    x: 5
- if: 
  condition: "x > 3"
  then:
    - print: "hello world"
`
	// Dash ID COLON NL (4)
	// INDENT (1)
	// ID COLON SCALAR NL (4)
	// DEDENT (1)
	// DASH ID COLON NL (4) good through here
	// INDENT (1)
	// ID COLON SCLAR  NL (4)
	// ID COLON NL (3)
	// INDENT (1)
	// DASH ID COLON SCALAR NL (5)
	// DEDENT DEDENT (2)

	file := test_util.OpenTestFile(t, test_util.BasicSetIfYAP, testFileDirPrefix)
	defer file.Close()

	lex := parser.NewLexer(file, test_util.BasicSetIfYAP)
	tokens, err := lex.Lex()
	for _, tok := range tokens {
		fmt.Println(tok.Kind.String())
	}

	// Expecting 5 tokens in print file

	assert.Nil(t, err)
	assert.Equal(t, 30, len(tokens))
}

func TestLexNoTabs(t *testing.T) {
	file := test_util.OpenTestFile(t, test_util.NoTabCharYAP, testFileDirPrefix)
	defer file.Close()

	lex := parser.NewLexer(file, test_util.NoTabCharYAP)
	_, err := lex.Lex()
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "tab"))

}
