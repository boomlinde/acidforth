package collection

import (
	"github.com/boomlinde/acidforth/machine"
	"sync"
)

type Collection struct {
	Mutex   sync.Mutex
	tickers []func()
	Machine *machine.Machine
	out1    float32
	out2    float32
	Playing bool
}

func (c *Collection) Register(ticker func()) {
	c.tickers = append(c.tickers, ticker)
}

func (c *Collection) Callback(buf [][]float32) {
	c.Machine.UpdateSafep()
	c.Mutex.Lock()
	for i := range buf[0] {
		for _, t := range c.tickers {
			t()
		}
		c.out1 = 0
		c.out2 = 0
		c.Machine.Run()
		if c.out1 > 1 {
			c.out1 = 1
		} else if c.out1 < -1 {
			c.out1 = -1
		}
		if c.out2 > 1 {
			c.out2 = 1
		} else if c.out2 < -1 {
			c.out2 = -1
		}
		buf[0][i] = c.out1
		buf[1][i] = c.out2
	}
	c.Mutex.Unlock()
}

func NewCollection() *Collection {
	col := &Collection{
		tickers: make([]func(), 0),
		Machine: machine.NewMachine(),
	}
	col.Machine.Register(">out1", func(s *machine.Stack) { col.out1 = float32(s.Pop()) })
	col.Machine.Register(">out2", func(s *machine.Stack) { col.out2 = float32(s.Pop()) })
	col.Machine.Register(">out", func(s *machine.Stack) {
		o := float32(s.Pop())
		col.out1 = o
		col.out2 = o
	})
	col.Machine.Compile(machine.TokenizeString(""))
	return col
}
