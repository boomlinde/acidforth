package synth

import (
	"github.com/boomlinde/acidforth/collection"
	"github.com/boomlinde/acidforth/machine"
	"math"
	"math/rand"
)

type Note struct {
	Tone   float64
	Octave float64
	Gate   bool
	Accent bool
	Slide  bool
}

type Seq struct {
	tempo       float64
	phase       float64
	phaseInc    float64
	pattern     []Note
	nextPattern []Note
	lastSeed    float64
	step        int
	baseNote    float64
	lastSlide   bool
	lastTone    float64
	slideFactor float64
	trigState   bool
	srate       float64
	length      int

	slideRate     float64
	currentTone   float64
	currentGate   float64
	currentAccent float64
}

func genNote(gen *rand.Rand) Note {
	return Note{
		Tone:   float64(gen.Intn(12)),
		Octave: float64(gen.Intn(3) - 1),
		Gate:   gen.Intn(9) > 1,
		Accent: gen.Intn(7) > 2,
		Slide:  gen.Intn(7) > 3,
	}
}

func genPattern(seed float64, queue chan []Note) {
	gen := rand.New(rand.NewSource(int64(seed)))

	p := make([]Note, 16)
	for i := range p {
		p[i] = genNote(gen)
	}

	queue <- p
}

func (s *Seq) Tick() {
	s.phase = s.phase + s.phaseInc
	s.currentTone += s.slideRate
	if s.phase > 1 {
		_, s.phase = math.Modf(s.phase)
		s.phase = math.Abs(s.phase)
		s.Trig()
	}
}

func (s *Seq) Trig() {
	if s.trigState {
		if s.step >= len(s.pattern) || s.step >= s.length {
			s.step = 0
			s.pattern = s.nextPattern
		}
		step := s.pattern[s.step]
		s.step += 1
		if step.Gate {
			s.currentGate = 1
		} else {
			s.currentGate = 0
		}
		if s.lastSlide {
			t := step.Tone + step.Octave*12 + s.baseNote
			s.slideRate = (t - s.lastTone) * s.slideFactor
		} else {
			s.currentTone = step.Tone + step.Octave*12 + s.baseNote
			s.slideRate = 0
			if step.Accent {
				s.currentAccent = 1
			} else {
				s.currentAccent = 0
			}
		}
		s.lastSlide = step.Slide
		s.lastTone = s.currentTone
	} else {
		step := s.pattern[s.step-1]
		if !step.Slide {
			s.currentGate = 0
		}
	}

	s.trigState = !s.trigState
}

func (s *Seq) SetTempo(tempo float64) {
	s.tempo = tempo
	s.slideFactor = 1 / (s.srate * 15.0 / tempo)
	s.phaseInc = 8 * tempo / 60 / s.srate
}

func (s *Seq) SetPattern(p []Note) {
	s.nextPattern = p
}

func NewSeq(name string, c *collection.Collection, srate float64) *Seq {
	se := &Seq{trigState: true, srate: srate, baseNote: 60, length: 16}
	se.SetTempo(140)
	c.Register(se.Tick)

	c.Machine.Register(name+".pitch", func(s *machine.Stack) {
		s.Push(se.currentTone)
	})
	c.Machine.Register(name+".gate", func(s *machine.Stack) {
		s.Push(se.currentGate)
	})
	c.Machine.Register(name+".accent", func(s *machine.Stack) {
		s.Push(se.currentAccent)
	})
	c.Machine.Register(name+".tune", func(s *machine.Stack) {
		se.baseNote = 60 + s.Pop()
	})
	c.Machine.Register(name+".tempo", func(s *machine.Stack) {
		tempo := s.Pop()
		if se.tempo != tempo {
			se.SetTempo(tempo)
		}
	})
	c.Machine.Register(name+".pattern", func(s *machine.Stack) {
		seed := s.Pop()
		if seed != se.lastSeed {
			se.lastSeed = seed
			queue := make(chan []Note, 1)
			go genPattern(seed, queue)
			se.nextPattern = <-queue
		}
	})
	c.Machine.Register(name+".len", func(s *machine.Stack) {
		se.length = int(s.Pop())
	})
	c.Machine.Register(name+".trig", func(s *machine.Stack) {
		if !se.trigState {
			s.Push(1)
		} else {
			s.Push(0)
		}
	})

	se.SetPattern([]Note{Note{0, 1, true, false, false}})

	return se
}
