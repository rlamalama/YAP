package parser

import (
	"fmt"
	"io"
	"strings"
)

type Lexer struct {
	scanner     *Scanner
	indentStack *Stack
	tokens      []*Token
}

func NewLexer(r io.Reader) *Lexer {

	return &Lexer{
		scanner:     NewScanner(r),
		indentStack: NewStackWithElem(0),
		tokens:      []*Token{},
	}
}

func (l *Lexer) Lex() ([]*Token, error) {
	for {
		line, ok := l.scanner.NextLine()
		if !ok {
			break
		}
		if isBlank(line) {
			continue
		}

		indent := countIndent(line)

		if err := l.handleIndent(indent); err != nil {
			return l.tokens, err
		}

		if err := l.lexLine(line, indent); err != nil {
			return l.tokens, err
		}
	}
	numIndent := l.indentStack.Length()
	for numIndent > 1 {
		l.indentStack.Pop()
		numIndent = l.indentStack.Length()
		l.emit(TokenDedent, "", l.scanner.line, 0)
	}

	return l.tokens, nil
}

func (l *Lexer) handleIndent(indent int) error {
	currLine := l.scanner.line
	prevIndent, ok := l.indentStack.Peek()

	fmt.Println("CURR INDENT --> ", prevIndent, "EVALUATED INDENT -->", indent, "CURR STACK", l.indentStack)
	if !ok {
		return fmt.Errorf("error parsing indents on line %d", currLine)
	}

	if indent > prevIndent {
		fmt.Println("PUSHINGGGGG")
		l.indentStack.Push(indent)
		l.emit(TokenIndent, "", currLine, indent)
		return nil
	}

	if indent < prevIndent {
		for {
			prevIndent, _ = l.indentStack.Peek()
			if !ok {
				return fmt.Errorf("error parsing indents on line %d", currLine)
			}
			if indent == prevIndent {
				break
			} else if indent < prevIndent {
				l.emit(TokenDedent, "", currLine, indent)
				prevIndent, _ = l.indentStack.Pop()
			}
			if l.indentStack.Length() == 0 {
				return fmt.Errorf("invalid formatting for indents on line %d with %d spaces vs %d prev", currLine, indent, prevIndent)
			}
		}
		if indent != prevIndent {
			return fmt.Errorf("invalid formatting for indents on line %d with %d spaces vs %d prev", currLine, indent, prevIndent)
		}
	}

	return nil
}

func (l *Lexer) lexLine(line string, indent int) error {
	i := indent
	col := indent + 1

	if line[i] == '-' {
		l.emit(TokenDash, "-", l.scanner.line, col)
		i++
		col++
	}

	for i < len(line) {
		switch {

		case isSpace(line[i]):
			i++
			col++

		case line[i] == ':':
			l.emit(TokenColon, ":", l.scanner.line, col)
			i++
			col++
		case isAlpha(line[i]):
			start := i
			for i < len(line) && isAlphaNum(line[i]) {
				i++
			}
			l.emit(TokenIdentifier, line[start:i], l.scanner.line, col)
			col += i - start
		case line[i] == '"':
			start := i + 1
			i++ // consume opening quote

			for i < len(line) && line[i] != '"' {
				i++
			}

			if i >= len(line) {
				return fmt.Errorf("unterminated string at line %d", l.scanner.line)
			}

			value := line[start:i]
			l.emit(TokenScalar, value, l.scanner.line, col)

			i++ // consume closing quote
			col += (i - start) + 1
		default:
			l.emit(TokenScalar, line[i:], l.scanner.line, col)
			i = len(line)
			col = i + 1
		}
	}
	l.emit(TokenNewline, "", l.scanner.line, col)

	return nil
}

func (l *Lexer) emit(tk TokenKind, val string, line, col int) error {
	token := NewToken(tk, val, line, col)
	l.tokens = append(l.tokens, token)
	return nil
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c == '_')
}

func isAlphaNum(c byte) bool {
	return isAlpha(c) || (c >= '0' && c <= '9')
}

func isSpace(c byte) bool {
	return c == ' '
}

func countIndent(s string) int {
	n := 0
	for n < len(s) && s[n] == ' ' {
		n++
	}
	return n
}

func isBlank(s string) bool {
	if strings.TrimSpace(s) == "" {
		return true
	}
	return false
}
