package machine

import (
	"strconv"
	"strings"
)

type Instruction func(*Stack)
type Program []Instruction
type Machine struct {
	program        Program
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

func (m *Machine) Compile(source []string) error {
	m.program = make(Program, 0)
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
			m.program = append(m.program, genFloatFunc(val))
		} else {
			m.program = append(m.program, ins)
		}
	}
	return nil
}

func (m *Machine) Run() {
	for _, v := range m.program {
		v(m.stack)
	}
}

func NewMachine() *Machine {
	m := &Machine{
		program:        make(Program, 0),
		stack:          NewStack(0xff),
		secondaryStack: NewStack(0xff),
		words:          make(map[string]Instruction),
	}
	basicInstructions(m)
	return m
}
