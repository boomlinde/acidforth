package synth

import (
	"github.com/boomlinde/acidforth/collection"
	"github.com/boomlinde/acidforth/machine"
)

type Delay struct {
	buffer []float64
	index  int
	sample float64
	length int
}

func (d *Delay) Tick() {
	d.buffer[d.index] = d.sample
	d.index += 1
	if d.index == len(d.buffer) {
		d.index -= len(d.buffer)
	}
}

func NewDelay(name string, c *collection.Collection, srate float64) {
	d := &Delay{length: int(srate), buffer: make([]float64, 5*int(srate))}

	c.Machine.Register(name, func(s *machine.Stack) {
		d.length = int(s.Pop() * srate)
		if d.length < 0 {
			d.length = 0
		} else if d.length > len(d.buffer) {
			d.length = len(d.buffer)
		}
	})

	c.Machine.Register(">"+name, func(s *machine.Stack) {
		d.sample = s.Pop()
	})

	c.Machine.Register(name+">", func(s *machine.Stack) {
		d.Tick()
		index := d.index - d.length
		if index < 0 {
			index += len(d.buffer)
		}

		s.Push(d.buffer[index])
	})
}
