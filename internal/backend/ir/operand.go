package ir

type OperandKind int

const (
	OperandLiteral OperandKind = iota
	OperandIdentifier
)

type Operand struct {
	Kind  OperandKind
	Value string
}
