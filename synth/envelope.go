package synth

import (
	"github.com/boomlinde/acidforth/collection"
	"github.com/boomlinde/acidforth/machine/stack"
)

const (
	ENV_R = iota
	ENV_A
	ENV_D
)

type Envelope struct {
	attack   float64
	decay    float64
	release  float64
	current  float64
	state    int
	gate     bool
	lastGate bool
}

func (e *Envelope) Tick() {
	switch {
	case e.state == ENV_A:
		e.current += e.attack
		if e.current > 1 {
			e.current = 1
			e.state = ENV_D
		}
	case e.state == ENV_D:
		e.current -= e.decay
		if e.current < 0 {
			e.current = 0
		}
	case e.state == ENV_R:
		e.current -= e.release
		if e.current < 0 {
			e.current = 0
		}
	}
}

func NewEnvelope(name string, c *collection.Collection, srate float64) *Envelope {
	e := &Envelope{}
	c.Register(e.Tick)

	c.Machine.Register(name, func(s *stack.Stack) {
		e.gate = s.Pop() != 0
		if e.lastGate != e.gate {
			e.lastGate = e.gate
			if e.gate {
				e.state = ENV_A
			} else {
				e.state = ENV_R
			}
		}
		s.Push(e.current)
	})

	c.Machine.Register(name+".a", func(s *stack.Stack) {
		e.attack = 1 / (s.Pop()*srate + 1)
	})
	c.Machine.Register(name+".d", func(s *stack.Stack) {
		e.decay = 1 / (s.Pop()*srate + 1)
	})
	c.Machine.Register(name+".r", func(s *stack.Stack) {
		e.release = 1 / (s.Pop()*srate + 1)
	})
	return e
}
