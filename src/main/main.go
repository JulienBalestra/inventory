package main

import (
	"net/http"
	"log"
	"encoding/json"
)

var CONF = CreateConfig()

func NotFound(w http.ResponseWriter, path string) {

	w.WriteHeader(404)
	b, j_error := json.Marshal(CONF.Urls)
	if j_error != nil {
		log.Println(j_error)
	}
	r, err := w.Write(b)
	if err != nil {
		log.Println(err)
	}

	log.Printf("%d GET %s: 404\n", r, path)
}

func HandlerNotFound(w http.ResponseWriter, r *http.Request) {
	if (r.Method == "GET") {
		NotFound(w, r.URL.Path)
	}
}

func HandlerMachines(w http.ResponseWriter, r *http.Request) {
	if (r.Method == "GET") {
		log.Printf("%s GET %s", FuncNameF(HandlerMachines), CONF.Urls.Machines)
		MarshalAndSend(w, GetMachines(false))
	}
}

func HandlerInterfaces(w http.ResponseWriter, r *http.Request) {
	if (r.Method == "GET") {
		log.Printf("%s GET %s", FuncNameF(HandlerInterfaces), CONF.Urls.Interfaces)
		MarshalAndSend(w, LocalIfaces())
	}
}

func HanderRoot(w http.ResponseWriter, r *http.Request) {
	if (r.Method == "GET") {
		log.Printf("%s GET %s", FuncNameF(HanderRoot), CONF.Urls.Root)

		FullMachine := GetMachines(true)
		MarshalAndSend(w, FullMachine)
	}
}

func HandlerProbe(w http.ResponseWriter, r *http.Request) {
	if (r.Method == "GET") {
		log.Printf("%s GET %s", FuncNameF(HandlerProbe), CONF.Urls.Probe)
		w.Write([]byte("{\"probe\": true}"))
	}
}

func main() {
	b, _ := json.Marshal(CONF)
	log.Printf("%s", string(b))
	http.HandleFunc("/", HandlerNotFound)
	http.HandleFunc(CONF.Urls.Root, HanderRoot)
	http.HandleFunc(CONF.Urls.Machines, HandlerMachines)
	http.HandleFunc(CONF.Urls.Interfaces, HandlerInterfaces)
	http.HandleFunc(CONF.Urls.Probe, HandlerProbe)
	http.ListenAndServe(CONF.Bind + ":" + CONF.Port, nil)
}
