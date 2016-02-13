package midi

import (
	"github.com/boomlinde/acidforth/collection"
	"github.com/boomlinde/acidforth/machine"
	"github.com/rakyll/portmidi"
)

type Hook struct {
	Value uint8
}

type Midi struct {
	KeyHooks       []*Hook
	ControlHooks   []*Hook
	MomentaryHooks []*Hook
	VelocityHooks  []*Hook
	ch             <-chan portmidi.Event
	patch          Hook
}

func NewMidi(ch <-chan portmidi.Event) *Midi {
	keyHooks := make([]*Hook, 128)
	controlHooks := make([]*Hook, 128)
	momentaryHooks := make([]*Hook, 128)
	velocityHooks := make([]*Hook, 128)
	for i := 0; i < 128; i++ {
		keyHooks[i] = &Hook{}
		controlHooks[i] = &Hook{}
		momentaryHooks[i] = &Hook{}
		velocityHooks[i] = &Hook{}
	}

	return &Midi{
		KeyHooks:       keyHooks,
		ControlHooks:   controlHooks,
		MomentaryHooks: momentaryHooks,
		VelocityHooks:  velocityHooks,
		ch:             ch,
		patch:          Hook{},
	}
}

func (m *Midi) Listen() {
	for event := range m.ch {
		msg := event.Status >> 4
		switch {
		case msg == 9:
			h := m.KeyHooks[event.Data1&0x7f]
			if h.Value == 0 {
				h.Value = 1
			} else {
				h.Value = 0
			}
			h = m.MomentaryHooks[event.Data1&0x7f]
			h.Value = 1
			h = m.VelocityHooks[event.Data1&0x7f]
			h.Value = uint8(event.Data2 & 0x7f)
		case msg == 11:
			h := m.ControlHooks[event.Data1]
			h.Value = uint8(event.Data2 & 0x7f)
		case msg == 8:
			h := m.MomentaryHooks[event.Data1&0x7f]
			h.Value = 0
		case msg == 12:
			m.patch.Value = uint8(event.Data2 & 0x7f)
		}
	}
}

func (m *Midi) Register(c *collection.Collection) {
	c.Machine.Register("patch", func(s *machine.Stack) {
		s.Push(float64(m.patch.Value))
	})
	c.Machine.Register("cc", func(s *machine.Stack) {
		h := m.ControlHooks[int(s.Pop())]
		s.Push(float64(h.Value) / 127)
	})
	c.Machine.Register("key", func(s *machine.Stack) {
		h := m.KeyHooks[int(s.Pop())]
		s.Push(float64(h.Value))
	})
	c.Machine.Register("mom", func(s *machine.Stack) {
		h := m.MomentaryHooks[int(s.Pop())]
		s.Push(float64(h.Value))
	})
	c.Machine.Register("vel", func(s *machine.Stack) {
		h := m.VelocityHooks[int(s.Pop())]
		s.Push(float64(h.Value) / 127)
	})
}
