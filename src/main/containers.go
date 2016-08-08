package main

import (
	"log"
	"encoding/json"
	"net/http"
)

type Rkt struct {
	Name   string
	Aci    string
	Pid    int
	Ip     string
	HostIP string
	HostID string
}

func retrive_local() Rkt {
	one := Rkt{Name: "One", Aci: "one.aci", Pid: 10, Ip: "1.1.1.1"}

	return one
}

func remote_containers(containers *[]Rkt) {
	var machines []Machine
	var content []byte
	var remotes []Rkt

	machines = get_machines()
	http.DefaultClient.Timeout = 2
	for _, m := range machines {

		content = fetch("http://" + m.PublicIP + ":8080/containers")

		if content == nil {
			log.Printf("remote_containers empty content: http://%s:8080/containers", m.PublicIP)
			continue
		}
		ret := json.Unmarshal(content, &remotes)
		if ret != nil {
			log.Println(ret)
			continue
		}
		for _, c := range remotes {
			c.HostID = m.ID
			c.HostIP = m.PublicIP
			*containers = append(*containers, c)
		}

	}
	http.DefaultClient.Timeout = 0
}

func get_containers(w http.ResponseWriter, remote bool) {
	var containers []Rkt

	if remote {
		remote_containers(&containers)
	} else {
		containers = append(containers, retrive_local())
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
