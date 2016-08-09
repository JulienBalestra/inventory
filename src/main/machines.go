package main

import (
	"log"
	"encoding/json"
)

type Machine struct {
	ID       string
	PublicIP string
	Metadata interface{}
	Version  string
}

func BrowseNodes(node EtcdNode, machines *[]Machine) {

	for _, node := range node.Nodes {
		if node.Dir == true {
			BrowseNodes(node, machines)

		} else {
			log.Printf("%s: %s", FuncName(BrowseNodes), node.Value)

			var one_machine Machine
			ret := json.Unmarshal([]byte(node.Value), &one_machine)
			if ret != nil {
				log.Println(ret)
				continue
			} else {
				*machines = append(*machines, one_machine)
			}
		}
	}
}

func GetMachines() []Machine {
	var machines []Machine
	var reply EtcdReply

	content := Fetch(CONF.EtcdAddress + CONF.FleetUrl + "/?recursive=true")

	ret := json.Unmarshal(content, &reply)
	if ret != nil {
		log.Println(ret)
		return machines
	}
	BrowseNodes(reply.Node, &machines)

	//RequestMachines(etcd_url, dir, &machines)
	return machines
}
