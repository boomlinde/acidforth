package synth

import (
	"github.com/boomlinde/acidforth/collection"
	"github.com/boomlinde/acidforth/machine/stack"
)

type Accumulator struct {
	total float64
}

func NewAccumulator(name string, c *collection.Collection) *Accumulator {
	a := &Accumulator{}

	c.Machine.Register(">"+name, func(s *stack.Stack) {
		a.total += s.Pop()
	})

	c.Machine.Register(name+">", func(s *stack.Stack) {
		s.Push(a.total)
		a.total = 0
	})
	return a
}
