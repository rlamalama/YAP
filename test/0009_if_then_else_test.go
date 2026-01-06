package test

import (
	"path/filepath"
	"testing"

	"github.com/rlamalama/YAP/cmd/yap/commands"
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
