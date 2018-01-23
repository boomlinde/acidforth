package synth

import (
	"github.com/boomlinde/acidforth/collection"
	"github.com/boomlinde/acidforth/machine"
)

type Register struct {
	value float64
}

func NewRegister(name string, c *collection.Collection) {
	a := &Register{}

	c.Machine.Register(">"+name, func(s *machine.Stack) {
		a.value = s.Pop()
	})

	c.Machine.Register(name+">", func(s *machine.Stack) {
		s.Push(a.value)
	})
}
