package test

import (
	"path/filepath"
	"testing"

	"github.com/rlamalama/YAP/cmd/yap/commands"
	test_util "github.com/rlamalama/YAP/test/test-util"
	"github.com/stretchr/testify/assert"
)

func TestOneLinePrint(t *testing.T) {
	filepath := filepath.Join(test_util.TestFiles, "0001-one-line-print.yap")
	args := []string{filepath}

	output := test_util.CaptureStdout(t, func() {
		commands.RunCmd(args)
	})

	expected := "hello world\n"
	assert.Equal(t, expected, output)
}
