package midi

import (
	"sync"
	"time"
)

type note struct {
	t        time.Duration
	velocity uint8
	key      uint8
}

type keyboard struct {
	held  map[uint8]*note
	last  *note
	gate  bool
	mtx   sync.Mutex
	ztime time.Time
}

func (k *keyboard) press(key uint8, velocity uint8) {
	k.held[key] = &note{time.Since(k.ztime), velocity, key}
	k.update()
}

func (k *keyboard) release(key uint8) {
	delete(k.held, key)
	k.update()
}

func (k *keyboard) update() {
	var last *note

	for _, no := range k.held {
		if last == nil || no.t > last.t {
			last = no
		}
	}

	if last != nil {
		(&k.mtx).Lock()
		k.gate = true
		k.last = last
		(&k.mtx).Unlock()
	} else {
		k.gate = false
	}
}

func newKeyboard() *keyboard {
	return &keyboard{held: make(map[uint8]*note), ztime: time.Now()}
}
