package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/boomlinde/acidforth/collection"
	"github.com/boomlinde/acidforth/machine"
	"github.com/boomlinde/acidforth/midi"
	"github.com/boomlinde/acidforth/synth"
	"github.com/gordonklaus/portaudio"
	"github.com/rakyll/portmidi"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"
)

var sfreq float64

func main() {
	var listMidi bool
	var midiInterface int
	var address string
	var m *midi.Midi
	var prompt float64
	var samples string
	var sfreqi int

	flag.BoolVar(&listMidi, "l", false, "List MIDI interfaces")
	flag.IntVar(&midiInterface, "m", -1, "Connect to MIDI interface ID")
	flag.StringVar(&address, "s", "", "HTTP server address. Leave unset to disable")
	flag.StringVar(&samples, "w", "", "Sample directory")
	flag.IntVar(&sfreqi, "r", 44100, "Sample rate")
	flag.Parse()
	args := flag.Args()

	sfreq = float64(sfreqi)

	if listMidi {
		portmidi.Initialize()
		defer portmidi.Terminate()
		deviceCount := portmidi.CountDevices()
		for i := 0; i < deviceCount; i++ {
			fmt.Println(i, portmidi.Info(portmidi.DeviceID(i)))
		}
		os.Exit(0)
	} else if midiInterface != -1 {
		portmidi.Initialize()
		defer portmidi.Terminate()
		in, err := portmidi.NewInputStream(portmidi.DeviceID(midiInterface), 1024)
		if err != nil {
			log.Fatal(err)
		}
		defer in.Close()
		m = midi.NewMidi(in.Listen())
	} else {
		m = midi.NewMidi(make(chan portmidi.Event))
	}

	pl := &sync.Mutex{}

	col := collection.NewCollection()

	m.Register(col)
	go m.Listen()

	addComponents(sfreq, col, samples)

	col.Machine.Register("prompt", func(s *machine.Stack) {
		pl.Lock()
		s.Push(prompt)
		pl.Unlock()
	})

	if len(args) > 0 {
		prg, err := ioutil.ReadFile(args[len(args)-1])
		if err != nil {
			log.Println("ERROR:", err)
		}
		log.Println("Compiling program")
		err = col.Machine.Build(prg)
		if err != nil {
			log.Println("ERROR:", err)
		}
	}

	log.Println("Running")

	if address != "" {
		go server(col, address)
	}

	portaudio.Initialize()
	defer portaudio.Terminate()
	stream, err := portaudio.OpenDefaultStream(0, 2, sfreq, 0, col.Callback)
	chk(err)
	defer stream.Close()
	stream.Start()

	reader := bufio.NewReader(os.Stdin)
	numberRe, err := regexp.Compile("[0-9]+\\.?[0-9]*")
	chk(err)

	for {
		text, err := reader.ReadString('\n')
		chk(err)
		found := numberRe.FindString(text)
		if found == "" {
			if col.Playing {
				log.Print("Stopping sequencer")
			} else {
				log.Println("Starting sequencer")
			}
			col.Playing = !col.Playing
		} else {
			tpr, err := strconv.ParseFloat(found, 64)
			if err != nil {
				log.Println(err)
				continue
			}
			pl.Lock()
			prompt = tpr
			pl.Unlock()
		}
	}
}

func addComponents(srate float64, c *collection.Collection, samples string) {
	dseqs := make([]synth.Triggable, 0, 16)
	for i := 1; i <= 8; i++ {
		synth.NewOperator(fmt.Sprintf("op%d", i), c, srate)
		synth.NewEnvelope(fmt.Sprintf("env%d", i), c, srate)
		synth.NewPulser(fmt.Sprintf("pulser%d", i), c, srate)
		synth.NewVactrol(fmt.Sprintf("vactrol%d", i), c, srate)
		dseqs = append(dseqs, synth.NewDSeq(fmt.Sprintf("dseq%d", i), c))
		dseqs = append(dseqs, synth.NewVSeq(fmt.Sprintf("vseq%d", i), c))
	}
	for _, r := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		synth.NewRegister(string(r), c)
	}
	for i := 1; i <= 4; i++ {
		synth.NewAccumulator(fmt.Sprintf("mix%d", i), c)
		synth.NewDelay(fmt.Sprintf("delay%d", i), c, srate)
		synth.NewITable(fmt.Sprintf("itab%d", i), c)
	}

	if samples != "" {
		files, err := filepath.Glob(filepath.Join(samples, "*.wav"))
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range files {
			synth.NewSampler(f, c, srate)
		}
	}

	synth.NewSeq("seq", c, srate, dseqs)
	synth.NewWaveTables(c)
	synth.NewShaper(c)

	c.Machine.Register("srate", func(s *machine.Stack) { s.Push(srate) })
	c.Machine.Register("m2f", func(s *machine.Stack) {
		s.Push(440 * math.Pow(2, (s.Pop()-69)/12))
	})
	c.Machine.Register("TET", func(s *machine.Stack) {
		scale := s.Pop()
		s.Push(440 * math.Pow(2, (s.Pop()-69)/scale))
	})
}

func chk(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
