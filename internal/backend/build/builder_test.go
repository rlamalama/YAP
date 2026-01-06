package build_test

import (
	"testing"

	"github.com/rlamalama/YAP/internal/backend/build"
	"github.com/rlamalama/YAP/internal/backend/ir"
	"github.com/rlamalama/YAP/internal/frontend/parser"
	"github.com/stretchr/testify/require"
)

func TestBuildPrint(t *testing.T) {
	expr := &parser.StringLiteral{Value: "hello"}
	stmts := []parser.Stmt{
		parser.PrintStmt{Expr: expr},
	}

	builder := build.New()
	irs, err := builder.Build(stmts)
	require.NoError(t, err)

	require.Equal(t, ir.OpPrint, irs[0].Op)
	require.Equal(t, expr, irs[0].Expr)
}

func TestBuildSet(t *testing.T) {
	name, val := "x", "Hello"
	expr := &parser.StringLiteral{Value: val}
	stmts := []parser.Stmt{
		parser.SetStmt{
			Assignment: []*parser.Assignment{
				{
					Name: name,
					Expr: expr,
				},
			},
		},
	}

	builder := build.New()
	irs, err := builder.Build(stmts)
	require.NoError(t, err)

	require.Equal(t, ir.OpSet, irs[0].Op)
	require.Equal(t, ir.OperandIdentifier, irs[0].Arg.Kind)
	require.Equal(t, name, irs[0].Arg.Value)
	require.Equal(t, expr, irs[0].Expr)
}
func TestBuildSetPrint(t *testing.T) {
	name, val := "x", "Hello"
	setExpr := &parser.StringLiteral{Value: val}
	printExpr := &parser.Identifier{Name: name}
	stmts := []parser.Stmt{
		parser.SetStmt{
			Assignment: []*parser.Assignment{
				{
					Name: name,
					Expr: setExpr,
				},
			},
		},
		parser.PrintStmt{
			Expr: printExpr,
		},
	}

	builder := build.New()
	irs, err := builder.Build(stmts)
	require.NoError(t, err)
	require.Equal(t, 2, len(irs))

	require.Equal(t, ir.OpSet, irs[0].Op)
	require.Equal(t, ir.OperandIdentifier, irs[0].Arg.Kind)
	require.Equal(t, name, irs[0].Arg.Value)
	require.Equal(t, setExpr, irs[0].Expr)

	require.Equal(t, ir.OpPrint, irs[1].Op)
	require.Equal(t, printExpr, irs[1].Expr)
}

func TestBuildBinaryExpr(t *testing.T) {
	// Test building: x = 10 + 5
	binExpr := &parser.BinaryExpr{
		Left:     &parser.NumericLiteral{Value: 10},
		Operator: "+",
		Right:    &parser.NumericLiteral{Value: 5},
	}
	stmts := []parser.Stmt{
		parser.SetStmt{
			Assignment: []*parser.Assignment{
				{
					Name: "x",
					Expr: binExpr,
				},
			},
		},
		parser.PrintStmt{
			Expr: &parser.Identifier{Name: "x"},
		},
	}

	builder := build.New()
	irs, err := builder.Build(stmts)
	require.NoError(t, err)
	require.Equal(t, 2, len(irs))

	// Check set instruction
	require.Equal(t, ir.OpSet, irs[0].Op)
	require.Equal(t, "x", irs[0].Arg.Value)
	require.Equal(t, binExpr, irs[0].Expr)

	// Check print instruction
	require.Equal(t, ir.OpPrint, irs[1].Op)
}

func TestBuildChainedBinaryExpr(t *testing.T) {
	// Test building: x = 10 + 10 - 15 (which is ((10 + 10) - 15))
	innerExpr := &parser.BinaryExpr{
		Left:     &parser.NumericLiteral{Value: 10},
		Operator: "+",
		Right:    &parser.NumericLiteral{Value: 10},
	}
	outerExpr := &parser.BinaryExpr{
		Left:     innerExpr,
		Operator: "-",
		Right:    &parser.NumericLiteral{Value: 15},
	}

	stmts := []parser.Stmt{
		parser.SetStmt{
			Assignment: []*parser.Assignment{
				{Name: "x", Expr: outerExpr},
			},
		},
	}

	builder := build.New()
	irs, err := builder.Build(stmts)
	require.NoError(t, err)
	require.Equal(t, 1, len(irs))

	require.Equal(t, ir.OpSet, irs[0].Op)
	require.Equal(t, "x", irs[0].Arg.Value)
	require.Equal(t, outerExpr, irs[0].Expr)
}

func TestBuildPrintBinaryExpr(t *testing.T) {
	// Test building: print x * 4
	binExpr := &parser.BinaryExpr{
		Left:     &parser.Identifier{Name: "x"},
		Operator: "*",
		Right:    &parser.NumericLiteral{Value: 4},
	}
	stmts := []parser.Stmt{
		parser.PrintStmt{Expr: binExpr},
	}

	builder := build.New()
	irs, err := builder.Build(stmts)
	require.NoError(t, err)
	require.Equal(t, 1, len(irs))

	require.Equal(t, ir.OpPrint, irs[0].Op)
	require.Equal(t, binExpr, irs[0].Expr)
}

func TestBuildBooleanLiteral(t *testing.T) {
	// Test building: flag = True, print flag
	boolExpr := &parser.BooleanLiteral{Value: true}
	stmts := []parser.Stmt{
		parser.SetStmt{
			Assignment: []*parser.Assignment{
				{
					Name: "flag",
					Expr: boolExpr,
				},
			},
		},
		parser.PrintStmt{
			Expr: &parser.Identifier{Name: "flag"},
		},
	}

	builder := build.New()
	irs, err := builder.Build(stmts)
	require.NoError(t, err)
	require.Equal(t, 2, len(irs))

	require.Equal(t, ir.OpSet, irs[0].Op)
	require.Equal(t, "flag", irs[0].Arg.Value)
	require.Equal(t, boolExpr, irs[0].Expr)

	require.Equal(t, ir.OpPrint, irs[1].Op)
}

func TestBuildComparisonExpr(t *testing.T) {
	// Test building: isGreater = a > b
	compExpr := &parser.BinaryExpr{
		Left:     &parser.Identifier{Name: "a"},
		Operator: ">",
		Right:    &parser.Identifier{Name: "b"},
	}
	stmts := []parser.Stmt{
		parser.SetStmt{
			Assignment: []*parser.Assignment{
				{
					Name: "isGreater",
					Expr: compExpr,
				},
			},
		},
	}

	builder := build.New()
	irs, err := builder.Build(stmts)
	require.NoError(t, err)
	require.Equal(t, 1, len(irs))

	require.Equal(t, ir.OpSet, irs[0].Op)
	require.Equal(t, "isGreater", irs[0].Arg.Value)
	require.Equal(t, compExpr, irs[0].Expr)
}

func TestBuildPrintComparisonExpr(t *testing.T) {
	// Test building: print a >= 10
	compExpr := &parser.BinaryExpr{
		Left:     &parser.Identifier{Name: "a"},
		Operator: ">=",
		Right:    &parser.NumericLiteral{Value: 10},
	}
	stmts := []parser.Stmt{
		parser.PrintStmt{Expr: compExpr},
	}

	builder := build.New()
	irs, err := builder.Build(stmts)
	require.NoError(t, err)
	require.Equal(t, 1, len(irs))

	require.Equal(t, ir.OpPrint, irs[0].Op)
	require.Equal(t, compExpr, irs[0].Expr)
}

func TestBuildIfThenElse(t *testing.T) {
	// Test building: if x > 5 then print "big" else print "small"
	condition := &parser.BinaryExpr{
		Left:     &parser.Identifier{Name: "x"},
		Operator: ">",
		Right:    &parser.NumericLiteral{Value: 5},
	}
	stmts := []parser.Stmt{
		parser.IfStmt{
			Condition: condition,
			Then: []parser.Stmt{
				parser.PrintStmt{Expr: &parser.StringLiteral{Value: "big"}},
			},
			Else: []parser.Stmt{
				parser.PrintStmt{Expr: &parser.StringLiteral{Value: "small"}},
			},
		},
	}

	builder := build.New()
	irs, err := builder.Build(stmts)
	require.NoError(t, err)

	// Expected instructions:
	// 0: JumpIfFalse (condition, jump to 3 if false)
	// 1: Print "big" (then block)
	// 2: Jump (to 4, skip else block)
	// 3: Print "small" (else block)
	require.Equal(t, 4, len(irs))

	// Check JumpIfFalse
	require.Equal(t, ir.OpJumpIfFalse, irs[0].Op)
	require.Equal(t, condition, irs[0].Expr)
	require.Equal(t, 3, irs[0].Arg.Offset) // Jump to else block

	// Check then block print
	require.Equal(t, ir.OpPrint, irs[1].Op)

	// Check Jump (skip else)
	require.Equal(t, ir.OpJump, irs[2].Op)
	require.Equal(t, 4, irs[2].Arg.Offset) // Jump past else block

	// Check else block print
	require.Equal(t, ir.OpPrint, irs[3].Op)
}

func TestBuildIfThenNoElse(t *testing.T) {
	// Test building: if x > 5 then print "big" (no else)
	condition := &parser.BinaryExpr{
		Left:     &parser.Identifier{Name: "x"},
		Operator: ">",
		Right:    &parser.NumericLiteral{Value: 5},
	}
	stmts := []parser.Stmt{
		parser.IfStmt{
			Condition: condition,
			Then: []parser.Stmt{
				parser.PrintStmt{Expr: &parser.StringLiteral{Value: "big"}},
			},
			Else: nil, // No else block
		},
	}

	builder := build.New()
	irs, err := builder.Build(stmts)
	require.NoError(t, err)

	// Expected instructions:
	// 0: JumpIfFalse (condition, jump to 2 if false)
	// 1: Print "big" (then block)
	require.Equal(t, 2, len(irs))

	// Check JumpIfFalse
	require.Equal(t, ir.OpJumpIfFalse, irs[0].Op)
	require.Equal(t, 2, irs[0].Arg.Offset) // Jump past then block

	// Check then block print
	require.Equal(t, ir.OpPrint, irs[1].Op)
}

func TestBuildNestedIf(t *testing.T) {
	// Test building nested if statements
	innerCondition := &parser.BinaryExpr{
		Left:     &parser.Identifier{Name: "x"},
		Operator: "<",
		Right:    &parser.NumericLiteral{Value: 20},
	}
	outerCondition := &parser.BinaryExpr{
		Left:     &parser.Identifier{Name: "x"},
		Operator: ">",
		Right:    &parser.NumericLiteral{Value: 5},
	}
	stmts := []parser.Stmt{
		parser.IfStmt{
			Condition: outerCondition,
			Then: []parser.Stmt{
				parser.IfStmt{
					Condition: innerCondition,
					Then: []parser.Stmt{
						parser.PrintStmt{Expr: &parser.StringLiteral{Value: "medium"}},
					},
					Else: []parser.Stmt{
						parser.PrintStmt{Expr: &parser.StringLiteral{Value: "large"}},
					},
				},
			},
			Else: []parser.Stmt{
				parser.PrintStmt{Expr: &parser.StringLiteral{Value: "small"}},
			},
		},
	}

	builder := build.New()
	irs, err := builder.Build(stmts)
	require.NoError(t, err)

	// Expected instructions:
	// 0: JumpIfFalse outer (jump to 6 if false)
	// 1: JumpIfFalse inner (jump to 4 if false)
	// 2: Print "medium"
	// 3: Jump (to 5, skip inner else)
	// 4: Print "large"
	// 5: Jump (to 7, skip outer else)
	// 6: Print "small"
	require.Equal(t, 7, len(irs))

	// Check outer JumpIfFalse
	require.Equal(t, ir.OpJumpIfFalse, irs[0].Op)
	require.Equal(t, 6, irs[0].Arg.Offset)

	// Check inner JumpIfFalse
	require.Equal(t, ir.OpJumpIfFalse, irs[1].Op)
	require.Equal(t, 4, irs[1].Arg.Offset)
}
