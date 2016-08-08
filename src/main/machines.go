package main

import (
	"log"
	"encoding/json"
	"os"
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
			log.Printf("browse_nodes: %s\n", node.Value)
			var one_machine Machine
			ret := json.Unmarshal([]byte(node.Value), &one_machine)
			if ret != nil {
				log.Println(ret)
				return
			}
			*machines = append(*machines, one_machine)
		}
	}
}

func BuildMachines(root_url string, key string, machines *[]Machine) {
	content := Fetch(root_url + key + "/?recursive=true")

	var reply EtcdReply
	ret := json.Unmarshal(content, &reply)
	if ret != nil {
		log.Println(ret)
		return
	}
	BrowseNodes(reply.Node, machines)
}

func GetMachines() []Machine {
	etcd_url := os.Getenv("ETCD_URL")
	if etcd_url == "" {
		etcd_url = "http://127.0.0.1:2379/v2/keys"
	}
	dir := os.Getenv("FLEET_DIR")

	if dir == "" {
		dir = "/_coreos.com/fleet/machines"
	}

	var machines []Machine
	BuildMachines(etcd_url, dir, &machines)
	return machines
}
