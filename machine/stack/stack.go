package stack

type Stack struct {
	index    int
	sizeMask int
	stack    []float64
}

func (s *Stack) Pop() float64 {
	s.index = (s.index - 1) & s.sizeMask
	return s.stack[s.index]
}

func (s *Stack) Push(val float64) {
	s.stack[s.index] = val
	s.index = (s.index + 1) & s.sizeMask
}

func NewStack(sizeMask int) *Stack {
	return &Stack{
		sizeMask: sizeMask,
		stack:    make([]float64, sizeMask+1),
	}
}
