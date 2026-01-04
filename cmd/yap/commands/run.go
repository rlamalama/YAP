package commands

import (
	"log"
	"os"

	"github.com/rlamalama/YAP/internal/backend/build"
	"github.com/rlamalama/YAP/internal/backend/vm"
	"github.com/rlamalama/YAP/internal/frontend/parser"
)

func RunCmd(args []string) {
	// Accessing flags
	file := args[0]

	_, err := os.Stat(file)
	if err != nil {
		log.Fatalf("error finding file: %+v", err)
	}

	parser := parser.NewParser(file)
	ast, err := parser.Parse()
	if err != nil {
		log.Fatalf("error parsing program: %+v", err)
	}

	builder := build.New()
	program, err := builder.Build(ast.Statements)

	if err != nil {
		log.Fatalf("error building program: %+v", err)
	}

	vm := vm.New(program)
	if err := vm.Run(); err != nil {
		log.Fatalf("error running program: %+v", err)
	}
}
