package synth

import (
	"github.com/boomlinde/acidforth/collection"
	"github.com/boomlinde/acidforth/machine"
	"math"
)

type Operator struct {
	phase    float64
	phaseInc float64
	looped   float64
}

func (o *Operator) Tick() {
	o.phase = o.phase + o.phaseInc
	if o.phase > 1 {
		_, o.phase = math.Modf(o.phase)
		o.phase = math.Abs(o.phase)
		o.looped = 1
	}
}

func NewOperator(name string, c *collection.Collection, srate float64) {
	o := &Operator{}
	c.Register(o.Tick)

	c.Machine.Register(name, func(s *machine.Stack) {
		o.phaseInc = s.Pop() / srate
		s.Push(o.phase)
	})

	c.Machine.Register(name+".rst", func(s *machine.Stack) {
		if s.Pop() != 0 {
			o.phase = 0
		}
	})
	c.Machine.Register(name+".cycle?", func(s *machine.Stack) {
		s.Push(o.looped)
		o.looped = 0
	})
}
