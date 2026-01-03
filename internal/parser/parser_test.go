package parser_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const testFileDir = "../../test/test-files"

// openTestFile caller should defer file.Close()
func OpenTestFile(t *testing.T, testFile string) *os.File {
	filepath := filepath.Join(testFileDir, testFile)
	file, err := os.Open(filepath)
	require.Nil(t, err)
	require.NotNil(t, file)
	return file
}
