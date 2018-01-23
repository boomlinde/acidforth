package synth

import (
	"github.com/boomlinde/acidforth/collection"
	"github.com/boomlinde/acidforth/machine"
)

type Accumulator struct {
	total float64
}

func NewAccumulator(name string, c *collection.Collection) {
	a := &Accumulator{}

	c.Machine.Register(">"+name, func(s *machine.Stack) {
		a.total += s.Pop()
	})

	c.Machine.Register(name+">", func(s *machine.Stack) {
		s.Push(a.total)
		a.total = 0
	})
}
