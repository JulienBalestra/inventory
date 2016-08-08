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

func get_method(w http.ResponseWriter, path string) {
	root_url := "/"
	rkt_url := "/containers"
	unit_url := "/units"
	machines_url := "/machines"

	if (path == root_url) {
		log.Printf("GET %s\n", root_url)
		get_containers(w, true)
	} else if (path == rkt_url) {
		log.Printf("GET %s\n", rkt_url)
		get_containers(w, false)
	} else if (path == unit_url) {
		log.Printf("GET %s\n", unit_url)
	} else if (path == machines_url) {
		log.Printf("GET %s\n", machines_url)
		m := get_machines()
		send_machines(w, m)

	} else {
		not_found(w, path)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if (r.Method == "GET") {
		get_method(w, r.URL.Path)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe("0.0.0.0:8080", nil)
}
