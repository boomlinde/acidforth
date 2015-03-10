package collection

import (
	"github.com/boomlinde/gobassline/machine"
	"sync"
)

type Collection struct {
	Mutex   sync.Mutex
	Tickers []func()
	Machine *machine.Machine
}

func (c *Collection) Register(ticker func()) {
	c.Tickers = append(c.Tickers, ticker)
}

func (c *Collection) Callback(buf [][]float32) {
	c.Mutex.Lock()
	for i := range buf[0] {
		for _, t := range c.Tickers {
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
		Tickers: make([]func(), 0),
		Machine: machine.NewMachine(),
	}
	col.Machine.Compile(machine.TokenizeString("0"))
	return col
}
