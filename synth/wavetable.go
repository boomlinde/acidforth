package synth

import (
	"github.com/boomlinde/acidforth/collection"
	"github.com/boomlinde/acidforth/machine"
	"math"
)

func waveTable(table []float64, phase float64) float64 {
	return table[int(phase*0x10000)&0xffff]
}

func NewWaveTables(c *collection.Collection) {
	sintab := make([]float64, 0x10000)
	tritab := make([]float64, 0x10000)

	for i := range sintab {
		phase := float64(i) * math.Pi / 0x8000
		sintab[i] = math.Sin(phase)
		tritab[i] = 2 * math.Asin(math.Sin(phase)) / math.Pi
	}
	c.Machine.Register("sintab", func(s *machine.Stack) {
		phase := s.Pop()
		s.Push(waveTable(sintab, phase))
	})
	c.Machine.Register("tritab", func(s *machine.Stack) {
		phase := s.Pop()
		s.Push(waveTable(tritab, phase))
	})
}
