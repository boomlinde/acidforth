package synth

import (
	"github.com/boomlinde/acidforth/collection"
	"github.com/boomlinde/acidforth/machine"
	"github.com/boomlinde/acidforth/synth/audio"
	"log"
	"strings"
)

type Sampler struct {
	index    float64
	rate     float64
	lastGate bool
	out      float64
	data     []float32
}

func (s *Sampler) Tick() {
	i := int(s.index)
	if i < len(s.data)-1 {
		s.index += s.rate
		s.out = float64(s.data[i])
	} else {
		s.out = 0
	}
}

func NewSampler(fname string, c *collection.Collection, srate float64) *Sampler {
	pathSplit := strings.Split(fname, "/")
	name := pathSplit[len(pathSplit)-1]

	log.Printf("Registering sampler: %s", name)

	sound, err := audio.ReadWavFile(fname)
	if err != nil {
		log.Fatal(err)
	}
	sampleRate := float64(sound.Rate)

	s := &Sampler{
		index: float64(len(sound.Data[0]) - 1),
		data:  sound.Data[0],
		rate:  sampleRate / srate,
	}

	c.Register(s.Tick)

	c.Machine.Register(name, func(st *machine.Stack) {
		gate := st.Pop() != 0
		if gate && !s.lastGate {
			s.index = 0
		}
		s.lastGate = gate
		st.Push(s.out)
	})

	c.Machine.Register(name+".rate", func(st *machine.Stack) {
		rate := st.Pop()
		s.rate = sampleRate / srate * rate
	})

	return s
}
