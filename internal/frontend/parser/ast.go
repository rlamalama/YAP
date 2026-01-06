package parser

type Program struct {
	Statements []Stmt
}

type StmtType int

const (
	StmtTypeUnknown StmtType = iota
	StmtTypePrint
	StmtTypeSet
	StmtTypeIf
)

// Stmt is the interface for all statements
type Stmt interface {
	stmt()
	Type() StmtType
	// Span() source.Span
}

type PrintStmt struct {
	Expr Value
}

func (PrintStmt) stmt()          {}
func (PrintStmt) Type() StmtType { return StmtTypePrint }

type SetStmt struct {
	Assignment []*Assignment
}

func (s SetStmt) Type() StmtType { return StmtTypeSet }

func (SetStmt) stmt() {}

type Assignment struct {
	Name string
	Expr Value
}

// IfStmt represents an if-then-else statement
type IfStmt struct {
	Condition Value  // The conditional expression
	Then      []Stmt // Statements to execute if condition is true
	Else      []Stmt // Statements to execute if condition is false (can be nil)
}

func (IfStmt) stmt()          {}
func (IfStmt) Type() StmtType { return StmtTypeIf }
