package collection

import (
	"github.com/boomlinde/acidforth/machine"
	"sync"
)

type Collection struct {
	Mutex   sync.Mutex
	tickers []func()
	Machine *machine.Machine
}

func (c *Collection) Register(ticker func()) {
	c.tickers = append(c.tickers, ticker)
}

func (c *Collection) Callback(buf [][]float32) {
	c.Mutex.Lock()
	for i := range buf[0] {
		for _, t := range c.tickers {
			t()
		}
		c.Machine.Run()
		v := float32(c.Machine.Last())
		buf[0][i] = v
	}
	c.Mutex.Unlock()
}

func NewCollection() *Collection {
	col := &Collection{
		tickers: make([]func(), 0),
		Machine: machine.NewMachine(),
	}
	col.Machine.Compile(machine.TokenizeString("0"))
	return col
}
