package synth

import (
	"github.com/boomlinde/acidforth/collection"
	"github.com/boomlinde/acidforth/machine"
	"math"
)

type Delay struct {
	buffer []float64
	index  int
	sample float64
	length int
}

func (d *Delay) Tick() {
	d.buffer[d.index] = d.sample
	d.index = int(math.Mod(float64(d.index+1), float64(d.length)))
}

func NewDelay(name string, c *collection.Collection, srate float64) *Delay {
	d := &Delay{length: 44100, buffer: make([]float64, int(srate))}
	c.Register(d.Tick)

	c.Machine.Register(name, func(s *machine.Stack) {
		d.length = int(math.Abs(s.Pop()) * srate)
		if d.length > len(d.buffer) {
			d.length = len(d.buffer)
		}
	})

	c.Machine.Register(">"+name, func(s *machine.Stack) {
		d.sample = s.Pop()
	})

	c.Machine.Register(name+">", func(s *machine.Stack) {
		s.Push(d.buffer[d.index])
	})

	return d
}
