package machine

import "github.com/boomlinde/gobassline/machine/stack"

func genFloatFunc(val float64) Instruction {
	return func(s *stack.Stack) { s.Push(val) }
}
