package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const TestFileDir = "test/test-files"

// openTestFile caller should defer file.Close()
func OpenTestFile(t *testing.T, testFile string, prefix string) *os.File {
	filepath := filepath.Join(prefix, TestFileDir, testFile)
	file, err := os.Open(filepath)
	require.Nil(t, err)
	require.NotNil(t, file)
	return file
}
