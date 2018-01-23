package synth

import (
	"github.com/boomlinde/acidforth/collection"
	"github.com/boomlinde/acidforth/machine"
)

type Pulser struct {
	length   float64
	curwait  float64
	value    float64
	lastsamp float64
}

func NewPulser(name string, c *collection.Collection, srate float64) {
	o := &Pulser{}
	c.Machine.Register(name, func(s *machine.Stack) {
		samp := s.Pop()
		if samp != o.lastsamp && samp > o.lastsamp {
			o.curwait = o.length * srate
			o.value = samp
		}
		o.lastsamp = samp
		if o.curwait > 0 {
			s.Push(o.value)
		} else {
			s.Push(0)
		}
		o.curwait -= 1
	})

	c.Machine.Register(name+".len", func(s *machine.Stack) {
		o.length = s.Pop()
	})
}
