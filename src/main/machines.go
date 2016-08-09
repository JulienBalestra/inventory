package main

import (
	"log"
	"encoding/json"
)

type Machine struct {
	ID         string
	PublicIP   string
	Metadata   interface{}
	Version    string

	Interfaces []Iface
}

type FullMachine struct {
	ID         string
	PublicIP   string
	Metadata   interface{}
	Version    string

	Interfaces []Iface
}

func BrowseNodes(node EtcdNode, machines *[]Machine, full bool) {

	for _, node := range node.Nodes {
		if node.Dir == true {
			BrowseNodes(node, machines, full)

		} else {
			log.Printf("%s: %s", FuncName(BrowseNodes), node.Value)

			var one_machine Machine
			ret := json.Unmarshal([]byte(node.Value), &one_machine)
			if ret != nil {
				log.Println(ret)
				continue
			} else {
				if full == true {
					RemoteAddIfaces(one_machine.PublicIP, &one_machine.Interfaces)
				}
				*machines = append(*machines, one_machine)
			}
		}
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
	BrowseNodes(reply.Node, &machines, full)

	return machines
}
