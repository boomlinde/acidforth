package main

import (
	"code.google.com/p/portaudio-go/portaudio"
	"github.com/boomlinde/gobassline/collection"
	"github.com/boomlinde/gobassline/machine"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const SFREQ = 44100

func main() {
	log.Println("Booting")
	ui := NewUI("./www")
	col := collection.NewCollection()
	addComponents(col)

	data, err := ioutil.ReadAll(os.Stdin)
	chk(err)
	chk(col.Machine.Compile(machine.TokenizeBytes(data)))
	log.Println("Program loaded and compiled")

	portaudio.Initialize()
	defer portaudio.Terminate()
	stream, err := portaudio.OpenDefaultStream(0, 1, SFREQ, 0, col.Callback)
	chk(err)
	defer stream.Close()
	stream.Start()
	log.Fatal((&http.Server{
		Addr:           ":8000",
		Handler:        ui,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}).ListenAndServe())
}
