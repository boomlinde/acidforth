package main

import (
	"fmt"
	"github.com/boomlinde/acidforth/collection"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func server(col *collection.Collection, address string) {
	mux := http.NewServeMux()

	mux.HandleFunc("/compiler", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		prg, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Could not read program", http.StatusInternalServerError)
			return
		}
		err = col.Machine.Build(prg)
		if err != nil {
			http.Error(w, "Could not compile program: "+err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "OK\n")
		log.Println("Running new program from", r.RemoteAddr)
	})

	mux.HandleFunc("/playback", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if col.Playing {
			log.Print("Stopping sequencer")
			fmt.Fprintf(w, "Stopping sequencer\n")
		} else {
			log.Println("Starting sequencer")
			fmt.Fprintf(w, "Starting sequencer\n")
		}
		col.Playing = !col.Playing
	})

	ws := &http.Server{
		Addr:           address,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	ws.ListenAndServe()
}
