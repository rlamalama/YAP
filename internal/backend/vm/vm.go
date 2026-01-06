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
	pc           int // program counter
}

func New(instructions []ir.Instruction) *VM {
	env := make(map[string]interface{})
	return &VM{instructions: instructions, env: env, pc: 0}
}

func (vm *VM) Run() *yaperror.YapError {
	for vm.pc < len(vm.instructions) {
		instr := vm.instructions[vm.pc]
		switch instr.Op {
		case ir.OpSet:
			name := instr.Arg.Value
			val, err := vm.evaluate(instr.Expr)
			if err != nil {
				return err
			}
			vm.env[name] = val
			vm.pc++

		case ir.OpPrint:
			val, err := vm.evaluate(instr.Expr)
			if err != nil {
				return err
			}
			fmt.Println(val)
			vm.pc++

		case ir.OpJumpIfFalse:
			val, err := vm.evaluate(instr.Expr)
			if err != nil {
				return err
			}
			boolVal, ok := val.(bool)
			if !ok {
				return yaperror.NewRuntimeError(fmt.Sprintf("condition must be a boolean, got %T", val))
			}
			if !boolVal {
				vm.pc = instr.Arg.Offset
			} else {
				vm.pc++
			}

		case ir.OpJump:
			vm.pc = instr.Arg.Offset

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

	case *parser.BooleanLiteral:
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

		// Handle numeric operations and comparisons
		leftInt, leftIsInt := left.(int)
		rightInt, rightIsInt := right.(int)

		if leftIsInt && rightIsInt {
			switch v.Operator {
			// Arithmetic operators
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
			// Comparison operators for integers
			case lexer.ComparisonGtOperator.String():
				return leftInt > rightInt, nil
			case lexer.ComparisonLtOperator.String():
				return leftInt < rightInt, nil
			case lexer.ComparisonGteOperator.String():
				return leftInt >= rightInt, nil
			case lexer.ComparisonLteOperator.String():
				return leftInt <= rightInt, nil
			case lexer.ComparisonEqOperator.String():
				return leftInt == rightInt, nil
			case lexer.ComparisonNeOperator.String():
				return leftInt != rightInt, nil
			}
		}

		// Handle string operations
		leftStr, leftIsStr := left.(string)
		rightStr, rightIsStr := right.(string)

		if leftIsStr && rightIsStr {
			switch v.Operator {
			case lexer.ArithmeticAdditionOperator.String():
				return leftStr + rightStr, nil
			// Comparison operators for strings
			case lexer.ComparisonEqOperator.String():
				return leftStr == rightStr, nil
			case lexer.ComparisonNeOperator.String():
				return leftStr != rightStr, nil
			case lexer.ComparisonGtOperator.String():
				return leftStr > rightStr, nil
			case lexer.ComparisonLtOperator.String():
				return leftStr < rightStr, nil
			case lexer.ComparisonGteOperator.String():
				return leftStr >= rightStr, nil
			case lexer.ComparisonLteOperator.String():
				return leftStr <= rightStr, nil
			}
		}

		// Handle boolean operations
		leftBool, leftIsBool := left.(bool)
		rightBool, rightIsBool := right.(bool)

		if leftIsBool && rightIsBool {
			switch v.Operator {
			case lexer.ComparisonEqOperator.String():
				return leftBool == rightBool, nil
			case lexer.ComparisonNeOperator.String():
				return leftBool != rightBool, nil
			}
		}

		return nil, yaperror.NewRuntimeError(fmt.Sprintf("unsupported operation: %T %s %T", left, v.Operator, right))

	default:
		return nil, yaperror.NewRuntimeError(fmt.Sprintf("unknown expression type: %T", expr))
	}
}
