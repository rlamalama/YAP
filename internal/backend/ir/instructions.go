package ir

type OpCode int

const (
	OpPrint OpCode = iota
	OpSet
)

type Instruction struct {
	Op   OpCode
	Arg  Operand
	Expr interface{} // Holds parser.Value for expression evaluation
}
