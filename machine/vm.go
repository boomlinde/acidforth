package machine

import (
	"github.com/boomlinde/gobassline/machine/stack"
	"strconv"
)

type Instruction func(*stack.Stack)
type Program []Instruction
type Machine struct {
	program        Program
	words          map[string]Instruction
	stack          *stack.Stack
	secondaryStack *stack.Stack
}

func (m *Machine) Register(name string, f Instruction) {
	m.words[name] = f
}

func (m *Machine) Compile(source []string) error {
	m.program = make(Program, 0)
	for _, word := range source {
		ins := m.words[word]
		if ins == nil {
			val, err := strconv.ParseFloat(word, 64)
			if err != nil {
				return err
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

func (m *Machine) Last() float64 { return m.stack.Pop() }

func NewMachine() *Machine {
	m := &Machine{
		program:        make(Program, 0),
		stack:          stack.NewStack(0xff),
		secondaryStack: stack.NewStack(0xff),
		words:          make(map[string]Instruction),
	}
	basicInstructions(m)
	return m
}
