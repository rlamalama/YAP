package vm_test

import (
	"testing"

	"github.com/rlamalama/YAP/internal/backend/ir"
	"github.com/rlamalama/YAP/internal/backend/vm"
	test_util "github.com/rlamalama/YAP/test/test-util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVMPrint(t *testing.T) {
	arg := "hi"
	vm := vm.New([]ir.Instruction{
		{Op: ir.OpPrint, Arg: ir.Operand{
			Kind:  ir.OperandLiteral,
			Value: arg,
		}},
	})

	output := test_util.CaptureStdout(t, func() {
		require.Nil(t, vm.Run())
	})

	assert.Equal(t, arg+"\n", output)
}

func TestVMSetAndPrint(t *testing.T) {
	key, arg := "x", "hi"
	vm := vm.New([]ir.Instruction{
		{Op: ir.OpSet, Arg: ir.Operand{
			Kind:  ir.OperandIdentifier,
			Value: key + "=" + arg,
		}},
		{Op: ir.OpPrint, Arg: ir.Operand{
			Kind:  ir.OperandIdentifier,
			Value: key,
		}},
	})

	output := test_util.CaptureStdout(t, func() {
		require.Nil(t, vm.Run())
	})

	assert.Equal(t, arg+"\n", output)
}
