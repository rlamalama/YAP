package parser

type Stack []int

func NewStack() *Stack {
	return &Stack{}
}

func NewStackWithElem(val int) *Stack {
	return &Stack{val}
}

func (s *Stack) Length() int { return len(*s) }
func (s *Stack) Push(val int) {
	*s = append(*s, val)
}

func (s *Stack) Pop() (int, bool) {
	if len(*s) == 0 {
		return 0, false
	}

	val := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]

	return val, true
}

func (s *Stack) Peek() (int, bool) {
	if len(*s) == 0 {
		return 0, false
	}

	return (*s)[len(*s)-1], true
}
