package build

import (
	"fmt"

	"github.com/rlamalama/YAP/internal/backend/ir"
	"github.com/rlamalama/YAP/internal/frontend/parser"
)

type Builder struct {
	instructions []ir.Instruction
}

func New() *Builder {
	return &Builder{}
}

func (b *Builder) Build(stmts []parser.Stmt) ([]ir.Instruction, error) {
	for _, stmt := range stmts {
		if err := b.buildStmt(stmt); err != nil {
			return nil, err
		}
	}
	return b.instructions, nil
}

func (b *Builder) buildStmt(stmt parser.Stmt) error {
	switch s := stmt.(type) {
	case parser.PrintStmt:
		b.instructions = append(b.instructions, ir.Instruction{
			Op:   ir.OpPrint,
			Expr: s.Expr,
		})

	case parser.SetStmt:
		for _, assignment := range s.Assignment {
			b.instructions = append(b.instructions, ir.Instruction{
				Op: ir.OpSet,
				Arg: ir.Operand{
					Kind:  ir.OperandIdentifier,
					Value: assignment.Name,
				},
				Expr: assignment.Expr,
			})
		}

	case parser.IfStmt:
		if err := b.buildIfStmt(s); err != nil {
			return err
		}

	default:
		return fmt.Errorf("unsupported statement %T", stmt)
	}
	return nil
}

func (b *Builder) buildIfStmt(s parser.IfStmt) error {
	// Emit jump-if-false instruction with placeholder offset
	jumpIfFalseIdx := len(b.instructions)
	b.instructions = append(b.instructions, ir.Instruction{
		Op:   ir.OpJumpIfFalse,
		Arg:  ir.Operand{Kind: ir.OperandOffset, Offset: 0}, // placeholder
		Expr: s.Condition,
	})

	// Build the "then" block
	for _, stmt := range s.Then {
		if err := b.buildStmt(stmt); err != nil {
			return err
		}
	}

	if len(s.Else) > 0 {
		// If there's an else block, we need to jump over it after the then block
		jumpOverElseIdx := len(b.instructions)
		b.instructions = append(b.instructions, ir.Instruction{
			Op:  ir.OpJump,
			Arg: ir.Operand{Kind: ir.OperandOffset, Offset: 0}, // placeholder
		})

		// Patch the jump-if-false to jump to the else block
		b.instructions[jumpIfFalseIdx].Arg.Offset = len(b.instructions)

		// Build the "else" block
		for _, stmt := range s.Else {
			if err := b.buildStmt(stmt); err != nil {
				return err
			}
		}

		// Patch the jump-over-else to jump to after the else block
		b.instructions[jumpOverElseIdx].Arg.Offset = len(b.instructions)
	} else {
		// No else block, just patch jump-if-false to after the then block
		b.instructions[jumpIfFalseIdx].Arg.Offset = len(b.instructions)
	}

	return nil
}
