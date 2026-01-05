package parser_test

import (
	"testing"

	"github.com/rlamalama/YAP/internal/frontend/parser"
	test_util "github.com/rlamalama/YAP/test/test-util"
	"github.com/stretchr/testify/assert"
)

const testFileDir = "../../.."

func TestParseEmptyFile(t *testing.T) {
	p := parser.NewParser(test_util.GetTestFilepath(test_util.EmptyFileYAP, testFileDir))
	prog, err := p.Parse()

	// Expecting 0 tokens in an empty file
	assert.Nil(t, err)
	assert.NotNil(t, prog)
	assert.Equal(t, 0, len(prog.Statements))
}

// - print: "hello world"
func TestParseOneLinePrintStatement(t *testing.T) {
	p := parser.NewParser(test_util.GetTestFilepath(test_util.OneLinePrintYAP, testFileDir))
	prog, err := p.Parse()

	// Expecting 1 in an empty file
	assert.Nil(t, err)
	assert.NotNil(t, prog)
	assert.Equal(t, 1, len(prog.Statements))

	stmt := prog.Statements[0]
	assert.Equal(t, parser.StmtTypePrint, stmt.Type())
	printStmt := stmt.(parser.PrintStmt)
	assert.Equal(t, "hello world", printStmt.Expr.String())
}

// 5 print statements
func TestParseMultiLinePrintStatement(t *testing.T) {
	p := parser.NewParser(test_util.GetTestFilepath(test_util.MultiLinePrintYAP, testFileDir))
	prog, err := p.Parse()

	// Expecting 1 in an empty file
	assert.Nil(t, err)
	assert.NotNil(t, prog)
	assert.Equal(t, 5, len(prog.Statements))
}
