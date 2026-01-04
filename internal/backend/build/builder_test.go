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
		parser.PrintStmt{Value: "hello"},
	}

	builder := build.New()
	irs, err := builder.Build(stmts)
	require.NoError(t, err)

	require.Equal(t, ir.OpPrint, irs[0].Op)
	require.Equal(t, "hello", irs[0].Arg)
}
