package main

import (
	"code.google.com/p/portaudio-go/portaudio"
	"github.com/boomlinde/acidforth/collection"
	"github.com/boomlinde/acidforth/machine"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const SFREQ = 44100

func main() {
	log.Println("Booting")
	col := collection.NewCollection()
	addComponents(SFREQ, col)

	log.Println("Waiting for source on stdin")
	data, err := ioutil.ReadAll(os.Stdin)
	chk(err)

	log.Println("Tokenizing source")
	tokens := machine.TokenizeBytes(data)

	log.Println("Expanding macros")
	tokens, err = machine.ExpandMacros(tokens)
	chk(err)

	log.Println("Parsing")
	chk(col.Machine.Compile(tokens))
	log.Println("Running")

	portaudio.Initialize()
	defer portaudio.Terminate()
	stream, err := portaudio.OpenDefaultStream(0, 1, SFREQ, 0, col.Callback)
	chk(err)
	defer stream.Close()
	stream.Start()

	for {
		time.Sleep(time.Second)
	}
}
