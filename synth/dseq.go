package synth

import (
	"github.com/boomlinde/acidforth/collection"
	"github.com/boomlinde/acidforth/machine"
)

type DSeq struct {
	index   uint32
	length  uint32
	Trigged bool
}

func (d *DSeq) Trig() {
	d.index++
	if d.index > d.length {
		d.index = 1
	}
	d.Trigged = true
}

func (d *DSeq) Rel() {
	d.Trigged = false
}

func NewDSeq(name string, c *collection.Collection) *DSeq {
	d := &DSeq{length: 16}
	c.Machine.Register(name, func(s *machine.Stack) {
		if !c.Playing {
			d.index = 0
		}
		pattern := uint32(s.Pop())
		if (pattern>>(d.length-d.index))&1 == 1 && d.Trigged && c.Playing {
			s.Push(1)
		} else {
			s.Push(0)
		}
	})

	c.Machine.Register(name+".len", func(s *machine.Stack) {
		d.length = uint32(s.Pop())
	})

	return d
}
