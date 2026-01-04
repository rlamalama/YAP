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
		{Op: ir.OpPrint, Arg: arg},
	})

	output := test_util.CaptureStdout(t, func() {

		require.NoError(t, vm.Run())
	})

	assert.Equal(t, arg+"\n", output)
}
