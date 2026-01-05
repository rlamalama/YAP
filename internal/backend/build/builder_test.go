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
