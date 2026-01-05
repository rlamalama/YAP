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

// Test parsing binary expressions
func TestParseBinaryExpression(t *testing.T) {
	p := parser.NewParser(test_util.GetTestFilepath(test_util.SetPrintBinaryExpYAP, testFileDir))
	prog, err := p.Parse()

	assert.Nil(t, err)
	assert.NotNil(t, prog)
	// 1 set statement + 3 print statements = 4 statements
	assert.Equal(t, 4, len(prog.Statements))

	// First statement should be a set with 3 assignments
	setStmt := prog.Statements[0].(parser.SetStmt)
	assert.Equal(t, 3, len(setStmt.Assignment))

	// x: 10 + 10 - 15 should be a BinaryExpr
	xAssign := setStmt.Assignment[0]
	assert.Equal(t, "x", xAssign.Name)
	xExpr, ok := xAssign.Expr.(*parser.BinaryExpr)
	assert.True(t, ok, "x expression should be BinaryExpr")
	// The expression is left-associative: ((10 + 10) - 15)
	assert.Equal(t, "-", xExpr.Operator)
	innerExpr, ok := xExpr.Left.(*parser.BinaryExpr)
	assert.True(t, ok, "inner expression should be BinaryExpr")
	assert.Equal(t, "+", innerExpr.Operator)

	// y: x * 4 should be a BinaryExpr
	yAssign := setStmt.Assignment[1]
	assert.Equal(t, "y", yAssign.Name)
	yExpr, ok := yAssign.Expr.(*parser.BinaryExpr)
	assert.True(t, ok, "y expression should be BinaryExpr")
	assert.Equal(t, "*", yExpr.Operator)

	// z: y / 5 should be a BinaryExpr
	zAssign := setStmt.Assignment[2]
	assert.Equal(t, "z", zAssign.Name)
	zExpr, ok := zAssign.Expr.(*parser.BinaryExpr)
	assert.True(t, ok, "z expression should be BinaryExpr")
	assert.Equal(t, "/", zExpr.Operator)

	// Third print statement: "hello" + " " + "world!"
	printStmt := prog.Statements[3].(parser.PrintStmt)
	concatExpr, ok := printStmt.Expr.(*parser.BinaryExpr)
	assert.True(t, ok, "print expression should be BinaryExpr")
	assert.Equal(t, "+", concatExpr.Operator)
}
