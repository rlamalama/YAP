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

func TestVMBooleanLiteralTrue(t *testing.T) {
	// Test: flag = True, print flag (should output true)
	v := vm.New([]ir.Instruction{
		{
			Op:   ir.OpSet,
			Arg:  ir.Operand{Kind: ir.OperandIdentifier, Value: "flag"},
			Expr: &parser.BooleanLiteral{Value: true},
		},
		{Op: ir.OpPrint, Expr: &parser.Identifier{Name: "flag"}},
	})

	output := test_util.CaptureStdout(t, func() {
		require.Nil(t, v.Run())
	})

	assert.Equal(t, "true\n", output)
}

func TestVMBooleanLiteralFalse(t *testing.T) {
	// Test: flag = False, print flag (should output false)
	v := vm.New([]ir.Instruction{
		{
			Op:   ir.OpSet,
			Arg:  ir.Operand{Kind: ir.OperandIdentifier, Value: "flag"},
			Expr: &parser.BooleanLiteral{Value: false},
		},
		{Op: ir.OpPrint, Expr: &parser.Identifier{Name: "flag"}},
	})

	output := test_util.CaptureStdout(t, func() {
		require.Nil(t, v.Run())
	})

	assert.Equal(t, "false\n", output)
}

func TestVMComparisonGreaterThan(t *testing.T) {
	// Test: a = 10, b = 5, isGreater = a > b (should output true)
	v := vm.New([]ir.Instruction{
		{
			Op:   ir.OpSet,
			Arg:  ir.Operand{Kind: ir.OperandIdentifier, Value: "a"},
			Expr: &parser.NumericLiteral{Value: 10},
		},
		{
			Op:   ir.OpSet,
			Arg:  ir.Operand{Kind: ir.OperandIdentifier, Value: "b"},
			Expr: &parser.NumericLiteral{Value: 5},
		},
		{
			Op:  ir.OpSet,
			Arg: ir.Operand{Kind: ir.OperandIdentifier, Value: "isGreater"},
			Expr: &parser.BinaryExpr{
				Left:     &parser.Identifier{Name: "a"},
				Operator: ">",
				Right:    &parser.Identifier{Name: "b"},
			},
		},
		{Op: ir.OpPrint, Expr: &parser.Identifier{Name: "isGreater"}},
	})

	output := test_util.CaptureStdout(t, func() {
		require.Nil(t, v.Run())
	})

	assert.Equal(t, "true\n", output)
}

func TestVMComparisonLessThan(t *testing.T) {
	// Test: a = 5, b = 10, isLess = a < b (should output true)
	v := vm.New([]ir.Instruction{
		{
			Op:   ir.OpSet,
			Arg:  ir.Operand{Kind: ir.OperandIdentifier, Value: "a"},
			Expr: &parser.NumericLiteral{Value: 5},
		},
		{
			Op:   ir.OpSet,
			Arg:  ir.Operand{Kind: ir.OperandIdentifier, Value: "b"},
			Expr: &parser.NumericLiteral{Value: 10},
		},
		{
			Op:  ir.OpSet,
			Arg: ir.Operand{Kind: ir.OperandIdentifier, Value: "isLess"},
			Expr: &parser.BinaryExpr{
				Left:     &parser.Identifier{Name: "a"},
				Operator: "<",
				Right:    &parser.Identifier{Name: "b"},
			},
		},
		{Op: ir.OpPrint, Expr: &parser.Identifier{Name: "isLess"}},
	})

	output := test_util.CaptureStdout(t, func() {
		require.Nil(t, v.Run())
	})

	assert.Equal(t, "true\n", output)
}

func TestVMComparisonEqual(t *testing.T) {
	// Test: a = 5, b = 5, isEqual = a == b (should output true)
	v := vm.New([]ir.Instruction{
		{
			Op:   ir.OpSet,
			Arg:  ir.Operand{Kind: ir.OperandIdentifier, Value: "a"},
			Expr: &parser.NumericLiteral{Value: 5},
		},
		{
			Op:   ir.OpSet,
			Arg:  ir.Operand{Kind: ir.OperandIdentifier, Value: "b"},
			Expr: &parser.NumericLiteral{Value: 5},
		},
		{
			Op:  ir.OpSet,
			Arg: ir.Operand{Kind: ir.OperandIdentifier, Value: "isEqual"},
			Expr: &parser.BinaryExpr{
				Left:     &parser.Identifier{Name: "a"},
				Operator: "==",
				Right:    &parser.Identifier{Name: "b"},
			},
		},
		{Op: ir.OpPrint, Expr: &parser.Identifier{Name: "isEqual"}},
	})

	output := test_util.CaptureStdout(t, func() {
		require.Nil(t, v.Run())
	})

	assert.Equal(t, "true\n", output)
}

func TestVMComparisonNotEqual(t *testing.T) {
	// Test: a = 10, b = 5, notEqual = a != b (should output true)
	v := vm.New([]ir.Instruction{
		{
			Op:   ir.OpSet,
			Arg:  ir.Operand{Kind: ir.OperandIdentifier, Value: "a"},
			Expr: &parser.NumericLiteral{Value: 10},
		},
		{
			Op:   ir.OpSet,
			Arg:  ir.Operand{Kind: ir.OperandIdentifier, Value: "b"},
			Expr: &parser.NumericLiteral{Value: 5},
		},
		{
			Op:  ir.OpSet,
			Arg: ir.Operand{Kind: ir.OperandIdentifier, Value: "notEqual"},
			Expr: &parser.BinaryExpr{
				Left:     &parser.Identifier{Name: "a"},
				Operator: "!=",
				Right:    &parser.Identifier{Name: "b"},
			},
		},
		{Op: ir.OpPrint, Expr: &parser.Identifier{Name: "notEqual"}},
	})

	output := test_util.CaptureStdout(t, func() {
		require.Nil(t, v.Run())
	})

	assert.Equal(t, "true\n", output)
}

func TestVMComparisonGreaterOrEqual(t *testing.T) {
	// Test: a = 10, print a >= 10 (should output true)
	v := vm.New([]ir.Instruction{
		{
			Op:   ir.OpSet,
			Arg:  ir.Operand{Kind: ir.OperandIdentifier, Value: "a"},
			Expr: &parser.NumericLiteral{Value: 10},
		},
		{
			Op: ir.OpPrint,
			Expr: &parser.BinaryExpr{
				Left:     &parser.Identifier{Name: "a"},
				Operator: ">=",
				Right:    &parser.NumericLiteral{Value: 10},
			},
		},
	})

	output := test_util.CaptureStdout(t, func() {
		require.Nil(t, v.Run())
	})

	assert.Equal(t, "true\n", output)
}

func TestVMComparisonLessOrEqual(t *testing.T) {
	// Test: b = 5, a = 10, isLessOrEqual = b <= a (should output true)
	v := vm.New([]ir.Instruction{
		{
			Op:   ir.OpSet,
			Arg:  ir.Operand{Kind: ir.OperandIdentifier, Value: "b"},
			Expr: &parser.NumericLiteral{Value: 5},
		},
		{
			Op:   ir.OpSet,
			Arg:  ir.Operand{Kind: ir.OperandIdentifier, Value: "a"},
			Expr: &parser.NumericLiteral{Value: 10},
		},
		{
			Op:  ir.OpSet,
			Arg: ir.Operand{Kind: ir.OperandIdentifier, Value: "isLessOrEqual"},
			Expr: &parser.BinaryExpr{
				Left:     &parser.Identifier{Name: "b"},
				Operator: "<=",
				Right:    &parser.Identifier{Name: "a"},
			},
		},
		{Op: ir.OpPrint, Expr: &parser.Identifier{Name: "isLessOrEqual"}},
	})

	output := test_util.CaptureStdout(t, func() {
		require.Nil(t, v.Run())
	})

	assert.Equal(t, "true\n", output)
}

func TestVMComparisonFalseResult(t *testing.T) {
	// Test: a = 10, b = 5, isEqual = a == b (should output false)
	v := vm.New([]ir.Instruction{
		{
			Op:   ir.OpSet,
			Arg:  ir.Operand{Kind: ir.OperandIdentifier, Value: "a"},
			Expr: &parser.NumericLiteral{Value: 10},
		},
		{
			Op:   ir.OpSet,
			Arg:  ir.Operand{Kind: ir.OperandIdentifier, Value: "b"},
			Expr: &parser.NumericLiteral{Value: 5},
		},
		{
			Op:  ir.OpSet,
			Arg: ir.Operand{Kind: ir.OperandIdentifier, Value: "isEqual"},
			Expr: &parser.BinaryExpr{
				Left:     &parser.Identifier{Name: "a"},
				Operator: "==",
				Right:    &parser.Identifier{Name: "b"},
			},
		},
		{Op: ir.OpPrint, Expr: &parser.Identifier{Name: "isEqual"}},
	})

	output := test_util.CaptureStdout(t, func() {
		require.Nil(t, v.Run())
	})

	assert.Equal(t, "false\n", output)
}

func TestVMBooleanComparison(t *testing.T) {
	// Test: flag1 = True, flag2 = True, areEqual = flag1 == flag2 (should output true)
	v := vm.New([]ir.Instruction{
		{
			Op:   ir.OpSet,
			Arg:  ir.Operand{Kind: ir.OperandIdentifier, Value: "flag1"},
			Expr: &parser.BooleanLiteral{Value: true},
		},
		{
			Op:   ir.OpSet,
			Arg:  ir.Operand{Kind: ir.OperandIdentifier, Value: "flag2"},
			Expr: &parser.BooleanLiteral{Value: true},
		},
		{
			Op:  ir.OpSet,
			Arg: ir.Operand{Kind: ir.OperandIdentifier, Value: "areEqual"},
			Expr: &parser.BinaryExpr{
				Left:     &parser.Identifier{Name: "flag1"},
				Operator: "==",
				Right:    &parser.Identifier{Name: "flag2"},
			},
		},
		{Op: ir.OpPrint, Expr: &parser.Identifier{Name: "areEqual"}},
	})

	output := test_util.CaptureStdout(t, func() {
		require.Nil(t, v.Run())
	})

	assert.Equal(t, "true\n", output)
}

func TestVMStringComparison(t *testing.T) {
	// Test: s1 = "hello", s2 = "hello", areEqual = s1 == s2 (should output true)
	v := vm.New([]ir.Instruction{
		{
			Op:   ir.OpSet,
			Arg:  ir.Operand{Kind: ir.OperandIdentifier, Value: "s1"},
			Expr: &parser.StringLiteral{Value: "hello"},
		},
		{
			Op:   ir.OpSet,
			Arg:  ir.Operand{Kind: ir.OperandIdentifier, Value: "s2"},
			Expr: &parser.StringLiteral{Value: "hello"},
		},
		{
			Op:  ir.OpSet,
			Arg: ir.Operand{Kind: ir.OperandIdentifier, Value: "areEqual"},
			Expr: &parser.BinaryExpr{
				Left:     &parser.Identifier{Name: "s1"},
				Operator: "==",
				Right:    &parser.Identifier{Name: "s2"},
			},
		},
		{Op: ir.OpPrint, Expr: &parser.Identifier{Name: "areEqual"}},
	})

	output := test_util.CaptureStdout(t, func() {
		require.Nil(t, v.Run())
	})

	assert.Equal(t, "true\n", output)
}

func TestVMUndefinedVariableError(t *testing.T) {
	// Simulates the scenario from 0008-comments-ignore-in-block.yap:
	// - set:
	//   - x: 10
	//     // - y: 5  (commented out, y is not defined)
	// - print: x
	// - print: y  (should error because y is undefined)
	v := vm.New([]ir.Instruction{
		{
			Op:   ir.OpSet,
			Arg:  ir.Operand{Kind: ir.OperandIdentifier, Value: "x"},
			Expr: &parser.NumericLiteral{Value: 10},
		},
		// y is NOT set (simulating it being commented out)
		{Op: ir.OpPrint, Expr: &parser.Identifier{Name: "x"}},
		{Op: ir.OpPrint, Expr: &parser.Identifier{Name: "y"}}, // y is undefined
	})

	err := v.Run()

	// Should error because y is undefined
	assert.NotNil(t, err, "should error because y is undefined")
	assert.Contains(t, err.Error(), "y", "error should mention undefined variable y")
}
