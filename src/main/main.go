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

type Root struct {
	Interfaces []Iface
	Machines   []Machine
}

func GetMethod(w http.ResponseWriter, path string) {
	switch  {
	case path == CONF.Urls.Root || path == CONF.Urls.Root + "/":
		log.Printf("GET %s\n", path)

		var root_data Root

		root_data.Machines = GetMachines()
		root_data.Interfaces = GetInterfaces(root_data.Machines)
		MarshalAndSend(w, root_data)
	case path == CONF.Urls.Interfaces || path == CONF.Urls.Interfaces + "/":
		log.Printf("GET %s\n", path)

		MarshalAndSend(w, GetInterfaces(nil))

	case path == CONF.Urls.Machines || path == CONF.Urls.Machines + "/":
		log.Printf("GET %s\n", path)

		MarshalAndSend(w, GetMachines())

	default:
		NotFound(w, path)
	}

}

func Handler(w http.ResponseWriter, r *http.Request) {
	if (r.Method == "GET") {
		GetMethod(w, r.URL.Path)
	}
}

func main() {
	b, _ := json.Marshal(CONF)
	log.Printf("%s", string(b))
	http.HandleFunc("/", Handler)
	http.ListenAndServe(CONF.Bind + ":" + CONF.Port, nil)
}
