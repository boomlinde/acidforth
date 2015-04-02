package main

import (
	"code.google.com/p/portaudio-go/portaudio"
	"flag"
	"fmt"
	"github.com/boomlinde/acidforth/collection"
	"github.com/boomlinde/acidforth/machine"
	"github.com/boomlinde/acidforth/midi"
	"github.com/boomlinde/acidforth/synth"
	"github.com/rakyll/portmidi"
	"io/ioutil"
	"log"
	"math"
	"os"
	"time"
)

const sfreq = 44100

func main() {
	var listMidi bool
	var midiInterface int
	var m *midi.Midi

	flag.BoolVar(&listMidi, "l", false, "List MIDI interfaces")
	flag.IntVar(&midiInterface, "m", -1, "Connect to MIDI interface ID")
	flag.Parse()

	portmidi.Initialize()
	defer portmidi.Terminate()

	if listMidi {
		deviceCount := portmidi.CountDevices()
		for i := 0; i < deviceCount; i++ {
			fmt.Println(i, portmidi.GetDeviceInfo(portmidi.DeviceId(i)))
		}
		os.Exit(0)
	}

	if midiInterface != -1 {
		in, err := portmidi.NewInputStream(portmidi.DeviceId(midiInterface), 1024)
		if err != nil {
			log.Fatal(err)
		}
		defer in.Close()
		m = midi.NewMidi(in.Listen())
	}

	col := collection.NewCollection()
	addComponents(sfreq, col)

	log.Println("Waiting for source on stdin")
	data, err := ioutil.ReadAll(os.Stdin)
	chk(err)

	tokens := machine.TokenizeBytes(data)

	if m != nil {
		tokens = m.GetHooks(col, tokens)
		go m.Listen()
	}

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
	for i := 1; i < 9; i++ {
		_ = synth.NewDSeq(fmt.Sprintf("dseq%d", i), c)
	}
	for _, v := range flag.Args() {
		_ = synth.NewSampler(v, c, srate)
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
