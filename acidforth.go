package main

import (
	"code.google.com/p/portaudio-go/portaudio"
	"fmt"
	"github.com/boomlinde/acidforth/collection"
	"github.com/boomlinde/acidforth/machine"
	"github.com/boomlinde/acidforth/synth"
	"io/ioutil"
	"log"
	"math"
	"os"
	"time"
)

const sfreq = 44100

func main() {
	log.Println("Booting")
	col := collection.NewCollection()
	addComponents(sfreq, col)

	log.Println("Waiting for source on stdin")
	data, err := ioutil.ReadAll(os.Stdin)
	chk(err)

	tokens := machine.TokenizeBytes(data)

	tokens, err = machine.ExpandMacros(tokens)
	chk(err)

	chk(col.Machine.Compile(tokens))
	log.Println("Running")

	portaudio.Initialize()
	defer portaudio.Terminate()
	stream, err := portaudio.OpenDefaultStream(0, 2, sfreq, 0, col.Callback)
	chk(err)
	defer stream.Close()
	stream.Start()

	for {
		time.Sleep(time.Second)
	}
}

func addComponents(srate float64, c *collection.Collection) {
	for i := 1; i < 9; i++ {
		_ = synth.NewOperator(fmt.Sprintf("op%d", i), c, srate)
		_ = synth.NewEnvelope(fmt.Sprintf("env%d", i), c, srate)
	}
	for _, r := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		_ = synth.NewRegister(string(r), c)
	}
	for i := 1; i < 5; i++ {
		_ = synth.NewAccumulator(fmt.Sprintf("mix%d", i), c)
		_ = synth.NewDelay(fmt.Sprintf("delay%d", i), c, srate)
	}

	_ = synth.NewSeq("seq", c, srate)

	synth.NewWaveTables(c)

	c.Machine.Register("srate", func(s *machine.Stack) { s.Push(srate) })
	c.Machine.Register("m2f", func(s *machine.Stack) {
		s.Push(440 * math.Pow(2, (s.Pop()-69)/12))
	})
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
