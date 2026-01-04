package test

import (
	"path/filepath"
	"testing"

	"github.com/rlamalama/YAP/cmd/yap/commands"
	test_util "github.com/rlamalama/YAP/test/test-util"
)

func TestEmptyFile(t *testing.T) {
	filepath := filepath.Join(test_util.TestFiles, "0000-empty-file.yap")
	args := []string{filepath}

	commands.RunCmd(args)
}
