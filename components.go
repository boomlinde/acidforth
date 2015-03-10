package main

import (
	"github.com/boomlinde/gobassline/collection"
	"github.com/boomlinde/gobassline/synth"
)

func addComponents(c *collection.Collection) {
	_ = synth.NewOperator("op1", c)
	_ = synth.NewOperator("op2", c)
	_ = synth.NewOperator("op3", c)
	_ = synth.NewOperator("op4", c)
	_ = synth.NewAccumulator("mix1", c)
	_ = synth.NewAccumulator("mix2", c)
	_ = synth.NewAccumulator("mix3", c)
	_ = synth.NewAccumulator("mix4", c)
	synth.NewWaveTables(c)
}
