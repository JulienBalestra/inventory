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

func HNotFound(w http.ResponseWriter, r *http.Request) {
	NotFound(w, r.URL.Path)
}

func HMachines(w http.ResponseWriter, r *http.Request) {
	if (r.Method == "GET") {
		log.Printf("%s GET %s", FuncNameF(HMachines), CONF.Urls.Machines)
		MarshalAndSend(w, GetMachines(false))
	}
}

func HInterfaces(w http.ResponseWriter, r *http.Request) {
	if (r.Method == "GET") {
		log.Printf("%s GET %s", FuncNameF(HInterfaces), CONF.Urls.Interfaces)
		MarshalAndSend(w, LocalIfaces())
	}
}

func HRoot(w http.ResponseWriter, r *http.Request) {
	if (r.Method == "GET") {
		log.Printf("%s GET %s", FuncNameF(HRoot), CONF.Urls.Root)

		FullMachine := GetMachines(true)
		MarshalAndSend(w, FullMachine)
	}
}

func HProbe(w http.ResponseWriter, r *http.Request) {
	if (r.Method == "GET") {
		log.Printf("%s GET %s", FuncNameF(HProbe), CONF.Urls.Probe)
		p := Probe()
		w.Write(p)
	}
}

func HHostname(w http.ResponseWriter, r *http.Request) {
	if (r.Method == "GET") {
		log.Printf("%s GET %s", FuncNameF(HHostname), CONF.Urls.Hostname)
		h := []byte(LocalHostname())
		w.Write(h)
	}
}

func HTangle(w http.ResponseWriter, r *http.Request) {
	if (r.Method == "POST") {
		log.Printf("%s POST %s", FuncNameF(HTangle), CONF.Urls.Tangle)
		if r.ContentLength > 0 {
			t := Tangle(r)
			w.Write(t)
		} else {
			log.Printf("%s POST %s EMPTY", FuncNameF(HTangle), CONF.Urls.Tangle)
		}
	}
}

func main() {
	b, _ := json.Marshal(CONF)
	http.DefaultClient.Timeout = CONF.HttpClientTimeout
	log.Printf("%s", string(b))
	http.HandleFunc("/", HNotFound)
	http.HandleFunc(CONF.Urls.Root, HRoot)
	http.HandleFunc(CONF.Urls.Machines, HMachines)
	http.HandleFunc(CONF.Urls.Interfaces, HInterfaces)
	http.HandleFunc(CONF.Urls.Hostname, HHostname)
	http.HandleFunc(CONF.Urls.Probe, HProbe)
	http.HandleFunc(CONF.Urls.Tangle, HTangle)
	http.ListenAndServe(CONF.Bind + ":" + CONF.Port, nil)
}
