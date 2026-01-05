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

	default:
		return fmt.Errorf("unsupported statement %T", stmt)
	}
	return nil
}
