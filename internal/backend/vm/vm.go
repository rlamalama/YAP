package vm

import (
	"fmt"
	"strings"

	"github.com/rlamalama/YAP/internal/backend/ir"
	yaperror "github.com/rlamalama/YAP/internal/error"
)

type VM struct {
	instructions []ir.Instruction
	env          map[string]string
}

func New(instructions []ir.Instruction) *VM {
	env := make(map[string]string)
	return &VM{instructions: instructions, env: env}
}

func (vm *VM) Run() *yaperror.YapError {
	for _, instr := range vm.instructions {
		switch instr.Op {
		case ir.OpSet:
			parts := strings.SplitN(instr.Arg.Value, "=", 2)
			if len(parts) != 2 {
				return yaperror.NewInvalidSetIR(instr.Arg.Value)
			}
			vm.env[parts[0]] = parts[1]
		case ir.OpPrint:
			switch instr.Arg.Kind {
			case ir.OperandLiteral:
				fmt.Println(instr.Arg.Value)
			case ir.OperandIdentifier:
				val, ok := vm.env[instr.Arg.Value]
				if !ok {
					return yaperror.NewUndefinedVariable(instr.Arg.Value)
				}
				fmt.Println(val)
			}

		default:
			return yaperror.NewUnknownOpcodeError(int(instr.Op))
		}
	}
	return nil
}
