package lexer_test

import (
	"strings"
	"testing"

	"github.com/rlamalama/YAP/internal/frontend/lexer"
	test_util "github.com/rlamalama/YAP/test/test-util"
	"github.com/stretchr/testify/assert"
)

const testFileDirPrefix = "../../.."

func TestLexEmptyFile(t *testing.T) {
	file := test_util.OpenTestFile(t, test_util.EmptyFileYAP, testFileDirPrefix)
	defer file.Close()

	lex := lexer.NewLexer(file, test_util.EmptyFileYAP)
	tokens, err := lex.Lex()

	// Expecting 0 tokens in an empty file
	assert.Nil(t, err)
	assert.Equal(t, 0, len(tokens))
}

// - print: "hello world"
func TestLexOneLinePrintStatement(t *testing.T) {
	file := test_util.OpenTestFile(t, test_util.OneLinePrintYAP, testFileDirPrefix)
	defer file.Close()

	lex := lexer.NewLexer(file, test_util.OneLinePrintYAP)
	tokens, err := lex.Lex()

	// Expecting 5 tokens in simple print file
	assert.Nil(t, err)
	assert.Equal(t, 5, len(tokens))
	assert.Equal(t, lexer.TokenDash, tokens[0].Kind)
	assert.Equal(t, "-", tokens[0].Value)
	assert.Equal(t, 1, tokens[0].Col)
	assert.Equal(t, 1, tokens[0].Line)

	assert.Equal(t, lexer.TokenKeyword, tokens[1].Kind)
	assert.Equal(t, lexer.KeywordPrint, tokens[1].Value)
	assert.Equal(t, 3, tokens[1].Col)
	assert.Equal(t, 1, tokens[1].Line)

	assert.Equal(t, lexer.TokenColon, tokens[2].Kind)
	assert.Equal(t, ":", tokens[2].Value)
	assert.Equal(t, 8, tokens[2].Col)
	assert.Equal(t, 1, tokens[2].Line)

	assert.Equal(t, lexer.TokenString, tokens[3].Kind)
	assert.Equal(t, "hello world", tokens[3].Value)
	assert.Equal(t, 10, tokens[3].Col)
	assert.Equal(t, 1, tokens[3].Line)

	assert.Equal(t, lexer.TokenNewline, tokens[4].Kind)
	assert.Equal(t, "", tokens[4].Value)
	assert.Equal(t, 23, tokens[4].Col)
	assert.Equal(t, 1, tokens[4].Line)
}

// Multiline should have 5 tokens per line, with 5 lines = 25 tokens
func TestLexMultiLinePrintStatement(t *testing.T) {
	file := test_util.OpenTestFile(t, test_util.MultiLinePrintYAP, testFileDirPrefix)
	defer file.Close()

	lex := lexer.NewLexer(file, test_util.MultiLinePrintYAP)
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
	// Dash ID COLON SCALAR NL (4)
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

	lex := lexer.NewLexer(file, test_util.BasicSetIfYAP)
	tokens, err := lex.Lex()
	assert.Nil(t, err)
	assert.Equal(t, 31, len(tokens))
}

func TestLexNoTabs(t *testing.T) {
	file := test_util.OpenTestFile(t, test_util.NoTabCharYAP, testFileDirPrefix)
	defer file.Close()

	lex := lexer.NewLexer(file, test_util.NoTabCharYAP)
	_, err := lex.Lex()
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "tab"))
}

func TestLexSetPrint(t *testing.T) {
	_ = `
- set:
  - x: "5"
  - y: 10
- print: x 
- print: y 
`
	expectedTok := []lexer.Token{
		{Kind: lexer.TokenDash},
		{Kind: lexer.TokenKeyword},
		{Kind: lexer.TokenColon},
		{Kind: lexer.TokenNewline},
		{Kind: lexer.TokenIndent},
		{Kind: lexer.TokenDash},
		{Kind: lexer.TokenIdentifier},
		{Kind: lexer.TokenColon},
		{Kind: lexer.TokenString},
		{Kind: lexer.TokenNewline},
		{Kind: lexer.TokenDash},
		{Kind: lexer.TokenIdentifier},
		{Kind: lexer.TokenColon},
		{Kind: lexer.TokenNumerical},
		{Kind: lexer.TokenNewline},
		{Kind: lexer.TokenDedent},
		{Kind: lexer.TokenDash},
		{Kind: lexer.TokenKeyword},
		{Kind: lexer.TokenColon},
		{Kind: lexer.TokenIdentifier},
		{Kind: lexer.TokenNewline},
		{Kind: lexer.TokenDash},
		{Kind: lexer.TokenKeyword},
		{Kind: lexer.TokenColon},
		{Kind: lexer.TokenIdentifier},
		{Kind: lexer.TokenNewline},
	}

	file := test_util.OpenTestFile(t, test_util.SetPrintYAP, testFileDirPrefix)
	defer file.Close()

	lex := lexer.NewLexer(file, test_util.SetPrintYAP)
	toks, err := lex.Lex()
	assert.Nil(t, err)
	for i, tok := range toks {
		assert.Equal(t, expectedTok[i].Kind.String(), tok.Kind.String(), i)
	}
}

func TestLexSetPrintBinaryExp(t *testing.T) {
	_ = `
- set: 
  - x: 10 + 10 - 15
  - y: x * 4 
  - z: y / 5
- print: x
- print: y
- print: "hello" + " " + "world!"
`
	// Expected tokens for binary expressions:
	// Line 1: - set: \n
	// Line 2: - x: 10 + 10 - 15 \n
	// Line 3: - y: x * 4 \n
	// Line 4: - z: y / 5 \n
	// Line 5: - print: x \n
	// Line 6: - print: y \n
	// Line 7: - print: "hello" + " " + "world!" \n

	expectedTok := []lexer.Token{
		// - set:
		{Kind: lexer.TokenDash},
		{Kind: lexer.TokenKeyword, Value: lexer.KeywordSet},
		{Kind: lexer.TokenColon},
		{Kind: lexer.TokenNewline},
		{Kind: lexer.TokenIndent},
		// - x: 10 + 10 - 15
		{Kind: lexer.TokenDash},
		{Kind: lexer.TokenIdentifier, Value: "x"},
		{Kind: lexer.TokenColon},
		{Kind: lexer.TokenNumerical, Value: "10"},
		{Kind: lexer.TokenOperator, Value: "+"},
		{Kind: lexer.TokenNumerical, Value: "10"},
		{Kind: lexer.TokenOperator, Value: "-"},
		{Kind: lexer.TokenNumerical, Value: "15"},
		{Kind: lexer.TokenNewline},
		// - y: x * 4
		{Kind: lexer.TokenDash},
		{Kind: lexer.TokenIdentifier, Value: "y"},
		{Kind: lexer.TokenColon},
		{Kind: lexer.TokenIdentifier, Value: "x"},
		{Kind: lexer.TokenOperator, Value: "*"},
		{Kind: lexer.TokenNumerical, Value: "4"},
		{Kind: lexer.TokenNewline},
		// - z: y / 5
		{Kind: lexer.TokenDash},
		{Kind: lexer.TokenIdentifier, Value: "z"},
		{Kind: lexer.TokenColon},
		{Kind: lexer.TokenIdentifier, Value: "y"},
		{Kind: lexer.TokenOperator, Value: "/"},
		{Kind: lexer.TokenNumerical, Value: "5"},
		{Kind: lexer.TokenNewline},
		{Kind: lexer.TokenDedent},
		// - print: x
		{Kind: lexer.TokenDash},
		{Kind: lexer.TokenKeyword, Value: lexer.KeywordPrint},
		{Kind: lexer.TokenColon},
		{Kind: lexer.TokenIdentifier, Value: "x"},
		{Kind: lexer.TokenNewline},
		// - print: y
		{Kind: lexer.TokenDash},
		{Kind: lexer.TokenKeyword, Value: lexer.KeywordPrint},
		{Kind: lexer.TokenColon},
		{Kind: lexer.TokenIdentifier, Value: "y"},
		{Kind: lexer.TokenNewline},
		// - print: "hello" + " " + "world!"
		{Kind: lexer.TokenDash},
		{Kind: lexer.TokenKeyword, Value: lexer.KeywordPrint},
		{Kind: lexer.TokenColon},
		{Kind: lexer.TokenString, Value: "hello"},
		{Kind: lexer.TokenOperator, Value: "+"},
		{Kind: lexer.TokenString, Value: " "},
		{Kind: lexer.TokenOperator, Value: "+"},
		{Kind: lexer.TokenString, Value: "world!"},
		{Kind: lexer.TokenNewline},
	}

	file := test_util.OpenTestFile(t, test_util.SetPrintBinaryExpYAP, testFileDirPrefix)
	defer file.Close()

	lex := lexer.NewLexer(file, test_util.SetPrintBinaryExpYAP)
	toks, err := lex.Lex()
	assert.Nil(t, err)
	assert.Equal(t, len(expectedTok), len(toks), "token count mismatch")
	for i, tok := range toks {
		assert.Equal(t, expectedTok[i].Kind.String(), tok.Kind.String(), "token kind mismatch at %d", i)
		if expectedTok[i].Value != "" {
			assert.Equal(t, expectedTok[i].Value, tok.Value, "token value mismatch at %d", i)
		}
	}
}

func TestLexBooleanComparison(t *testing.T) {
	file := test_util.OpenTestFile(t, test_util.BooleanComparisonYAP, testFileDirPrefix)
	defer file.Close()

	lex := lexer.NewLexer(file, test_util.BooleanComparisonYAP)
	toks, err := lex.Lex()
	assert.Nil(t, err)

	// Verify specific tokens for comparison operators and booleans
	// Find and verify key tokens: >, ==, <=, !=, >=, <, True, False

	// Helper to find token by value
	findToken := func(value string) bool {
		for _, tok := range toks {
			if tok.Value == value {
				return true
			}
		}
		return false
	}

	// Verify comparison operators are correctly lexed
	assert.True(t, findToken(">"), "should have > operator")
	assert.True(t, findToken("=="), "should have == operator")
	assert.True(t, findToken("<="), "should have <= operator")
	assert.True(t, findToken("!="), "should have != operator")
	assert.True(t, findToken(">="), "should have >= operator")
	assert.True(t, findToken("<"), "should have < operator")

	// Verify boolean keywords are correctly lexed
	assert.True(t, findToken(lexer.KeywordTrue), "should have True keyword")
	assert.True(t, findToken(lexer.KeywordFalse), "should have False keyword")

	// Verify specific two-character operators are parsed correctly (not as two tokens)
	// Check that we have the == operator as a single token
	for _, tok := range toks {
		if tok.Value == "==" {
			assert.Equal(t, lexer.TokenOperator, tok.Kind, "== should be a single operator token")
		}
		if tok.Value == "!=" {
			assert.Equal(t, lexer.TokenOperator, tok.Kind, "!= should be a single operator token")
		}
		if tok.Value == "<=" {
			assert.Equal(t, lexer.TokenOperator, tok.Kind, "<= should be a single operator token")
		}
		if tok.Value == ">=" {
			assert.Equal(t, lexer.TokenOperator, tok.Kind, ">= should be a single operator token")
		}
	}
}

// Test file content:
// - print: "hello"
// //print: "hello"
// // print: "hello"
// - print: "hello" // world
func TestLexCommentsIgnored(t *testing.T) {
	file := test_util.OpenTestFile(t, test_util.CommentsIgnoredYAP, testFileDirPrefix)
	defer file.Close()

	lex := lexer.NewLexer(file, test_util.CommentsIgnoredYAP)
	toks, err := lex.Lex()
	assert.Nil(t, err)

	// Expected tokens for the comments test file:
	// Line 1: - print: "hello" \n (5 tokens: Dash, Keyword, Colon, String, Newline)
	// Line 2: //print: "hello" -> completely ignored, only newline emitted (1 token: Newline)
	// Line 3: // print: "hello" -> completely ignored, only newline emitted (1 token: Newline)
	// Line 4: - print: "hello" // world -> (5 tokens: Dash, Keyword, Colon, String, Newline)
	// Total: 12 tokens

	expectedTok := []lexer.Token{
		// Line 1: - print: "hello"
		{Kind: lexer.TokenDash, Value: "-"},
		{Kind: lexer.TokenKeyword, Value: lexer.KeywordPrint},
		{Kind: lexer.TokenColon, Value: ":"},
		{Kind: lexer.TokenString, Value: "hello"},
		{Kind: lexer.TokenNewline, Value: ""},
		// Line 2: //print: "hello" - full line comment, only newline
		{Kind: lexer.TokenComment, Value: ""},
		{Kind: lexer.TokenNewline, Value: ""},
		// Line 3: // print: "hello" - full line comment, only newline
		{Kind: lexer.TokenComment, Value: ""},
		{Kind: lexer.TokenNewline, Value: ""},
		// Line 4: - print: "hello" // world
		{Kind: lexer.TokenDash, Value: "-"},
		{Kind: lexer.TokenKeyword, Value: lexer.KeywordPrint},
		{Kind: lexer.TokenColon, Value: ":"},
		{Kind: lexer.TokenString, Value: "hello"},
		{Kind: lexer.TokenComment, Value: ""},
		{Kind: lexer.TokenNewline, Value: ""},
		// Line 5: - s- set:
		{Kind: lexer.TokenDash, Value: "-"},
		{Kind: lexer.TokenKeyword, Value: lexer.KeywordSet},
		{Kind: lexer.TokenColon, Value: ":"},
		{Kind: lexer.TokenNewline, Value: ""},
		//   - x: 10
		{Kind: lexer.TokenIndent, Value: ""},
		{Kind: lexer.TokenDash, Value: "-"},
		{Kind: lexer.TokenIdentifier, Value: "x"},
		{Kind: lexer.TokenColon, Value: ":"},
		{Kind: lexer.TokenNumerical, Value: "10"},
		{Kind: lexer.TokenNewline, Value: ""},
		//   - y: 5
		{Kind: lexer.TokenDash, Value: "-"},
		{Kind: lexer.TokenIdentifier, Value: "y"},
		{Kind: lexer.TokenColon, Value: ":"},
		{Kind: lexer.TokenNumerical, Value: "5"},
		{Kind: lexer.TokenNewline, Value: ""},
		// - print: x // + y
		{Kind: lexer.TokenDedent, Value: ""},
		{Kind: lexer.TokenDash, Value: "-"},
		{Kind: lexer.TokenKeyword, Value: lexer.KeywordPrint},
		{Kind: lexer.TokenColon, Value: ":"},
		{Kind: lexer.TokenIdentifier, Value: "x"},
		{Kind: lexer.TokenComment, Value: ""},
		{Kind: lexer.TokenNewline, Value: ""},
	}

	assert.Equal(t, len(expectedTok), len(toks), "token count mismatch")
	for i, tok := range toks {
		assert.Equal(t, expectedTok[i].Kind.String(), tok.Kind.String(), "token kind mismatch at %d", i)
		if expectedTok[i].Value != "" {
			assert.Equal(t, expectedTok[i].Value, tok.Value, "token value mismatch at %d", i)
		}
	}

	// Verify that no comment content appears as tokens
	for _, tok := range toks {
		assert.NotEqual(t, "world", tok.Value, "comment content 'world' should not appear as a token")
		assert.NotContains(t, tok.Value, "//", "comment marker should not appear in any token value")
	}
}

// Test file content:
// - set:
//   - x: 10
//     // - y: 5
//
// - print: x
// - print: y
func TestLexCommentsIgnoreInBlock(t *testing.T) {
	file := test_util.OpenTestFile(t, test_util.CommentsIgnoreInBlockYAP, testFileDirPrefix)
	defer file.Close()

	lex := lexer.NewLexer(file, test_util.CommentsIgnoreInBlockYAP)
	toks, err := lex.Lex()
	assert.Nil(t, err)

	// Expected tokens:
	// Line 1: - set: \n
	// Line 2: INDENT - x: 10 \n
	// Line 3: // - y: 5 -> Comment, Newline
	// Line 4: DEDENT - print: x \n
	// Line 5: - print: y \n

	expectedTok := []lexer.Token{
		// Line 1: - set:
		{Kind: lexer.TokenDash, Value: "-"},
		{Kind: lexer.TokenKeyword, Value: lexer.KeywordSet},
		{Kind: lexer.TokenColon, Value: ":"},
		{Kind: lexer.TokenNewline, Value: ""},
		// Line 2: - x: 10 (indented)
		{Kind: lexer.TokenIndent, Value: ""},
		{Kind: lexer.TokenDash, Value: "-"},
		{Kind: lexer.TokenIdentifier, Value: "x"},
		{Kind: lexer.TokenColon, Value: ":"},
		{Kind: lexer.TokenNumerical, Value: "10"},
		{Kind: lexer.TokenNewline, Value: ""},
		// Line 3: // - y: 5 (comment line inside block)
		{Kind: lexer.TokenComment, Value: ""},
		{Kind: lexer.TokenNewline, Value: ""},
		// Line 4: - print: x (dedented)
		{Kind: lexer.TokenDedent, Value: ""},
		{Kind: lexer.TokenDash, Value: "-"},
		{Kind: lexer.TokenKeyword, Value: lexer.KeywordPrint},
		{Kind: lexer.TokenColon, Value: ":"},
		{Kind: lexer.TokenIdentifier, Value: "x"},
		{Kind: lexer.TokenNewline, Value: ""},
		// Line 5: - print: y
		{Kind: lexer.TokenDash, Value: "-"},
		{Kind: lexer.TokenKeyword, Value: lexer.KeywordPrint},
		{Kind: lexer.TokenColon, Value: ":"},
		{Kind: lexer.TokenIdentifier, Value: "y"},
		{Kind: lexer.TokenNewline, Value: ""},
	}

	assert.Equal(t, len(expectedTok), len(toks), "token count mismatch")
	for i, tok := range toks {
		assert.Equal(t, expectedTok[i].Kind.String(), tok.Kind.String(), "token kind mismatch at %d", i)
		if expectedTok[i].Value != "" {
			assert.Equal(t, expectedTok[i].Value, tok.Value, "token value mismatch at %d", i)
		}
	}
}
