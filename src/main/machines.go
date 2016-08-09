package main

import (
	"log"
	"encoding/json"
	"os"
)

type Machine struct {
	ID         string
	PublicIP   string
	Metadata   interface{}
	Version    string
	Hostname   string

	Interfaces []Iface
}

func SetHostname(one_machine *Machine) {
	var err error
	one_machine.Hostname, err = os.Hostname()
	if err != nil {
		log.Println(err)
	}
}

func ConstructMachine(ch chan <- Machine, node EtcdNode, full bool) {
	log.Printf("%s %s", FuncNameF(ConstructMachine), node.Value)

	var one_machine Machine
	for _, n := range node.Nodes {
		log.Printf("%s %s", FuncNameF(ConstructMachine), n.Value)
		ret := json.Unmarshal([]byte(n.Value), &one_machine)
		if ret != nil {
			log.Println(ret)
		}
		if full == true {
			SetHostname(&one_machine)
			RemoteIfaces(one_machine.PublicIP, &one_machine.Interfaces)
		}
		ch <- one_machine
	}
}

func GetMachines(full bool) []Machine {
	var machines []Machine
	var reply EtcdReply

	content := Fetch(CONF.EtcdAddress + CONF.FleetUrl + "/?recursive=true")
	ret := json.Unmarshal(content, &reply)
	if ret != nil {
		log.Println(ret)
		return machines
	}

	nb_nodes := len(reply.Node.Nodes)
	log.Printf("%s %d machines", FuncNameF(GetMachines), nb_nodes)
	if nb_nodes > 0 {
		ch := make(chan Machine)
		for i, node := range reply.Node.Nodes {
			log.Printf("%s starting %d/%d", FuncNameF(GetMachines), i + 1, nb_nodes)
			go ConstructMachine(ch, node, full)
		}
		for i := range reply.Node.Nodes {
			log.Printf("%s waiting %d/%d", FuncNameF(GetMachines), i + 1, nb_nodes)
			machines = append(machines, <-ch)
			log.Printf("%s finished %d/%d", FuncNameF(GetMachines), i + 1, nb_nodes)
		}
	}
	log.Printf("%s return [%d]Machine for %d expected", FuncNameF(GetMachines), len(machines), nb_nodes)

	return machines
}
