package test

import (
	"path/filepath"
	"testing"

	"github.com/rlamalama/YAP/internal/backend/build"
	"github.com/rlamalama/YAP/internal/backend/vm"
	"github.com/rlamalama/YAP/internal/frontend/parser"
	test_util "github.com/rlamalama/YAP/test/test-util"
	"github.com/stretchr/testify/assert"
)

// Test that commenting out a variable definition causes a runtime error
// when that variable is referenced
func TestCommentsIgnoreInBlockUndefinedVar(t *testing.T) {
	filepath := filepath.Join(test_util.TestFilesDir, test_util.CommentsIgnoreInBlockYAP)

	p := parser.NewParser(filepath)
	ast, err := p.Parse()
	assert.Nil(t, err, "parsing should succeed")
	assert.NotNil(t, ast)

	builder := build.New()
	program, err := builder.Build(ast.Statements)
	assert.Nil(t, err, "building should succeed")

	output := test_util.CaptureStdout(t, func() {
		v := vm.New(program)
		err = v.Run()

		// Should error because 	output := test_util.CaptureStdout(t, func() {is undefined (it was commented out)
		assert.NotNil(t, err, "should error because y is undefined")
		assert.Contains(t, err.Error(), "y", "error should mention undefined variable y")

	})
	expected :=
		`10
`
	assert.Equal(t, expected, output)

}
