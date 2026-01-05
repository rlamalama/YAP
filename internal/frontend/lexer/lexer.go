package lexer

import (
	"io"
	"strings"

	yaperror "github.com/rlamalama/YAP/internal/error"
)

type Lexer struct {
	filename    string
	scanner     *Scanner
	indentStack *Stack
	tokens      []*Token
}

func NewLexer(r io.Reader, filename string) *Lexer {
	return &Lexer{
		filename:    filename,
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

	if !ok {
		return yaperror.NewInvalidIndentError(l.filename, currLine, 1, indent, 0)
	}

	if indent > prevIndent {
		l.indentStack.Push(indent)
		l.emit(TokenIndent, "", currLine, indent)
		return nil
	}

	if indent < prevIndent {
		for {
			prevIndent, ok = l.indentStack.Peek()
			if !ok {
				return yaperror.NewInvalidIndentError(l.filename, currLine, 1, indent, prevIndent)
			}
			if indent == prevIndent {
				break
			} else if indent < prevIndent {
				l.emit(TokenDedent, "", currLine, indent)
				prevIndent, _ = l.indentStack.Pop()
			}
			if l.indentStack.Length() == 0 {
				return yaperror.NewInvalidIndentError(l.filename, currLine, 1, indent, prevIndent)
			}
		}
		if indent != prevIndent {
			return yaperror.NewInvalidIndentError(l.filename, currLine, 1, indent, prevIndent)
		}
	}

	return nil
}

func (l *Lexer) lexLine(line string, indent int) error {
	i := indent
	col := indent + 1

	if isDash(line[i]) {
		l.emit(TokenDash, "-", l.scanner.line, col)
		i++
		col++
	}

lineLoop:
	for i < len(line) {

		switch {
		case isTab(line[i]):
			return yaperror.NewTabCharError(l.filename, l.scanner.line, col)

		// Whitespace is not considered a token and ignored
		case isSpace(line[i]):
			i++
			col++

		case isColon(line[i]):
			l.emit(TokenColon, ":", l.scanner.line, col)
			i++
			col++

		// Keyword or Identifier
		case isAlpha(line[i]):
			start := i
			for i < len(line) && isAlphaNum(line[i]) {
				i++
			}
			val := line[start:i]
			tk := TokenIdentifier
			if IsKeyword(val) {
				tk = TokenKeyword
			}
			l.emit(tk, val, l.scanner.line, col)
			col += i - start

		case isQuote(line[i]):
			startCol := col
			start := i + 1
			i++ // consume opening quote

			for i < len(line) && !isQuote(line[i]) {
				i++
			}

			if i >= len(line) {
				return yaperror.NewUnterminatedStringError(l.filename, l.scanner.line, startCol)
			}

			value := line[start:i]
			l.emit(TokenString, value, l.scanner.line, col)

			i++ // consume closing quote
			col += (i - start) + 1
		case isNum(line[i]):
			start := i
			for i < len(line) && isNum(line[i]) {
				i++
			}
			l.emit(TokenNumerical, line[start:i], l.scanner.line, col)
			col += i - start
		case isComment(line[i]) && i < len(line)-1 && isComment(line[i+1]):
			l.emit(TokenComment, "", l.scanner.line, col)
			break lineLoop
		case StartsWithOperator(line[i]):
			// Check for two-character comparison operators first
			if i < len(line)-1 {
				twoChar := line[i : i+2]
				if IsComparisonOperator(twoChar) {
					l.emit(TokenOperator, twoChar, l.scanner.line, col)
					i += 2
					col += 2
					continue
				}
			}
			// Single character operator (arithmetic or single-char comparison)
			l.emit(TokenOperator, string(line[i]), l.scanner.line, col)
			i++
			col++
		default:
			return yaperror.NewInvalidTokenError(l.filename, l.scanner.line, col)
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

func isNum(c byte) bool {
	return (c >= '0' && c <= '9')
}

func isAlphaNum(c byte) bool {
	return isAlpha(c) || isNum(c)
}

func isSpace(c byte) bool {
	return c == ' '
}

func isTab(c byte) bool {
	return c == '\t'
}

func countIndent(s string) int {
	n := 0
	for n < len(s) && s[n] == ' ' {
		n++
	}
	return n
}

func isBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}

func isDash(c byte) bool {
	return c == '-'
}

func isColon(c byte) bool {
	return c == ':'
}

func isQuote(c byte) bool {
	return c == '"'
}

// Comment is //
func isComment(c byte) bool {
	return c == '/'
}
