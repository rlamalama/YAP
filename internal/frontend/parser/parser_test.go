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

// Test parsing boolean and comparison expressions
func TestParseBooleanComparison(t *testing.T) {
	p := parser.NewParser(test_util.GetTestFilepath(test_util.BooleanComparisonYAP, testFileDir))
	prog, err := p.Parse()

	assert.Nil(t, err)
	assert.NotNil(t, prog)
	// 1 set statement + 8 print statements = 9 statements
	assert.Equal(t, 9, len(prog.Statements))

	// First statement should be a set with 8 assignments
	setStmt := prog.Statements[0].(parser.SetStmt)
	assert.Equal(t, 8, len(setStmt.Assignment))

	// isGreater: a > b should be a BinaryExpr with ">" operator
	isGreaterAssign := setStmt.Assignment[2]
	assert.Equal(t, "isGreater", isGreaterAssign.Name)
	isGreaterExpr, ok := isGreaterAssign.Expr.(*parser.BinaryExpr)
	assert.True(t, ok, "isGreater expression should be BinaryExpr")
	assert.Equal(t, ">", isGreaterExpr.Operator)

	// isEqual: a == b should be a BinaryExpr with "==" operator
	isEqualAssign := setStmt.Assignment[3]
	assert.Equal(t, "isEqual", isEqualAssign.Name)
	isEqualExpr, ok := isEqualAssign.Expr.(*parser.BinaryExpr)
	assert.True(t, ok, "isEqual expression should be BinaryExpr")
	assert.Equal(t, "==", isEqualExpr.Operator)

	// isLessOrEqual: b <= a should be a BinaryExpr with "<=" operator
	isLessOrEqualAssign := setStmt.Assignment[4]
	assert.Equal(t, "isLessOrEqual", isLessOrEqualAssign.Name)
	isLessOrEqualExpr, ok := isLessOrEqualAssign.Expr.(*parser.BinaryExpr)
	assert.True(t, ok, "isLessOrEqual expression should be BinaryExpr")
	assert.Equal(t, "<=", isLessOrEqualExpr.Operator)

	// notEqual: a != b should be a BinaryExpr with "!=" operator
	notEqualAssign := setStmt.Assignment[5]
	assert.Equal(t, "notEqual", notEqualAssign.Name)
	notEqualExpr, ok := notEqualAssign.Expr.(*parser.BinaryExpr)
	assert.True(t, ok, "notEqual expression should be BinaryExpr")
	assert.Equal(t, "!=", notEqualExpr.Operator)

	// flag: True should be a BooleanLiteral with value true
	flagAssign := setStmt.Assignment[6]
	assert.Equal(t, "flag", flagAssign.Name)
	flagExpr, ok := flagAssign.Expr.(*parser.BooleanLiteral)
	assert.True(t, ok, "flag expression should be BooleanLiteral")
	assert.Equal(t, true, flagExpr.Value)

	// notFlag: False should be a BooleanLiteral with value false
	notFlagAssign := setStmt.Assignment[7]
	assert.Equal(t, "notFlag", notFlagAssign.Name)
	notFlagExpr, ok := notFlagAssign.Expr.(*parser.BooleanLiteral)
	assert.True(t, ok, "notFlag expression should be BooleanLiteral")
	assert.Equal(t, false, notFlagExpr.Value)

	// Last print statement: a >= 10 should be a BinaryExpr
	lastPrintStmt := prog.Statements[7].(parser.PrintStmt)
	gteExpr, ok := lastPrintStmt.Expr.(*parser.BinaryExpr)
	assert.True(t, ok, "print expression should be BinaryExpr")
	assert.Equal(t, ">=", gteExpr.Operator)
}

// Test parsing comments are ignored
// Test file content:
// - print: "hello"
// //print: "hello"
// // print: "hello"
// - print: "hello" // world
func TestParseCommentsIgnored(t *testing.T) {
	p := parser.NewParser(test_util.GetTestFilepath(test_util.CommentsIgnoredYAP, testFileDir))
	prog, err := p.Parse()

	assert.Nil(t, err)
	assert.NotNil(t, prog)
	// Lines 2 and 3 are comment-only lines
	assert.Equal(t, 4, len(prog.Statements))

	// First statement: - print: "hello"
	stmt1 := prog.Statements[0]
	assert.Equal(t, parser.StmtTypePrint, stmt1.Type())
	printStmt1 := stmt1.(parser.PrintStmt)
	strLit1, ok := printStmt1.Expr.(*parser.StringLiteral)
	assert.True(t, ok, "first print expression should be StringLiteral")
	assert.Equal(t, "hello", strLit1.Value)

	// Second statement: - print: "hello" // world
	// The "// world" comment should be ignored
	stmt2 := prog.Statements[1]
	assert.Equal(t, parser.StmtTypePrint, stmt2.Type())
	printStmt2 := stmt2.(parser.PrintStmt)
	strLit2, ok := printStmt2.Expr.(*parser.StringLiteral)
	assert.True(t, ok, "second print expression should be StringLiteral")
	assert.Equal(t, "hello", strLit2.Value)

	// set stmt simple
	stmt3 := prog.Statements[2]
	assert.Equal(t, parser.StmtTypeSet, stmt3.Type())

	// - print: x // + y
	stmt4 := prog.Statements[3]
	assert.Equal(t, parser.StmtTypePrint, stmt4.Type())
	printStmt3 := stmt4.(parser.PrintStmt)
	strLit3, ok := printStmt3.Expr.(*parser.Identifier)

	assert.True(t, ok, "second print expression should be StringLiteral")
	assert.Equal(t, "x", strLit3.Name)
}

// Test parsing comments inside a set block causes y to be undefined
// Test file content:
// - set:
//   - x: 10
//     // - y: 5
//
// - print: x
// - print: y
func TestParseCommentsIgnoreInBlock(t *testing.T) {
	p := parser.NewParser(test_util.GetTestFilepath(test_util.CommentsIgnoreInBlockYAP, testFileDir))
	prog, err := p.Parse()

	assert.Nil(t, err)
	assert.NotNil(t, prog)
	// 1 set statement + 2 print statements = 3 statements
	assert.Equal(t, 3, len(prog.Statements))

	// First statement should be a set with only 1 assignment (y is commented out)
	setStmt := prog.Statements[0].(parser.SetStmt)
	assert.Equal(t, 1, len(setStmt.Assignment), "only x should be assigned, y is commented out")
	assert.Equal(t, "x", setStmt.Assignment[0].Name)

	// Second statement: - print: x
	stmt2 := prog.Statements[1].(parser.PrintStmt)
	xIdent, ok := stmt2.Expr.(*parser.Identifier)
	assert.True(t, ok)
	assert.Equal(t, "x", xIdent.Name)

	// Third statement: - print: y (y is referenced but not defined)
	stmt3 := prog.Statements[2].(parser.PrintStmt)
	yIdent, ok := stmt3.Expr.(*parser.Identifier)
	assert.True(t, ok)
	assert.Equal(t, "y", yIdent.Name)
}

// Test parsing if-then-else statement
func TestParseIfThenElse(t *testing.T) {
	p := parser.NewParser(test_util.GetTestFilepath(test_util.IfThenElseYAP, testFileDir))
	prog, err := p.Parse()

	assert.Nil(t, err)
	assert.NotNil(t, prog)
	// 1 set statement + 1 if statement = 2 statements
	assert.Equal(t, 2, len(prog.Statements))

	// First statement should be a set
	setStmt := prog.Statements[0].(parser.SetStmt)
	assert.Equal(t, 1, len(setStmt.Assignment))
	assert.Equal(t, "x", setStmt.Assignment[0].Name)

	// Second statement should be an if statement
	ifStmt, ok := prog.Statements[1].(parser.IfStmt)
	assert.True(t, ok, "second statement should be IfStmt")
	assert.Equal(t, parser.StmtTypeIf, ifStmt.Type())

	// Verify condition is a binary expression (x > 5)
	condExpr, ok := ifStmt.Condition.(*parser.BinaryExpr)
	assert.True(t, ok, "condition should be BinaryExpr")
	assert.Equal(t, ">", condExpr.Operator)

	// Verify then block has 2 statements (print and nested if)
	assert.Equal(t, 2, len(ifStmt.Then), "then block should have 2 statements")

	// First then statement is print
	thenPrint, ok := ifStmt.Then[0].(parser.PrintStmt)
	assert.True(t, ok)
	strLit, ok := thenPrint.Expr.(*parser.StringLiteral)
	assert.True(t, ok)
	assert.Equal(t, "x is big", strLit.Value)

	// Second then statement is nested if
	nestedIf, ok := ifStmt.Then[1].(parser.IfStmt)
	assert.True(t, ok, "second then statement should be nested IfStmt")

	// Verify nested if condition (x < 20)
	nestedCond, ok := nestedIf.Condition.(*parser.BinaryExpr)
	assert.True(t, ok)
	assert.Equal(t, "<", nestedCond.Operator)

	// Verify nested if has then and else blocks
	assert.Equal(t, 1, len(nestedIf.Then))
	assert.Equal(t, 1, len(nestedIf.Else))

	// Verify outer else block has 1 statement
	assert.Equal(t, 1, len(ifStmt.Else), "else block should have 1 statement")
	elsePrint, ok := ifStmt.Else[0].(parser.PrintStmt)
	assert.True(t, ok)
	elseStrLit, ok := elsePrint.Expr.(*parser.StringLiteral)
	assert.True(t, ok)
	assert.Equal(t, "x is small", elseStrLit.Value)
}
