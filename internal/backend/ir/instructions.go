package ir

type OpCode int

const (
	OpPrint OpCode = iota
)

type Instruction struct {
	Op  OpCode
	Arg string
}
