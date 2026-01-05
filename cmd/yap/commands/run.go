package commands

import (
	"log"
	"os"
	"strings"

	"github.com/rlamalama/YAP/internal/backend/build"
	"github.com/rlamalama/YAP/internal/backend/vm"
	"github.com/rlamalama/YAP/internal/frontend/parser"
)

const FileExtYAP = ".yap"

func RunCmd(args []string) {
	// Accessing flags
	file := args[0]

	_, err := os.Stat(file)
	if err != nil {
		log.Fatalf("error finding file: %+v", err)
	}
	if !strings.HasSuffix(strings.ToLower(file), FileExtYAP) {
		log.Fatalf("file %s must be a .yap or .YAP file", file)
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
