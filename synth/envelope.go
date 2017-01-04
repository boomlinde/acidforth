package synth

import (
	"github.com/boomlinde/acidforth/collection"
	"github.com/boomlinde/acidforth/machine"
)

const (
	ENV_R = iota
	ENV_A
	ENV_D
)

type Envelope struct {
	attack   float64
	decay    float64
	sustain  float64
	release  float64
	current  float64
	state    int
	gate     bool
	lastGate bool
}

func (e *Envelope) Tick() {
	switch e.state {
	case ENV_A:
		e.current += e.attack
		if e.current > 1 {
			e.current = 1
			e.state = ENV_D
		}
	case ENV_D:
		e.current -= e.decay * (1 - e.sustain)
		if e.current < e.sustain {
			e.current = e.sustain
		}
	case ENV_R:
		e.current -= e.release
		if e.current < 0 {
			e.current = 0
		}
	}
}

func getrate(s *machine.Stack, srate float64) float64 {
	return 1 / (s.Pop()*srate + 1)
}

func NewEnvelope(name string, c *collection.Collection, srate float64) *Envelope {
	e := &Envelope{}
	c.Register(e.Tick)

	c.Machine.Register(name, func(s *machine.Stack) {
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

	c.Machine.Register(name+".a", func(s *machine.Stack) {
		e.attack = getrate(s, srate)
	})
	c.Machine.Register(name+".d", func(s *machine.Stack) {
		e.decay = getrate(s, srate)
	})
	c.Machine.Register(name+".s", func(s *machine.Stack) {
		e.sustain = s.Pop()
	})
	c.Machine.Register(name+".r", func(s *machine.Stack) {
		e.release = getrate(s, srate)
	})
	c.Machine.Register(name+".adsr", func(s *machine.Stack) {
		e.release = getrate(s, srate)
		e.sustain = s.Pop()
		e.decay = getrate(s, srate)
		e.attack = getrate(s, srate)
	})
	return e
}
