package main

import (
	"net/http"
	"log"
	"encoding/json"
)

var CONF = create_config()

func not_found(w http.ResponseWriter, path string) {

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

func get_method(w http.ResponseWriter, path string) {
	switch  {
	case path == CONF.Urls.Root || path == CONF.Urls.Root + "/":
		log.Printf("GET %s\n", path)

		var root_data Root

		root_data.Machines = get_machines()
		root_data.Interfaces = get_interfaces(root_data.Machines)
		marshal_send(w, root_data)
	case path == CONF.Urls.Interfaces || path == CONF.Urls.Interfaces + "/":
		log.Printf("GET %s\n", path)

		marshal_send(w, get_interfaces(nil))

	case path == CONF.Urls.Machines || path == CONF.Urls.Machines + "/":
		log.Printf("GET %s\n", path)

		marshal_send(w, get_machines())
	default:
		not_found(w, path)
	}

}

func handler(w http.ResponseWriter, r *http.Request) {
	if (r.Method == "GET") {
		get_method(w, r.URL.Path)
	}
}

func main() {
	b, _ := json.Marshal(CONF)
	log.Printf("%s", string(b))
	http.HandleFunc("/", handler)
	http.ListenAndServe(CONF.Bind + ":" + CONF.Port, nil)
}
