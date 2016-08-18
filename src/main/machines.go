package main

import (
	"log"
	"encoding/json"
	"strings"
	"time"
)

type Machine struct {
	ID          string
	PublicIP    string
	Metadata    interface{}
	Version     string

	Hostname    string
	Interfaces  []Iface
	Alive       bool

	Connections []ConnectStatus
}

type QueryData struct {
	machines []Machine
	reply    EtcdReply
	all_ips  []string
	fts      []func(m *Machine, re QueryData)
}

func MakeMachine(d QueryData, ch chan <- Machine, node EtcdNode) {
	log.Printf("%s %s", FuncNameF(MakeMachine),
		strings.TrimPrefix(node.Key, CONF.FleetMachineUrl + "/"))

	var m Machine

	if len(node.Nodes) != 1 {
		log.Printf("%s warning of node number %d != 1", FuncNameF(MakeMachine), len(node.Nodes))
	}

	for _, n := range node.Nodes {
		ret := json.Unmarshal([]byte(n.Value), &m)
		if ret != nil {
			log.Println(ret)
		} else if IsAlive(m.PublicIP) {
			m.Alive = true
			for _, ft := range d.fts {
				ft(&m, d)
			}
		}

		ch <- m
		break
	}
}

func closer(ch chan Machine, count chan int, size int) {
	defer close(ch)
	defer close(count)

	log.Printf("%s waiting for %.2fs or task done", FuncNameF(closer), CONF.GoRoutineTimeout.Seconds())

	end := time.Now().Add(CONF.GoRoutineTimeout)

	for time.Now().Unix() < end.Unix() {
		if len(count) >= size {
			log.Printf("%s task done", FuncNameF(closer))
			break
		}
		log.Printf("%s sleep %.1fs count(%d) < size(%d)",
			FuncNameF(closer), CONF.GoRoutineSleep.Seconds(), len(count), size)
		time.Sleep(CONF.GoRoutineSleep)
	}
	log.Printf("%s closing channel used as %d/%d",
		FuncNameF(closer), len(count), size)
	if len(ch) > 0 {
		log.Panicf("%s ERROR closing channel with content %d", FuncNameF(closer), len(ch))
	}
}

func AggregateMachines(ch_machine chan Machine, machines *[]Machine, nb_nodes int) {

	ch_count := make(chan int, nb_nodes)
	go closer(ch_machine, ch_count, nb_nodes)

	for i := 0; i < nb_nodes; i++ {
		log.Printf("%s waiting %d/%d", FuncNameF(GetMachines), i + 1, nb_nodes)
		*machines = append(*machines, <-ch_machine)
		ch_count <- i
		log.Printf("%s finished %d/%d", FuncNameF(GetMachines), i + 1, nb_nodes)
	}
}

func MachineNb(reply EtcdReply) int {
	nb_nodes := 0
	for _, node := range reply.Node.Nodes {
		if len(node.Nodes) == 1 {
			nb_nodes++
		} else {
			log.Printf("%s skipping %s", FuncNameF(MachineNb), node.Key)
		}
	}
	log.Printf("%s %d machines", FuncNameF(MachineNb), nb_nodes)
	return nb_nodes
}

func StartRoutine(d QueryData, machines *[]Machine) {
	nb_nodes := MachineNb(d.reply)

	ch_machines := make(chan Machine, nb_nodes)
	for i, node := range d.reply.Node.Nodes {
		log.Printf("%s starting %d/%d", FuncNameF(StartRoutine), i + 1, nb_nodes)
		go MakeMachine(d, ch_machines, node)
	}
	AggregateMachines(ch_machines, machines, nb_nodes)
}

func GetMachines(full bool) []Machine {
	var machines []Machine
	var d QueryData

	content, err := Fetch(CONF.EtcdAddress + CONF.FleetMachineUrl + "/?recursive=true")
	if err != nil {
		log.Printf("%s error %v", FuncNameF(GetMachines), err)
		return machines
	}

	ret := json.Unmarshal(content, &d.reply)
	if ret != nil {
		log.Println(ret)
		return machines
	}

	if full {
		d.fts = append(d.fts, RemoteIfaces)
		StartRoutine(d, &d.machines)
		log.Printf("%s Full start\n\n", FuncNameF(GetMachines))
		d.all_ips = GetSomeIPv4(&d.machines, IsWantedPrefix)
		d.fts = append(d.fts, RemoteTangle)
		d.fts = append(d.fts, RemoteHostname)
		StartRoutine(d, &machines)
	} else {
		StartRoutine(d, &machines)
	}

	log.Printf("%s return [%d]Machine", FuncNameF(GetMachines), len(machines))

	return machines
}
