package main

import (
	"github.com/boomlinde/gobassline/collection"
	"github.com/boomlinde/gobassline/machine/stack"
	"github.com/boomlinde/gobassline/synth"
	"math"
)

func addComponents(srate float64, c *collection.Collection) {
	_ = synth.NewOperator("op1", c, srate)
	_ = synth.NewOperator("op2", c, srate)
	_ = synth.NewOperator("op3", c, srate)
	_ = synth.NewOperator("op4", c, srate)
	_ = synth.NewOperator("op5", c, srate)
	_ = synth.NewOperator("op6", c, srate)
	_ = synth.NewOperator("op7", c, srate)
	_ = synth.NewOperator("op8", c, srate)
	_ = synth.NewAccumulator("mix1", c)
	_ = synth.NewAccumulator("mix2", c)
	_ = synth.NewAccumulator("mix3", c)
	_ = synth.NewAccumulator("mix4", c)
	_ = synth.NewRegister("A", c)
	_ = synth.NewRegister("B", c)
	_ = synth.NewRegister("C", c)
	_ = synth.NewRegister("D", c)
	synth.NewWaveTables(c)

	c.Machine.Register("srate", func(s *stack.Stack) { s.Push(srate) })
	c.Machine.Register("m2f", func(s *stack.Stack) {
		s.Push(440 * math.Pow(2, (s.Pop()-69)/12))
	})
}
