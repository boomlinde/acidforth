package seq

import (
	"github.com/boomlinde/gobassline/collection"
	"github.com/boomlinde/gobassline/machine/stack"
	"math"
)

type Note struct {
	Tone   float64
	Octave float64
	Gate   bool
	Accent bool
	Slide  bool
}

type Seq struct {
	Tempo       float64
	Phase       float64
	PhaseInc    float64
	Pattern     []Note
	Step        int
	BaseNote    float64
	LastSlide   bool
	LastTone    float64
	SlideFactor float64
	TrigState   bool
	srate       float64

	SlideRate     float64
	CurrentTone   float64
	CurrentGate   float64
	CurrentAccent float64
}

func (s *Seq) Tick() {
	s.Phase = s.Phase + s.PhaseInc
	s.CurrentTone += s.SlideRate
	if s.Phase > 1 {
		_, s.Phase = math.Modf(s.Phase)
		s.Phase = math.Abs(s.Phase)
		s.Trig()
	}
}

func (s *Seq) Trig() {
	if s.TrigState {
		if s.Step >= len(s.Pattern) {
			s.Step = 0
		}
		step := s.Pattern[s.Step]
		s.Step += 1
		if step.Gate {
			s.CurrentGate = 1
		} else {
			s.CurrentGate = 0
		}
		if s.LastSlide {
			t := step.Tone + step.Octave*12 + s.BaseNote
			s.SlideRate = (t - s.LastTone) * s.SlideFactor
		} else {
			s.CurrentTone = step.Tone + step.Octave*12 + s.BaseNote
			s.SlideRate = 0
			if step.Accent {
				s.CurrentAccent = 1
			} else {
				s.CurrentAccent = 0
			}
		}
		s.LastSlide = step.Slide
		s.LastTone = s.CurrentTone
	} else {
		step := s.Pattern[s.Step-1]
		if !step.Slide {
			s.CurrentGate = 0
		}
	}
	s.TrigState = !s.TrigState
}

func (s *Seq) SetTempo(tempo float64) {
	s.Tempo = tempo
	s.SlideFactor = 1 / (s.srate * 15.0 / tempo)
	s.PhaseInc = 8 * tempo / 60 / s.srate
}

func (s *Seq) SetPattern(p []Note) {
	s.Pattern = p
}

func NewSeq(name string, c *collection.Collection, srate float64) *Seq {
	se := &Seq{TrigState: true, srate: srate, BaseNote: 40}
	se.SetTempo(140)
	c.Register(se.Tick)

	c.Machine.Register(name+".pitch", func(s *stack.Stack) {
		s.Push(se.CurrentTone)
	})
	c.Machine.Register(name+".gate", func(s *stack.Stack) {
		s.Push(se.CurrentGate)
	})
	c.Machine.Register(name+".accent", func(s *stack.Stack) {
		s.Push(se.CurrentAccent)
	})

	se.SetPattern([]Note{
		Note{7, 1, true, false, true},
		Note{0, 1, true, false, false},
		Note{0, 1, true, true, false},
		Note{10, 0, true, true, false},
		Note{0, 1, true, false, false},
		Note{0, 0, true, false, false},
		Note{1, 2, true, true, false},
		Note{4, 1, true, false, true},
		Note{7, 1, true, false, true},
		Note{7, 1, true, false, false},
		Note{0, 1, true, false, false},
		Note{10, 0, true, true, false},
		Note{0, 1, true, true, false},
		Note{0, 0, true, false, false},
		Note{1, 2, true, true, false},
		Note{0, 0, false, false, false},
	})

	return se
}
