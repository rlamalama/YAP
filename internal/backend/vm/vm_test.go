package vm_test

import (
	"testing"

	"github.com/rlamalama/YAP/internal/backend/ir"
	"github.com/rlamalama/YAP/internal/backend/vm"
	"github.com/rlamalama/YAP/internal/frontend/parser"
	test_util "github.com/rlamalama/YAP/test/test-util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVMPrint(t *testing.T) {
	arg := "hi"
	vm := vm.New([]ir.Instruction{
		{Op: ir.OpPrint, Expr: &parser.StringLiteral{Value: arg}},
	})

	output := test_util.CaptureStdout(t, func() {
		require.Nil(t, vm.Run())
	})

	assert.Equal(t, arg+"\n", output)
}

func TestVMSetAndPrint(t *testing.T) {
	key, arg := "x", "hi"
	vm := vm.New([]ir.Instruction{
		{
			Op:   ir.OpSet,
			Arg:  ir.Operand{Kind: ir.OperandIdentifier, Value: key},
			Expr: &parser.StringLiteral{Value: arg},
		},
		{Op: ir.OpPrint, Expr: &parser.Identifier{Name: key}},
	})

	output := test_util.CaptureStdout(t, func() {
		require.Nil(t, vm.Run())
	})

	assert.Equal(t, arg+"\n", output)
}

func TestVMBinaryExprAddition(t *testing.T) {
	// Test: x = 10 + 5, print x (should output 15)
	binExpr := &parser.BinaryExpr{
		Left:     &parser.NumericLiteral{Value: 10},
		Operator: "+",
		Right:    &parser.NumericLiteral{Value: 5},
	}
	v := vm.New([]ir.Instruction{
		{
			Op:   ir.OpSet,
			Arg:  ir.Operand{Kind: ir.OperandIdentifier, Value: "x"},
			Expr: binExpr,
		},
		{Op: ir.OpPrint, Expr: &parser.Identifier{Name: "x"}},
	})

	output := test_util.CaptureStdout(t, func() {
		require.Nil(t, v.Run())
	})

	assert.Equal(t, "15\n", output)
}

func TestVMBinaryExprChained(t *testing.T) {
	// Test: x = 10 + 10 - 15 (should output 5)
	innerExpr := &parser.BinaryExpr{
		Left:     &parser.NumericLiteral{Value: 10},
		Operator: "+",
		Right:    &parser.NumericLiteral{Value: 10},
	}
	outerExpr := &parser.BinaryExpr{
		Left:     innerExpr,
		Operator: "-",
		Right:    &parser.NumericLiteral{Value: 15},
	}

	v := vm.New([]ir.Instruction{
		{
			Op:   ir.OpSet,
			Arg:  ir.Operand{Kind: ir.OperandIdentifier, Value: "x"},
			Expr: outerExpr,
		},
		{Op: ir.OpPrint, Expr: &parser.Identifier{Name: "x"}},
	})

	output := test_util.CaptureStdout(t, func() {
		require.Nil(t, v.Run())
	})

	assert.Equal(t, "5\n", output)
}

func TestVMBinaryExprWithVariables(t *testing.T) {
	// Test: x = 5, y = x * 4 (should output 20)
	v := vm.New([]ir.Instruction{
		{
			Op:   ir.OpSet,
			Arg:  ir.Operand{Kind: ir.OperandIdentifier, Value: "x"},
			Expr: &parser.NumericLiteral{Value: 5},
		},
		{
			Op:  ir.OpSet,
			Arg: ir.Operand{Kind: ir.OperandIdentifier, Value: "y"},
			Expr: &parser.BinaryExpr{
				Left:     &parser.Identifier{Name: "x"},
				Operator: "*",
				Right:    &parser.NumericLiteral{Value: 4},
			},
		},
		{Op: ir.OpPrint, Expr: &parser.Identifier{Name: "y"}},
	})

	output := test_util.CaptureStdout(t, func() {
		require.Nil(t, v.Run())
	})

	assert.Equal(t, "20\n", output)
}

func TestVMBinaryExprDivision(t *testing.T) {
	// Test: y = 20, z = y / 5 (should output 4)
	v := vm.New([]ir.Instruction{
		{
			Op:   ir.OpSet,
			Arg:  ir.Operand{Kind: ir.OperandIdentifier, Value: "y"},
			Expr: &parser.NumericLiteral{Value: 20},
		},
		{
			Op:  ir.OpSet,
			Arg: ir.Operand{Kind: ir.OperandIdentifier, Value: "z"},
			Expr: &parser.BinaryExpr{
				Left:     &parser.Identifier{Name: "y"},
				Operator: "/",
				Right:    &parser.NumericLiteral{Value: 5},
			},
		},
		{Op: ir.OpPrint, Expr: &parser.Identifier{Name: "z"}},
	})

	output := test_util.CaptureStdout(t, func() {
		require.Nil(t, v.Run())
	})

	assert.Equal(t, "4\n", output)
}

func TestVMBinaryExprStringConcat(t *testing.T) {
	// Test: "hello" + " " + "world!" (should output "hello world!")
	innerExpr := &parser.BinaryExpr{
		Left:     &parser.StringLiteral{Value: "hello"},
		Operator: "+",
		Right:    &parser.StringLiteral{Value: " "},
	}
	outerExpr := &parser.BinaryExpr{
		Left:     innerExpr,
		Operator: "+",
		Right:    &parser.StringLiteral{Value: "world!"},
	}

	v := vm.New([]ir.Instruction{
		{Op: ir.OpPrint, Expr: outerExpr},
	})

	output := test_util.CaptureStdout(t, func() {
		require.Nil(t, v.Run())
	})

	assert.Equal(t, "hello world!\n", output)
}

func TestVMPrintBinaryExprWithVariables(t *testing.T) {
	// Test: x = 5, z = 4, print x * z (should output 20)
	v := vm.New([]ir.Instruction{
		{
			Op:   ir.OpSet,
			Arg:  ir.Operand{Kind: ir.OperandIdentifier, Value: "x"},
			Expr: &parser.NumericLiteral{Value: 5},
		},
		{
			Op:   ir.OpSet,
			Arg:  ir.Operand{Kind: ir.OperandIdentifier, Value: "z"},
			Expr: &parser.NumericLiteral{Value: 4},
		},
		{
			Op: ir.OpPrint,
			Expr: &parser.BinaryExpr{
				Left:     &parser.Identifier{Name: "x"},
				Operator: "*",
				Right:    &parser.Identifier{Name: "z"},
			},
		},
	})

	output := test_util.CaptureStdout(t, func() {
		require.Nil(t, v.Run())
	})

	assert.Equal(t, "20\n", output)
}
