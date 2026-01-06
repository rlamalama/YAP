package test

import (
	"path/filepath"
	"testing"

	"github.com/rlamalama/YAP/cmd/yap/commands"
	"github.com/rlamalama/YAP/internal/frontend/parser"
	test_util "github.com/rlamalama/YAP/test/test-util"
	"github.com/stretchr/testify/assert"
)

func TestIfThenElse(t *testing.T) {
	filepath := filepath.Join(test_util.TestFilesDir, test_util.IfThenElseYAP)
	args := []string{filepath}

	output := test_util.CaptureStdout(t, func() {
		commands.RunCmd(args)
	})

	expected :=
		`x is big
well not that big
`
	assert.Equal(t, expected, output)
}

// Test empty else block - should work fine (else block is optional and can be empty)
func TestEmptyElse(t *testing.T) {
	filepath := filepath.Join(test_util.TestFilesDir, test_util.EmptyElseYAP)
	args := []string{filepath}

	output := test_util.CaptureStdout(t, func() {
		commands.RunCmd(args)
	})

	// x = 10, x > 5 is true, so "x is big" prints, empty else does nothing
	expected := "x is big\n"
	assert.Equal(t, expected, output)
}

// Test empty then block - should work fine (then block can be empty)
func TestEmptyThen(t *testing.T) {
	filepath := filepath.Join(test_util.TestFilesDir, test_util.EmptyThenYAP)
	args := []string{filepath}

	output := test_util.CaptureStdout(t, func() {
		commands.RunCmd(args)
	})

	// x = 10, x > 5 is true, so empty then runs (nothing prints)
	// else block is not executed
	expected := ""
	assert.Equal(t, expected, output)
}

// Test hanging else (else without if) - should error
func TestHangingElseError(t *testing.T) {
	fp := filepath.Join(test_util.TestFilesDir, test_util.HangingElseYAP)

	p := parser.NewParser(fp)
	_, err := p.Parse()

	// Should error because else is not a valid top-level statement
	assert.NotNil(t, err, "hanging else should produce an error")
}

// Test hanging then (then without if) - should error
func TestHangingThenError(t *testing.T) {
	fp := filepath.Join(test_util.TestFilesDir, test_util.HangingThenYAP)

	p := parser.NewParser(fp)
	_, err := p.Parse()

	// Should error because then is not a valid top-level statement
	assert.NotNil(t, err, "hanging then should produce an error")
}
