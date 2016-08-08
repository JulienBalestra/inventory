package main

import (
	"log"
	"encoding/json"
	"net/http"
	"os"
)

type Machine struct {
	ID       string
	PublicIP string
	Metadata interface{}
	Version  string
}

func browse_nodes(node EtcdNode, machines *[]Machine) {
	for _, node := range node.Nodes {
		if node.Dir == true {
			browse_nodes(node, machines)
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

func build_machines(root_url string, key string, machines *[]Machine) {
	content := fetch(root_url + key + "/?recursive=true")

	var reply EtcdReply
	ret := json.Unmarshal(content, &reply)
	if ret != nil {
		log.Println(ret)
		return
	}
	browse_nodes(reply.Node, machines)
}

func get_machines() []Machine {
	etcd_url := os.Getenv("ETCD_URL")
	if etcd_url == "" {
		etcd_url = "http://127.0.0.1:2379/v2/keys"
	}
	dir := os.Getenv("FLEET_DIR")

	if dir == "" {
		dir = "/_coreos.com/fleet/machines"
	}

	var machines []Machine
	build_machines(etcd_url, dir, &machines)
	return machines
}

func send_machines(w http.ResponseWriter, m []Machine) {
	b, err_marshal := json.Marshal(m)
	if err_marshal != nil {
		log.Println(err_marshal)
		return
	}
	_, err := w.Write(b)
	if err != nil {
		log.Println(err)
	}
}