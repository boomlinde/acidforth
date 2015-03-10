package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type UI struct {
	*http.ServeMux
	controllers map[string]float64
}

func (u *UI) registerController(name string, callback func(float64)) {
	log.Printf("Registering controller: %s", name)
	var c float64
	u.controllers[name] = c
	u.ServeMux.HandleFunc(fmt.Sprintf("/controls/%s", name), func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var input float64
			chk(json.NewDecoder(r.Body).Decode(&input))
			callback(input)
			u.controllers[name] = input
		}
		data, err := json.Marshal(u.controllers[name])
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Write(data)
	})
}

func (u *UI) registerReader(name string, callback func() float64) {
	log.Printf("Registering reader: %s", name)
	u.ServeMux.HandleFunc(fmt.Sprintf("/readers/%s", name), func(w http.ResponseWriter, r *http.Request) {
		data, _ := json.Marshal(callback())
		w.Write(data)
	})

}

func NewUI(staticDir string) *UI {
	ui := &UI{}
	ui.ServeMux = http.NewServeMux()
	ui.controllers = make(map[string]float64)
	ui.ServeMux.Handle("/", http.FileServer(http.Dir(staticDir)))
	return ui
}
