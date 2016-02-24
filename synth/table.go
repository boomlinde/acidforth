package synth

import (
	"github.com/boomlinde/acidforth/collection"
	"github.com/boomlinde/acidforth/machine"
	"math"
)

type ITable struct {
	length int
	values [32]float64
}

func NewITable(name string, c *collection.Collection) {
	it := &ITable{length: 8}

	var bottom float64
	var top float64
	var mix float64

	c.Machine.Register(name, func(s *machine.Stack) {
		index := math.Abs(s.Pop())
		bottom = math.Floor(index)
		top = math.Ceil(index)
		mix = index - bottom
		s.Push(it.values[int(bottom)%it.length]*(1-mix) + it.values[int(top)%it.length]*mix)
	})

	c.Machine.Register(name+".set", func(s *machine.Stack) {
		length := s.Pop()
		if length > 32 {
			length = 32
		} else if length < 1 {
			length = 1
		}
		it.length = int(length)
		for i := it.length - 1; i >= 0; i-- {
			it.values[i] = s.Pop()
		}
	})

	c.Machine.Register(name+".len", func(s *machine.Stack) {
		s.Push(float64(it.length))
	})
}
