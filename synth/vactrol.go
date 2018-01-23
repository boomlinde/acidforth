package synth

import (
	"github.com/boomlinde/acidforth/collection"
	"github.com/boomlinde/acidforth/machine"
)

func NewVactrol(name string, c *collection.Collection, srate float64) {
	var v float64
	var attack float64 = 1
	var decay float64 = 1

	c.Machine.Register(name, func(s *machine.Stack) {
		in := s.Pop()
		ret := v + in*in*0.02*attack*(srate*in-v)
		v = ret - ret*0.0015*decay
		s.Push(ret / srate)
	})
	c.Machine.Register(name+".decay", func(s *machine.Stack) {
		decay = s.Pop()
	})
	c.Machine.Register(name+".attack", func(s *machine.Stack) {
		attack = s.Pop()
	})
}
