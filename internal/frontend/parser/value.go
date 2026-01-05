package parser

import "fmt"

type Value interface {
	value()
	String() string
}

type StringLiteral struct {
	Value string
}

func (*StringLiteral) value()           {}
func (s *StringLiteral) String() string { return s.Value }

// Numbers only support int right now
type NumericLiteral struct {
	Value int
}

func (*NumericLiteral) value()           {}
func (n *NumericLiteral) String() string { return fmt.Sprintf("%d", n.Value) }

type Identifier struct {
	Name string
}

func (*Identifier) value()           {}
func (i *Identifier) String() string { return i.Name }
