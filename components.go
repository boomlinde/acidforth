package main

import (
	"fmt"
	"github.com/boomlinde/acidforth/collection"
	"github.com/boomlinde/acidforth/machine/stack"
	"github.com/boomlinde/acidforth/seq"
	"github.com/boomlinde/acidforth/synth"
	"math"
)

func addComponents(srate float64, c *collection.Collection) {
	for i := 1; i < 9; i++ {
		_ = synth.NewOperator(fmt.Sprintf("op%d", i), c, srate)
		_ = synth.NewEnvelope(fmt.Sprintf("env%d", i), c, srate)
	}
	for _, r := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		_ = synth.NewRegister(string(r), c)
	}
	for i := 1; i < 5; i++ {
		_ = synth.NewAccumulator(fmt.Sprintf("mix%d", i), c)
	}

	_ = seq.NewSeq("seq", c, srate)

	synth.NewWaveTables(c)

	c.Machine.Register("srate", func(s *stack.Stack) { s.Push(srate) })
	c.Machine.Register("m2f", func(s *stack.Stack) {
		s.Push(440 * math.Pow(2, (s.Pop()-69)/12))
	})
}
