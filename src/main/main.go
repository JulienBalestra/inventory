package main

import (
	"net/http"
	"log"
)

func not_found(w http.ResponseWriter, path string) {
	resp := []byte("404 not found\n")

	w.WriteHeader(404)
	r, err := w.Write(resp)
	if err != nil {
		log.Printf("error: %v\n", err)
		panic(err)
	}

	log.Printf("%d GET %s: 404\n", r, path)
}

type Root struct {
	Interfaces []Iface
	Machines   []Machine
}

func get_method(w http.ResponseWriter, path string) {
	root_url := "/"
	containers_url := "/containers"
	unit_url := "/units"
	machines_url := "/machines"
	interfaces := "/interfaces"

	if (path == root_url) {
		var root_data Root
		log.Printf("GET %s\n", root_url)
		//get_containers(w, true)
		root_data.Machines = get_machines()
		root_data.Interfaces = get_interfaces(root_data.Machines)
		marshal_send(w, root_data)

	} else if (path == containers_url) {
		log.Printf("GET %s\n", containers_url)
		//get_containers(w, false)

	} else if (path == unit_url) {
		log.Printf("GET %s\n", unit_url)

	} else if (path == interfaces) {
		log.Printf("GET %s\n", interfaces)
		marshal_send(w, get_interfaces(nil))

	} else if (path == machines_url) {
		log.Printf("GET %s\n", machines_url)
		m := get_machines()
		marshal_send(w, m)

	} else {
		not_found(w, path)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if (r.Method == "GET") {
		get_method(w, r.URL.Path)
	}
}

var PORT = "8080"

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe("0.0.0.0:" + PORT, nil)
}
