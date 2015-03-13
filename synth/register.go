package synth

import (
	"github.com/boomlinde/acidforth/collection"
	"github.com/boomlinde/acidforth/machine/stack"
)

type Register struct {
	Value float64
}

func NewRegister(name string, c *collection.Collection) *Register {
	a := &Register{}

	c.Machine.Register(">"+name, func(s *stack.Stack) {
		a.Value = s.Pop()
	})

	c.Machine.Register(name+">", func(s *stack.Stack) {
		s.Push(a.Value)
		a.Value = 0
	})
	return a
}
