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
	Attack   float64
	Decay    float64
	Release  float64
	Current  float64
	State    int
	Gate     bool
	LastGate bool
}

func (e *Envelope) Tick() {
	switch {
	case e.State == ENV_A:
		e.Current += e.Attack
		if e.Current > 1 {
			e.Current = 1
			e.State = ENV_D
		}
	case e.State == ENV_D:
		e.Current -= e.Decay
		if e.Current < 0 {
			e.Current = 0
		}
	case e.State == ENV_R:
		e.Current -= e.Release
		if e.Current < 0 {
			e.Current = 0
		}
	}
}

func NewEnvelope(name string, c *collection.Collection, srate float64) *Envelope {
	e := &Envelope{}
	c.Register(e.Tick)

	c.Machine.Register(name, func(s *stack.Stack) {
		e.Gate = s.Pop() != 0
		if e.LastGate != e.Gate {
			e.LastGate = e.Gate
			if e.Gate {
				e.State = ENV_A
			} else {
				e.State = ENV_R
			}
		}
		s.Push(e.Current)
	})

	c.Machine.Register(name+".a", func(s *stack.Stack) {
		e.Attack = 1 / (s.Pop()*srate + 1)
	})
	c.Machine.Register(name+".d", func(s *stack.Stack) {
		e.Decay = 1 / (s.Pop()*srate + 1)
	})
	c.Machine.Register(name+".r", func(s *stack.Stack) {
		e.Release = 1 / (s.Pop()*srate + 1)
	})
	return e
}
