package parser_test

import (
	"testing"

	"github.com/rlamalama/YAP/internal/frontend/parser"
	test_util "github.com/rlamalama/YAP/test/test-util"
	"github.com/stretchr/testify/assert"
)

const testFileDir = "../../.."

func TestParseEmptyFile(t *testing.T) {
	p := parser.NewParser(test_util.GetTestFilepath("0000-empty-file.yap", testFileDir))
	prog, err := p.Parse()

	// Expecting 0 tokens in an empty file
	assert.Nil(t, err)
	assert.NotNil(t, prog)
	assert.Equal(t, 0, len(prog.Statements))
}

// - print: "hello world"
func TestParseOneLinePrintStatement(t *testing.T) {
	p := parser.NewParser(test_util.GetTestFilepath("0001-one-line-print.yap", testFileDir))
	prog, err := p.Parse()

	// Expecting 1 in an empty file
	assert.Nil(t, err)
	assert.NotNil(t, prog)
	assert.Equal(t, 1, len(prog.Statements))
	assert.Equal(t, 1, len(prog.Statements))

	stmt := prog.Statements[0]
	assert.Equal(t, parser.StmtTypePrint, stmt.Type())
	printStmt := stmt.(parser.PrintStmt)
	assert.Equal(t, "hello world", printStmt.Value)
}
