package synth

import (
	"github.com/boomlinde/acidforth/collection"
	"github.com/boomlinde/acidforth/machine"
)

type DSeq struct {
	index  uint32
	length uint32
	last   bool
}

func NewDSeq(name string, c *collection.Collection) *DSeq {
	d := &DSeq{length: 8}
	c.Machine.Register(name, func(s *machine.Stack) {
		if !c.Playing {
			d.index = 0
		}
		pattern := uint32(s.Pop())
		now := s.Pop() != 0
		if now && !d.last {
			if (pattern>>(d.length-1-d.index))&1 == 1 {
				s.Push(1)
			} else {
				s.Push(0)
			}
			d.index++
			if d.index >= d.length {
				d.index = 0
			}
		} else {
			s.Push(0)
		}
		d.last = now
	})

	c.Machine.Register(name+".len", func(s *machine.Stack) {
		d.length = uint32(s.Pop())
	})

	return d
}
