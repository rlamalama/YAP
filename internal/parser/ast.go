package parser

// Program represents a YAMLScript file
// type Program struct {
// 	Statements []Stmt
// }

// Stmt is the interface for all statements
// type Stmt interface {
// 	stmt()
// 	Span() source.Span
// }

//
// // Concrete statements
//
// type PrintStmt struct {
// 	Value Expr
// 	span  source.Span
// }
//
// func (*PrintStmt) stmt()               {}
// func (s *PrintStmt) Span() source.Span { return s.span }
//
// type ExitStmt struct {
// 	Code int
// 	span source.Span
// }
//
// func (*ExitStmt) stmt()               {}
// func (s *ExitStmt) Span() source.Span { return s.span }
//
// type SetStmt struct {
// 	Name  string
// 	Value Expr
// 	span  source.Span
// }
//
// func (*SetStmt) stmt()               {}
// func (s *SetStmt) Span() source.Span { return s.span }
//
// type IfStmt struct {
// 	Condition Expr
// 	Then      []Stmt
// 	Else      []Stmt
// 	span      source.Span
// }
//
// func (*IfStmt) stmt()               {}
// func (s *IfStmt) Span() source.Span { return s.span }
//
// // Expressions are raw strings at this stage
// type Expr struct {
// 	Raw  string
// 	span source.Span
// }
