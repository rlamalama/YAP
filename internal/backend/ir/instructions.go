package ir

type OpCode int

const (
	OpPrint OpCode = iota
	OpSet
	OpJumpIfFalse // Jump to Arg.Offset if Expr evaluates to false
	OpJump        // Unconditional jump to Arg.Offset
)

type Instruction struct {
	Op   OpCode
	Arg  Operand
	Expr interface{} // Holds parser.Value for expression evaluation
}
