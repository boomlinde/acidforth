package midi

import (
	"github.com/boomlinde/acidforth/collection"
	"github.com/boomlinde/acidforth/machine"
	"github.com/rakyll/portmidi"
)

type Midi struct {
	KeyHooks       [128]uint8
	ControlHooks   [128]uint8
	MomentaryHooks [128]uint8
	VelocityHooks  [128]uint8
	ch             <-chan portmidi.Event
	patch          uint8
}

func NewMidi(ch <-chan portmidi.Event) *Midi {
	return &Midi{ch: ch}
}

func (m *Midi) Listen() {
	for event := range m.ch {
		msg := event.Status >> 4
		switch {
		case msg == 9:
			if m.KeyHooks[event.Data1&0x7f] == 0 {
				m.KeyHooks[event.Data1&0x7f] = 1
			} else {
				m.KeyHooks[event.Data1&0x7f] = 0
			}
			m.MomentaryHooks[event.Data1&0x7f] = 1
			m.VelocityHooks[event.Data1&0x7f] = uint8(event.Data2 & 0x7f)
		case msg == 11:
			m.ControlHooks[event.Data1] = uint8(event.Data2 & 0x7f)
		case msg == 8:
			m.MomentaryHooks[event.Data1&0x7f] = 0
		case msg == 12:
			m.patch = uint8(event.Data2 & 0x7f)
		}
	}
}

func (m *Midi) Register(c *collection.Collection) {
	c.Machine.Register("patch", func(s *machine.Stack) {
		s.Push(float64(m.patch))
	})
	c.Machine.Register("cc", func(s *machine.Stack) {
		h := m.ControlHooks[uint8(s.Pop())&0x7f]
		s.Push(float64(h) / 127)
	})
	c.Machine.Register("key", func(s *machine.Stack) {
		h := m.KeyHooks[uint8(s.Pop())&0x7f]
		s.Push(float64(h))
	})
	c.Machine.Register("mom", func(s *machine.Stack) {
		h := m.MomentaryHooks[uint8(s.Pop())&0x7f]
		s.Push(float64(h))
	})
	c.Machine.Register("vel", func(s *machine.Stack) {
		h := m.VelocityHooks[uint8(s.Pop())&0x7f]
		s.Push(float64(h) / 127)
	})
}
