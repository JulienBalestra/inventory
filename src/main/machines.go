package main

import (
	"log"
	"encoding/json"
	"strings"
	"time"
)

type Machine struct {
	ID         string
	PublicIP   string
	Metadata   interface{}
	Version    string

	Hostname   string
	Interfaces []Iface
	Alive      bool
}

func MakeMachine(ch chan <- Machine, node EtcdNode, full bool) {
	log.Printf("%s %s", FuncNameF(MakeMachine), strings.TrimPrefix(node.Key, CONF.FleetUrl + "/"))

	var one_machine Machine

	if len(node.Nodes) != 1 {
		log.Printf("%s warning of node number %d != 1", FuncNameF(MakeMachine), len(node.Nodes))
	}

	for _, n := range node.Nodes {
		ret := json.Unmarshal([]byte(n.Value), &one_machine)
		if ret != nil {
			log.Println(ret)
		}
		if full == true {
			if is_alive(one_machine.PublicIP) == false {
				log.Printf("%s %s is dead", FuncNameF(MakeMachine), one_machine)
				one_machine.Alive = false
			} else {
				one_machine.Alive = true
				RemoteHostname(one_machine.PublicIP, &one_machine)
				RemoteIfaces(one_machine.PublicIP, &one_machine.Interfaces)
			}
		} else {
			one_machine.Alive = false
		}
		ch <- one_machine
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

func GetMachines(full bool) []Machine {
	var machines []Machine
	var reply EtcdReply

	content, err := Fetch(CONF.EtcdAddress + CONF.FleetUrl + "/?recursive=true")
	if err != nil {
		log.Printf("%s error %v", FuncNameF(GetMachines), err)
		return machines
	}

	ret := json.Unmarshal(content, &reply)
	if ret != nil {
		log.Println(ret)
		return machines
	}
	nb_nodes := MachineNb(reply)
	if nb_nodes > 0 {
		ch := make(chan Machine, nb_nodes)
		for i, node := range reply.Node.Nodes {
			log.Printf("%s starting %d/%d", FuncNameF(GetMachines), i + 1, nb_nodes)
			go MakeMachine(ch, node, full)
		}
		AggregateMachines(ch, &machines, nb_nodes)
	}
	log.Printf("%s return [%d/%d]Machine", FuncNameF(GetMachines), len(machines), nb_nodes)

	return machines
}
