package parser

import (
	"bufio"
	"io"
)

type Scanner struct {
	r    *bufio.Scanner
	line int
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		r:    bufio.NewScanner(r),
		line: 0,
	}
}

func (s *Scanner) NextLine() (string, bool) {
	if !s.r.Scan() {
		return "", false
	}
	s.line++
	return s.r.Text(), true
}
