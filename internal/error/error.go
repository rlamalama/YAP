package error

import (
	"fmt"
	"strings"
)

// Severity represents the severity level of a diagnostic
type Severity int

const (
	SeverityError Severity = iota
	SeverityWarning
	SeverityNote
	SeverityHint
)

func (s Severity) String() string {
	switch s {
	case SeverityError:
		return "error"
	case SeverityWarning:
		return "warning"
	case SeverityNote:
		return "note"
	case SeverityHint:
		return "hint"
	default:
		return "unknown"
	}
}

// Phase represents which phase of compilation the error occurred in
type Phase int

const (
	PhaseLexer Phase = iota
	PhaseParser
	PhaseBuilder
	PhaseRuntime
)

func (p Phase) String() string {
	switch p {
	case PhaseLexer:
		return "lexer"
	case PhaseParser:
		return "parser"
	case PhaseBuilder:
		return "builder"
	case PhaseRuntime:
		return "runtime"
	default:
		return "unknown"
	}
}

// ErrorCode provides a unique identifier for each error type
type ErrorCode int

const (
	// Lexer errors (1000-1999)
	ErrInvalidCharacter ErrorCode = 1000 + iota
	ErrUnterminatedString
	ErrInvalidIndentation
	ErrTabCharacter
	ErrInvalidEscapeSequence
	ErrUnexpectedEOF
	ErrInvalidNumber

	// Parser errors (2000-2999)
	ErrUnexpectedToken ErrorCode = 2000 + iota
	ErrExpectedToken
	ErrUnknownStatement
	ErrMissingExpression
	ErrInvalidSyntax
	ErrUnexpectedEndOfInput
	ErrMissingColon
	ErrMissingValue

	// Builder/Semantic errors (3000-3999)
	ErrUnsupportedStatement ErrorCode = 3000 + iota
	ErrUndefinedVariable
	ErrUndefinedFunction
	ErrTypeMismatch
	ErrInvalidOperation
	ErrDuplicateDefinition
	ErrInvalidAssignment
	ErrInvalidArgCount

	// Runtime errors (4000-4999)
	ErrUnknownOpcode ErrorCode = 4000 + iota
	ErrDivisionByZero
	ErrStackUnderflow
	ErrStackOverflow
	ErrNullReference
	ErrOutOfBounds
	ErrInvalidType
	ErrIOError
)

// Position represents a location in the source code
type Position struct {
	File   string
	Line   int
	Column int
}

func (p Position) String() string {
	if p.File != "" {
		return fmt.Sprintf("%s:%d:%d", p.File, p.Line, p.Column)
	}
	return fmt.Sprintf("%d:%d", p.Line, p.Column)
}

// Span represents a range in the source code
type Span struct {
	Start Position
	End   Position
}

func (s Span) String() string {
	if s.Start.File != "" {
		return fmt.Sprintf("%s:%d:%d-%d:%d",
			s.Start.File, s.Start.Line, s.Start.Column,
			s.End.Line, s.End.Column)
	}
	return fmt.Sprintf("%d:%d-%d:%d",
		s.Start.Line, s.Start.Column,
		s.End.Line, s.End.Column)
}

// YapError is the main error type for the YAP compiler
type YapError struct {
	Code     ErrorCode
	Severity Severity
	Phase    Phase
	Position Position
	Span     *Span
	Message  string
	Context  string // The source line where the error occurred
	Notes    []string
}

// Error implements the error interface
func (e *YapError) Error() string {
	var sb strings.Builder

	// Format: file:line:col: severity: message
	sb.WriteString(e.Position.String())
	sb.WriteString(": ")
	sb.WriteString(e.Severity.String())
	sb.WriteString(": ")
	sb.WriteString(e.Message)

	return sb.String()
}

// FullError returns a detailed error message with context
func (e *YapError) FullError() string {
	var sb strings.Builder

	// Main error line
	sb.WriteString(e.Error())
	sb.WriteString("\n")

	// Source context with pointer
	if e.Context != "" {
		sb.WriteString("    ")
		sb.WriteString(e.Context)
		sb.WriteString("\n")

		// Add pointer to error location
		if e.Position.Column > 0 {
			sb.WriteString("    ")
			sb.WriteString(strings.Repeat(" ", e.Position.Column-1))
			sb.WriteString("^\n")
		}
	}

	// Additional notes
	for _, note := range e.Notes {
		sb.WriteString("note: ")
		sb.WriteString(note)
		sb.WriteString("\n")
	}

	return sb.String()
}

// AddNote adds a note to the error
func (e *YapError) AddNote(note string) *YapError {
	e.Notes = append(e.Notes, note)
	return e
}

// WithContext adds source context to the error
func (e *YapError) WithContext(line string) *YapError {
	e.Context = line
	return e
}

// WithSpan adds a span to the error
func (e *YapError) WithSpan(start, end Position) *YapError {
	e.Span = &Span{Start: start, End: end}
	return e
}

// IsError returns true if this is an error (not warning/note)
func (e *YapError) IsError() bool {
	return e.Severity == SeverityError
}

// IsWarning returns true if this is a warning
func (e *YapError) IsWarning() bool {
	return e.Severity == SeverityWarning
}

// ErrorList holds multiple errors for batch reporting
type ErrorList struct {
	errors []*YapError
}

// NewErrorList creates a new error list
func NewErrorList() *ErrorList {
	return &ErrorList{
		errors: make([]*YapError, 0),
	}
}

// Add adds an error to the list
func (el *ErrorList) Add(err *YapError) {
	el.errors = append(el.errors, err)
}

// AddError is a convenience method to add a simple error
func (el *ErrorList) AddError(code ErrorCode, phase Phase, pos Position, msg string) {
	el.Add(&YapError{
		Code:     code,
		Severity: SeverityError,
		Phase:    phase,
		Position: pos,
		Message:  msg,
	})
}

// AddWarning is a convenience method to add a warning
func (el *ErrorList) AddWarning(code ErrorCode, phase Phase, pos Position, msg string) {
	el.Add(&YapError{
		Code:     code,
		Severity: SeverityWarning,
		Phase:    phase,
		Position: pos,
		Message:  msg,
	})
}

// HasErrors returns true if there are any errors (not just warnings)
func (el *ErrorList) HasErrors() bool {
	for _, e := range el.errors {
		if e.IsError() {
			return true
		}
	}
	return false
}

// Len returns the number of errors in the list
func (el *ErrorList) Len() int {
	return len(el.errors)
}

// Errors returns all errors in the list
func (el *ErrorList) Errors() []*YapError {
	return el.errors
}

// Error implements the error interface
func (el *ErrorList) Error() string {
	if len(el.errors) == 0 {
		return ""
	}
	if len(el.errors) == 1 {
		return el.errors[0].Error()
	}

	var sb strings.Builder
	for i, err := range el.errors {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(err.Error())
	}
	return sb.String()
}

// FullError returns all errors with full context
func (el *ErrorList) FullError() string {
	if len(el.errors) == 0 {
		return ""
	}

	var sb strings.Builder
	for i, err := range el.errors {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(err.FullError())
	}
	return sb.String()
}

// Summary returns a summary of errors and warnings
func (el *ErrorList) Summary() string {
	errorCount := 0
	warningCount := 0

	for _, e := range el.errors {
		if e.IsError() {
			errorCount++
		} else if e.IsWarning() {
			warningCount++
		}
	}

	if errorCount == 0 && warningCount == 0 {
		return "no errors"
	}

	parts := []string{}
	if errorCount > 0 {
		parts = append(parts, fmt.Sprintf("%d error(s)", errorCount))
	}
	if warningCount > 0 {
		parts = append(parts, fmt.Sprintf("%d warning(s)", warningCount))
	}

	return strings.Join(parts, ", ")
}

// --- Convenience constructors for common errors ---

// Lexer error constructors

func NewInvalidCharError(file string, line, col int, char byte) *YapError {
	return &YapError{
		Code:     ErrInvalidCharacter,
		Severity: SeverityError,
		Phase:    PhaseLexer,
		Position: Position{File: file, Line: line, Column: col},
		Message:  fmt.Sprintf("invalid character %q", char),
	}
}

func NewUnterminatedStringError(file string, line, col int) *YapError {
	return &YapError{
		Code:     ErrUnterminatedString,
		Severity: SeverityError,
		Phase:    PhaseLexer,
		Position: Position{File: file, Line: line, Column: col},
		Message:  "unterminated string literal",
	}
}

func NewInvalidIndentError(file string, line, col int, got, expected int) *YapError {
	return &YapError{
		Code:     ErrInvalidIndentation,
		Severity: SeverityError,
		Phase:    PhaseLexer,
		Position: Position{File: file, Line: line, Column: col},
		Message:  fmt.Sprintf("invalid indentation: got %d spaces, expected %d", got, expected),
	}
}

func NewTabCharError(file string, line, col int) *YapError {
	return &YapError{
		Code:     ErrTabCharacter,
		Severity: SeverityError,
		Phase:    PhaseLexer,
		Position: Position{File: file, Line: line, Column: col},
		Message:  "tab character is not allowed, use spaces for indentation",
	}
}

// Parser error constructors

func NewUnexpectedTokenError(file string, line, col int, got, expected string) *YapError {
	return &YapError{
		Code:     ErrUnexpectedToken,
		Severity: SeverityError,
		Phase:    PhaseParser,
		Position: Position{File: file, Line: line, Column: col},
		Message:  fmt.Sprintf("unexpected token %q, expected %s", got, expected),
	}
}

func NewExpectedTokenError(file string, line, col int, expected string) *YapError {
	return &YapError{
		Code:     ErrExpectedToken,
		Severity: SeverityError,
		Phase:    PhaseParser,
		Position: Position{File: file, Line: line, Column: col},
		Message:  fmt.Sprintf("expected %s", expected),
	}
}

func NewUnknownStatementError(file string, line, col int, stmt string) *YapError {
	return &YapError{
		Code:     ErrUnknownStatement,
		Severity: SeverityError,
		Phase:    PhaseParser,
		Position: Position{File: file, Line: line, Column: col},
		Message:  fmt.Sprintf("unknown statement %q", stmt),
	}
}

// Builder/Semantic error constructors

func NewUnsupportedStatementError(stmtType string) *YapError {
	return &YapError{
		Code:     ErrUnsupportedStatement,
		Severity: SeverityError,
		Phase:    PhaseBuilder,
		Message:  fmt.Sprintf("unsupported statement type: %s", stmtType),
	}
}

func NewUndefinedVariableError(file string, line, col int, name string) *YapError {
	return &YapError{
		Code:     ErrUndefinedVariable,
		Severity: SeverityError,
		Phase:    PhaseBuilder,
		Position: Position{File: file, Line: line, Column: col},
		Message:  fmt.Sprintf("undefined variable %q", name),
	}
}

func NewTypeMismatchError(file string, line, col int, expected, got string) *YapError {
	return &YapError{
		Code:     ErrTypeMismatch,
		Severity: SeverityError,
		Phase:    PhaseBuilder,
		Position: Position{File: file, Line: line, Column: col},
		Message:  fmt.Sprintf("type mismatch: expected %s, got %s", expected, got),
	}
}

// Runtime error constructors

func NewUnknownOpcodeError(opcode int) *YapError {
	return &YapError{
		Code:     ErrUnknownOpcode,
		Severity: SeverityError,
		Phase:    PhaseRuntime,
		Message:  fmt.Sprintf("unknown opcode: %d", opcode),
	}
}

func NewDivisionByZeroError() *YapError {
	return &YapError{
		Code:     ErrDivisionByZero,
		Severity: SeverityError,
		Phase:    PhaseRuntime,
		Message:  "division by zero",
	}
}

func NewStackUnderflowError() *YapError {
	return &YapError{
		Code:     ErrStackUnderflow,
		Severity: SeverityError,
		Phase:    PhaseRuntime,
		Message:  "stack underflow",
	}
}

func NewOutOfBoundsError(index, length int) *YapError {
	return &YapError{
		Code:     ErrOutOfBounds,
		Severity: SeverityError,
		Phase:    PhaseRuntime,
		Message:  fmt.Sprintf("index out of bounds: %d (length: %d)", index, length),
	}
}

// Warning constructors

func NewDeprecatedWarning(file string, line, col int, feature, alternative string) *YapError {
	return &YapError{
		Code:     ErrInvalidSyntax, // Could add a dedicated warning code
		Severity: SeverityWarning,
		Phase:    PhaseParser,
		Position: Position{File: file, Line: line, Column: col},
		Message:  fmt.Sprintf("%q is deprecated, use %q instead", feature, alternative),
	}
}

func NewUnusedVariableWarning(file string, line, col int, name string) *YapError {
	return &YapError{
		Code:     ErrUndefinedVariable, // Could add a dedicated warning code
		Severity: SeverityWarning,
		Phase:    PhaseBuilder,
		Position: Position{File: file, Line: line, Column: col},
		Message:  fmt.Sprintf("unused variable %q", name),
	}
}
