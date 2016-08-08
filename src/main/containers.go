package main

import (
	"log"
	"encoding/json"
	"net/http"
)

type Rkt struct {
	Name string
	Aci  string
	Pid  int
	Ip   string
}

func retrive_local() Rkt {
	one := Rkt{Name: "One", Aci: "one.aci", Pid: 10, Ip: "1.1.1.1"}

	return one
}

func get_containers(w http.ResponseWriter) {
	var containers []Rkt
	var machines []Machine

	containers = append(containers, retrive_local())
	machines = get_machines()
	for _, m := range machines {
		log.Println(m.PublicIP)
	}

	b, err_j := json.Marshal(containers)
	if err_j != nil {
		log.Println(err_j)
		return
	}
	_, err := w.Write(b)
	if err != nil {
		log.Println(err)
	}
}
