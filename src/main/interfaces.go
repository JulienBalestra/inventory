package main

import (
	"net"
	"log"
	"strings"
	"strconv"
	"encoding/json"
)

type Iface struct {
	IPv4    string
	CIDR    string
	Netmask int

	HostIP  string
	HostID  string
}

func local_interfaces() []Iface {
	var ifaces []Iface
	var iface Iface

	interfaces, err := net.InterfaceAddrs()
	if err != nil {
		log.Println(err)
	}
	for _, i := range interfaces {
		ip, network, err := net.ParseCIDR(i.String())
		if err != nil {
			log.Println(err)
			continue
		}
		if !ip.IsLoopback() && ip.To4() != nil {
			iface.CIDR = network.String()
			iface.Netmask, _ = strconv.Atoi(strings.Split(network.String(), "/")[1])
			iface.IPv4 = ip.String()
			ifaces = append(ifaces, iface)
		}
	}
	return ifaces
}

func remote_interfaces(machines []Machine, interfaces *[]Iface) {
	var content []byte
	var remote_host []Iface

	for _, m := range machines {

		content = fetch(create_request(m.PublicIP, CONF.Urls.Interfaces))

		if content == nil {
			log.Printf("remote_interfaces with empty content: %s",
				create_request(m.PublicIP, CONF.Urls.Interfaces))
			continue
		}
		ret := json.Unmarshal(content, &remote_host)
		if ret != nil {
			log.Println(ret)
			continue
		}
		for _, i := range remote_host {
			i.HostID = m.ID
			i.HostIP = m.PublicIP
			*interfaces = append(*interfaces, i)
		}
	}
}

func get_interfaces(machines []Machine) []Iface {
	var ifaces []Iface

	if machines != nil {
		remote_interfaces(machines, &ifaces)
	} else {
		ifaces = local_interfaces()
	}

	return ifaces
}
