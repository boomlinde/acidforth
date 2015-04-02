package midi

import (
	"github.com/boomlinde/acidforth/collection"
	"github.com/boomlinde/acidforth/machine"
	"github.com/rakyll/portmidi"
	"log"
	"strconv"
	"strings"
	"sync"
)

type Hook struct {
	Value float64
	Lock  *sync.Mutex
}

func (h *Hook) Bind(s *machine.Stack) {
	h.Lock.Lock()
	s.Push(h.Value)
	h.Lock.Unlock()
}

type Midi struct {
	KeyHooks       map[int64]*Hook
	ControlHooks   map[int64]*Hook
	MomentaryHooks map[int64]*Hook
	ch             <-chan portmidi.Event
	patch          Hook
}

func NewMidi(ch <-chan portmidi.Event) *Midi {
	return &Midi{
		KeyHooks:       make(map[int64]*Hook),
		ControlHooks:   make(map[int64]*Hook),
		MomentaryHooks: make(map[int64]*Hook),
		ch:             ch,
		patch:          Hook{Lock: &sync.Mutex{}},
	}
}

func (m *Midi) Listen() {
	for event := range m.ch {
		msg := event.Status >> 4
		switch {
		case msg == 9:
			h, ok := m.KeyHooks[event.Data1]
			if ok {
				h.Lock.Lock()
				if h.Value == 0 {
					h.Value = 1
				} else {
					h.Value = 0
				}
				h.Lock.Unlock()
			}
			h, ok = m.MomentaryHooks[event.Data1]
			if ok {
				h.Lock.Lock()
				h.Value = 1
				h.Lock.Unlock()
			}
		case msg == 11:
			h, ok := m.ControlHooks[event.Data1]
			if ok {
				h.Lock.Lock()
				h.Value = float64(event.Data2) / 127
				h.Lock.Unlock()
			}
		case msg == 8:
			h, ok := m.MomentaryHooks[event.Data1]
			if ok {
				h.Lock.Lock()
				h.Value = 0
				h.Lock.Unlock()
			}
		case msg == 12:
			m.patch.Lock.Lock()
			m.patch.Value = float64(event.Data1)
			m.patch.Lock.Unlock()
		}
	}
}

func (m *Midi) GetHooks(c *collection.Collection, tokens []string) []string {
	out := make([]string, 0)
	for _, token := range tokens {
		if token[0] == '#' {
			d := strings.Split(token[1:], ":")
			name := d[0]
			log.Println("Registering:", name)
			if len(d) != 3 {
				log.Fatal("Badly formatted MIDI hook")
			}
			ctype := d[1]
			val, err := strconv.ParseInt(d[2], 0, 64)
			if err != nil {
				log.Fatal(err)
			}
			hook := &Hook{Lock: &sync.Mutex{}}
			switch {
			case ctype == "cc":
				m.ControlHooks[val] = hook
			case ctype == "key":
				m.KeyHooks[val] = hook
			case ctype == "mom":
				m.MomentaryHooks[val] = hook
			default:
				log.Fatal("Unsupported hook type")
			}
			c.Machine.Register(name, hook.Bind)
			c.Machine.Register("patch", func(s *machine.Stack) {
				m.patch.Lock.Lock()
				s.Push(m.patch.Value)
				m.patch.Lock.Unlock()
			})
		} else {
			out = append(out, token)
		}
	}
	return out
}
