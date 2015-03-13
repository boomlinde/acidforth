package synth

import (
	"github.com/boomlinde/acidforth/collection"
	"github.com/boomlinde/acidforth/machine/stack"
)

type Accumulator struct {
	Total float64
}

func NewAccumulator(name string, c *collection.Collection) *Accumulator {
	a := &Accumulator{}

	c.Machine.Register(">"+name, func(s *stack.Stack) {
		a.Total += s.Pop()
	})

	c.Machine.Register(name+">", func(s *stack.Stack) {
		s.Push(a.Total)
		a.Total = 0
	})
	return a
}
