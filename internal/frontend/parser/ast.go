package parser

type Program struct {
	Statements []Stmt
}

type StmtType int

const (
	StmtTypeUnknown StmtType = iota
	StmtTypePrint
	StmtTypeSet
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
