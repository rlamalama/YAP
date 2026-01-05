package build_test

import (
	"testing"

	"github.com/rlamalama/YAP/internal/backend/build"
	"github.com/rlamalama/YAP/internal/backend/ir"
	"github.com/rlamalama/YAP/internal/frontend/parser"
	"github.com/stretchr/testify/require"
)

func TestBuildPrint(t *testing.T) {
	stmts := []parser.Stmt{
		parser.PrintStmt{Expr: &parser.StringLiteral{Value: "hello"}},
	}

	builder := build.New()
	irs, err := builder.Build(stmts)
	require.NoError(t, err)

	require.Equal(t, ir.OpPrint, irs[0].Op)
	require.Equal(t, ir.OperandLiteral, irs[0].Arg.Kind)
	require.Equal(t, "hello", irs[0].Arg.Value)
}

func TestBuildSet(t *testing.T) {
	name, val := "x", "Hello"
	stmts := []parser.Stmt{
		parser.SetStmt{
			Assignment: []*parser.Assignment{
				{
					Name: name,
					Expr: &parser.StringLiteral{
						Value: val,
					},
				},
			},
		},
	}

	builder := build.New()
	irs, err := builder.Build(stmts)
	require.NoError(t, err)

	require.Equal(t, ir.OpSet, irs[0].Op)
	require.Equal(t, ir.OperandIdentifier, irs[0].Arg.Kind)
	require.Equal(t, name+"="+val, irs[0].Arg.Value)
}
func TestBuildSetPrint(t *testing.T) {
	name, val := "x", "Hello"
	stmts := []parser.Stmt{
		parser.SetStmt{
			Assignment: []*parser.Assignment{
				{
					Name: name,
					Expr: &parser.StringLiteral{
						Value: val,
					},
				},
			},
		},
		parser.PrintStmt{
			Expr: &parser.Identifier{
				Name: name,
			},
		},
	}

	builder := build.New()
	irs, err := builder.Build(stmts)
	require.NoError(t, err)
	require.Equal(t, 2, len(irs))

	require.Equal(t, ir.OpSet, irs[0].Op)
	require.Equal(t, ir.OperandIdentifier, irs[0].Arg.Kind)
	require.Equal(t, name+"="+val, irs[0].Arg.Value)

	require.Equal(t, ir.OpPrint, irs[1].Op)
	require.Equal(t, ir.OperandIdentifier, irs[1].Arg.Kind)
	require.Equal(t, name, irs[1].Arg.Value)
}
