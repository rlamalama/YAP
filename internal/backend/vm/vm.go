package vm

import (
	"fmt"

	"github.com/rlamalama/YAP/internal/backend/ir"
)

type VM struct {
	instructions []ir.Instruction
}

func New(instructions []ir.Instruction) *VM {
	return &VM{instructions: instructions}
}

func (vm *VM) Run() error {
	for _, instr := range vm.instructions {
		switch instr.Op {
		case ir.OpPrint:
			fmt.Println(instr.Arg)
		default:
			return fmt.Errorf("unknown opcode %v", instr.Op)
		}
	}
	return nil
}
