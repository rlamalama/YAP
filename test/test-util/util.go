package test_util

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const TestFileDir = "test/test-files"
const TestFiles = "test-files"

// openTestFile caller should defer file.Close()
func OpenTestFile(t *testing.T, testFile string, prefix string) *os.File {
	filepath := GetTestFilepath(testFile, prefix)
	file, err := os.Open(filepath)
	require.Nil(t, err)
	require.NotNil(t, file)
	return file
}

// openTestFile caller should defer file.Close()
func GetTestFilepath(testFile string, prefix string) string {
	return filepath.Join(prefix, TestFileDir, testFile)
}

func CaptureStdout(t *testing.T, fn func()) string {
	t.Helper()

	old := os.Stdout

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}

	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if err != nil {
		t.Fatalf("copy: %v", err)
	}

	return buf.String()
}
