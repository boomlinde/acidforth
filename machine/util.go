package machine

func genFloatFunc(val float64) Instruction {
	return func(s *Stack) { s.Push(val) }
}
