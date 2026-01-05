package lexer

type OperatorType string

const (
	OperatorComparison = "comparison"
	OperatorArithmetic = "arithmetic"
)

type Operator interface {
	operator()
	String() string
}

type ComparisonOperator string

func (ComparisonOperator) operator()        {}
func (c ComparisonOperator) String() string { return string(c) }

const (
	ComparisonGtOperator  ComparisonOperator = ">"
	ComparisonLtOperator  ComparisonOperator = "<"
	ComparisonGteOperator ComparisonOperator = ">="
	ComparisonLteOperator ComparisonOperator = "<="
	ComparisonEqOperator  ComparisonOperator = "=="
	ComparisonNeOperator  ComparisonOperator = "!="
)

func StartsWithOperator(c byte) bool {
	return IsArithmeticOperator(c) || IsComparisonOperator(string(c))
}

func IsComparisonOperator(c string) bool {
	switch c {
	case ComparisonGtOperator.String(), ComparisonLtOperator.String(), ComparisonGteOperator.String(),
		ComparisonLteOperator.String(), ComparisonEqOperator.String(), ComparisonNeOperator.String():
		return true
	default:
		return false
	}
}

type ArithmeticOperator byte

func (ArithmeticOperator) operator()        {}
func (a ArithmeticOperator) String() string { return string(a) }
func (a ArithmeticOperator) Byte() byte     { return byte(a) }

const (
	ArithmeticAdditionOperator       ArithmeticOperator = '+'
	ArithmeticSubtractionOperator    ArithmeticOperator = '-'
	ArithmeticMultiplicationOperator ArithmeticOperator = '*'
	ArithmeticDivisionOperator       ArithmeticOperator = '/'
)

func IsArithmeticOperator(c byte) bool {
	switch c {
	case ArithmeticAdditionOperator.Byte(), ArithmeticSubtractionOperator.Byte(),
		ArithmeticMultiplicationOperator.Byte(), ArithmeticDivisionOperator.Byte():
		return true
	default:
		return false
	}
}
