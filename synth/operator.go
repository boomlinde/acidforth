package synth

import (
	"github.com/boomlinde/gobassline/collection"
	"github.com/boomlinde/gobassline/machine/stack"
	"math"
)

type Operator struct {
	Phase    float64
	PhaseInc float64
	Looped   float64
}

func (o *Operator) Tick() {
	o.Phase = o.Phase + o.PhaseInc
	if o.Phase > 1 {
		_, o.Phase = math.Modf(o.Phase)
		o.Phase = math.Abs(o.Phase)
		o.Looped = 1
	}
}

func NewOperator(name string, c *collection.Collection, srate float64) *Operator {
	o := &Operator{}
	c.Register(o.Tick)

	c.Machine.Register(name, func(s *stack.Stack) {
		o.PhaseInc = s.Pop() / srate
		s.Push(o.Phase)
	})

	c.Machine.Register(name+".rst", func(s *stack.Stack) {
		if s.Pop() != 0 {
			o.Phase = 0
		}
	})
	c.Machine.Register(name+".cycle?", func(s *stack.Stack) {
		s.Push(o.Looped)
		o.Looped = 0
	})
	return o
}
