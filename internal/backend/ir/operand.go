package ir

type OperandKind int

const (
	OperandLiteral OperandKind = iota
	OperandIdentifier
	OperandOffset
)

type Operand struct {
	Kind   OperandKind
	Value  string
	Offset int // Used for jump targets
}
