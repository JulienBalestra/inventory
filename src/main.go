package main

import (
	"net/http"
	"log"
	"encoding/json"
	"os"
	"io/ioutil"
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

type Rkt struct {
	Name string
	Aci  string
	Pid  int
	Ip   string
}

func get_containers(w http.ResponseWriter) {
	one := Rkt{Name: "One", Aci: "one.aci", Pid: 10, Ip: "1.1.1.1"}

	b, err_j := json.Marshal(one)
	if err_j != nil {
		log.Println(err_j)
		return
	}
	_, err := w.Write(b)
	if err != nil {
		log.Println(err)
	}

}

type EtcdReply struct {
	Node   EtcdNode
	Nodes  []EtcdNode
	Action string
}

type EtcdNode struct {
	CreatedIndex  int
	ModifiedIndex int
	Value         string
	Key           string

	Nodes         []EtcdNode
	Dir           bool
}

type Machines struct {
	Machines []Machine
}

type Machine struct {
	ID       string
	PublicIP string
	Metadata string
	Version  string
}

func fetch(url string) []byte {
	log.Printf("GET %s ...\n", url)
	r, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	b, err_read := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err_read != nil {
		log.Println(err_read)
		return nil
	}
	return b
}

func build_machines(root_url string, key string, machines *Machines) {
	content := fetch(root_url + key)

	var reply EtcdReply
	ret := json.Unmarshal(content, &reply)
	if ret != nil {
		log.Println(ret)
		return
	}
	for _, node := range reply.Node.Nodes {
		if node.Dir == true {
			log.Printf("dir == %s\n", node.Key)
			build_machines(root_url, node.Key, machines)
		} else {
			log.Printf("value == %s\n", node.Value)
			var m Machine
			ret = json.Unmarshal([]byte(node.Value), &m)
			machines.Machines = append(machines.Machines, m)
		}
	}
}

func get_machines(etcd_url string) Machines {
	dir := os.Getenv("FLEET_DIR")
	if dir == "" {
		dir = "/_coreos.com/fleet/machines"
	}

	var machines Machines
	build_machines(etcd_url, dir, &machines)
	for _, m := range machines.Machines {
		log.Printf("%s %s", m.ID, m.PublicIP)
	}
	return machines
}

func send_machines(w http.ResponseWriter, m Machines) {
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

func get_method(w http.ResponseWriter, path string) {
	root_url := "/"
	rkt_url := "/containers"
	unit_url := "/units"
	machines_url := "/machines"
	etcd_url := os.Getenv("ETCD_URL")
	if etcd_url == "" {
		etcd_url = "http://127.0.0.1:2379/v2/keys"
	}

	if (path == root_url) {
		log.Printf("GET %s\n", root_url)
	} else if (path == rkt_url) {
		log.Printf("GET %s\n", rkt_url)
		get_containers(w)
	} else if (path == unit_url) {
		log.Printf("GET %s\n", unit_url)
	} else if (path == machines_url) {
		log.Printf("GET %s\n", machines_url)
		m := get_machines(etcd_url)
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
