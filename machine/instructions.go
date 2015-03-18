package machine

import (
	"fmt"
	"math"
)

func basicInstructions(m *Machine) {
	m.Register("drop", func(s *Stack) {
		_ = s.Pop()
	})
	m.Register("dup", func(s *Stack) {
		v := s.Pop()
		s.Push(v)
		s.Push(v)
	})
	m.Register("swap", func(s *Stack) {
		b := s.Pop()
		a := s.Pop()
		s.Push(b)
		s.Push(a)
	})
	m.Register("rot", func(s *Stack) {
		c := s.Pop()
		b := s.Pop()
		a := s.Pop()
		s.Push(b)
		s.Push(c)
		s.Push(a)
	})
	m.Register("*", func(s *Stack) {
		s.Push(s.Pop() * s.Pop())
	})
	m.Register("+", func(s *Stack) {
		s.Push(s.Pop() + s.Pop())
	})
	m.Register("-", func(s *Stack) {
		b := s.Pop()
		a := s.Pop()
		s.Push(a - b)
	})
	m.Register("/", func(s *Stack) {
		b := s.Pop()
		a := s.Pop()
		s.Push(a / b)
	})
	m.Register("%", func(s *Stack) {
		b := s.Pop()
		a := s.Pop()
		s.Push(math.Mod(a, b))
	})
	m.Register("pi", func(s *Stack) {
		s.Push(math.Pi)
	})
	m.Register("=", func(s *Stack) {
		b := s.Pop()
		a := s.Pop()
		if b == a {
			s.Push(1)
		} else {
			s.Push(0)
		}
	})
	m.Register("<", func(s *Stack) {
		b := s.Pop()
		a := s.Pop()
		if a < b {
			s.Push(1)
		} else {
			s.Push(0)
		}
	})
	m.Register(">=", func(s *Stack) {
		b := s.Pop()
		a := s.Pop()
		if a >= b {
			s.Push(1)
		} else {
			s.Push(0)
		}
	})
	m.Register("<=", func(s *Stack) {
		b := s.Pop()
		a := s.Pop()
		if a <= b {
			s.Push(1)
		} else {
			s.Push(0)
		}
	})
	m.Register(">", func(s *Stack) {
		b := s.Pop()
		a := s.Pop()
		if a > b {
			s.Push(1)
		} else {
			s.Push(0)
		}
	})
	m.Register("~", func(s *Stack) {
		a := s.Pop()
		if a == 0 {
			s.Push(1)
		} else {
			s.Push(0)
		}
	})
	m.Register(".", func(s *Stack) {
		fmt.Println(s.Pop())
	})
	m.Register("push", func(s *Stack) {
		m.secondaryStack.Push(s.Pop())
	})
	m.Register("pop", func(s *Stack) {
		s.Push(m.secondaryStack.Pop())
	})
}
