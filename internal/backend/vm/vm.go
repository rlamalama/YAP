package vm

import (
	"fmt"

	"github.com/rlamalama/YAP/internal/backend/ir"
	yaperror "github.com/rlamalama/YAP/internal/error"
	"github.com/rlamalama/YAP/internal/frontend/lexer"
	"github.com/rlamalama/YAP/internal/frontend/parser"
)

type VM struct {
	instructions []ir.Instruction
	env          map[string]interface{}
}

func New(instructions []ir.Instruction) *VM {
	env := make(map[string]interface{})
	return &VM{instructions: instructions, env: env}
}

func (vm *VM) Run() *yaperror.YapError {
	for _, instr := range vm.instructions {
		switch instr.Op {
		case ir.OpSet:
			name := instr.Arg.Value
			val, err := vm.evaluate(instr.Expr)
			if err != nil {
				return err
			}
			vm.env[name] = val

		case ir.OpPrint:
			val, err := vm.evaluate(instr.Expr)
			if err != nil {
				return err
			}
			fmt.Println(val)

		default:
			return yaperror.NewUnknownOpcodeError(int(instr.Op))
		}
	}
	return nil
}

func (vm *VM) evaluate(expr interface{}) (interface{}, *yaperror.YapError) {
	switch v := expr.(type) {
	case *parser.NumericLiteral:
		return v.Value, nil

	case *parser.StringLiteral:
		return v.Value, nil

	case *parser.Identifier:
		val, ok := vm.env[v.Name]
		if !ok {
			return nil, yaperror.NewUndefinedVariable(v.Name)
		}
		return val, nil

	case *parser.BinaryExpr:
		left, err := vm.evaluate(v.Left)
		if err != nil {
			return nil, err
		}
		right, err := vm.evaluate(v.Right)
		if err != nil {
			return nil, err
		}

		// Handle numeric operations
		leftInt, leftIsInt := left.(int)
		rightInt, rightIsInt := right.(int)

		if leftIsInt && rightIsInt {
			switch v.Operator {
			case lexer.ArithmeticAdditionOperator.String():
				return leftInt + rightInt, nil
			case lexer.ArithmeticSubtractionOperator.String():
				return leftInt - rightInt, nil
			case lexer.ArithmeticMultiplicationOperator.String():
				return leftInt * rightInt, nil
			case lexer.ArithmeticDivisionOperator.String():
				if rightInt == 0 {
					return nil, yaperror.NewRuntimeError("division by zero")
				}
				return leftInt / rightInt, nil
			}
		}

		// Handle string concatenation
		leftStr, leftIsStr := left.(string)
		rightStr, rightIsStr := right.(string)

		if leftIsStr && rightIsStr && v.Operator == lexer.ArithmeticAdditionOperator.String() {
			return leftStr + rightStr, nil
		}

		return nil, yaperror.NewRuntimeError(fmt.Sprintf("unsupported operation: %T %s %T", left, v.Operator, right))

	default:
		return nil, yaperror.NewRuntimeError(fmt.Sprintf("unknown expression type: %T", expr))
	}
}
