package test

import (
	"path/filepath"
	"testing"

	"github.com/rlamalama/YAP/cmd/yap/commands"
	test_util "github.com/rlamalama/YAP/test/test-util"
	"github.com/stretchr/testify/assert"
)

func TestMultiLinePrint(t *testing.T) {
	filepath := filepath.Join(test_util.TestFilesDir, test_util.MultiLinePrintYAP)
	args := []string{filepath}

	output := test_util.CaptureStdout(t, func() {
		commands.RunCmd(args)
	})

	expected :=
		`hello world
lots
lots
test
123ff
`
	assert.Equal(t, expected, output)
}
