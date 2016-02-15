package machine

import (
	"strconv"
	"strings"
	"sync"
)

type Instruction func(*Stack)
type Program []Instruction
type Machine struct {
	program        *Program
	safep          *Program
	plock          *sync.Mutex
	words          map[string]Instruction
	stack          *Stack
	secondaryStack *Stack
}

func (m *Machine) Register(name string, f Instruction) {
	m.words[name] = f
}

func StripComments(source []string) []string {
	var comment bool
	out := make([]string, 0)
	for _, word := range source {
		if word == "(" {
			comment = true
			continue
		} else if comment {
			if strings.HasSuffix(word, ")") {
				comment = false
			}
			continue
		}
		out = append(out, word)
	}
	return out
}

func (m *Machine) Build(source []byte) error {
	tokens := TokenizeBytes(source)
	tokens = StripComments(tokens)
	tokens, err := ExpandMacros(tokens)
	if err != nil {
		return err
	}
	return m.Compile(tokens)
}

func (m *Machine) Compile(source []string) error {
	program := make(Program, 0)
	for _, word := range source {
		ins := m.words[word]
		if ins == nil {
			var val float64
			val, err := strconv.ParseFloat(word, 64)
			if err != nil {
				vi, err := strconv.ParseUint(word, 0, 32)
				if err != nil {
					if word[len(word)-1] == 'b' {
						vi, err = strconv.ParseUint(word[:len(word)-1], 2, 32)
						if err != nil {
							return err
						}
					} else {
						return err
					}
				}
				val = float64(vi)
			}
			program = append(program, genFloatFunc(val))
		} else {
			program = append(program, ins)
		}
	}

	m.plock.Lock()
	m.program = &program
	m.plock.Unlock()
	return nil
}

func (m *Machine) Run() {
	for _, v := range *m.safep {
		v(m.stack)
	}
}

func (m *Machine) UpdateSafep() {
	m.plock.Lock()
	m.safep = m.program
	m.plock.Unlock()
}

func NewMachine() *Machine {
	program := make(Program, 0)
	m := &Machine{
		program:        &program,
		plock:          &sync.Mutex{},
		stack:          NewStack(0xff),
		secondaryStack: NewStack(0xff),
		words:          make(map[string]Instruction),
	}
	basicInstructions(m)
	return m
}
