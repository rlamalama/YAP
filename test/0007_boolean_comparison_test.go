package test

import (
	"path/filepath"
	"testing"

	"github.com/rlamalama/YAP/cmd/yap/commands"
	test_util "github.com/rlamalama/YAP/test/test-util"
	"github.com/stretchr/testify/assert"
)

func TestBooleanComparison(t *testing.T) {
	filepath := filepath.Join(test_util.TestFilesDir, test_util.BooleanComparisonYAP)
	args := []string{filepath}

	output := test_util.CaptureStdout(t, func() {
		commands.RunCmd(args)
	})

	expected :=
		`true
false
true
true
true
false
true
false
`
	assert.Equal(t, expected, output)
}
