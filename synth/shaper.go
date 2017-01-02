package synth

import (
	"github.com/boomlinde/acidforth/collection"
	"github.com/boomlinde/acidforth/machine"
	"math"
)

func shaper(input float64, shape float64) float64 {
	log := false
	mul := 1.0
	out := 0.0

	input = math.Max(math.Min(input, 1), -1)
	shape = math.Max(math.Min(shape, 1), -1) * 10

	if input < 0 {
		mul = -1
	}

	input *= mul
	log = shape < 0
	if log {
		input = 1 - input
		shape *= -1
	}

	out = math.Pow(input, shape+1)
	if log {
		out = 1 - out
	}
	return out * mul
}

func NewShaper(c *collection.Collection) {
	c.Machine.Register("shaper", func(s *machine.Stack) {
		phase := s.Pop()
		input := s.Pop()
		s.Push(shaper(input, phase))
	})
}
