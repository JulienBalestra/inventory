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
}

func LocalIfaces() []Iface {
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
	log.Printf("%s with %d ifaces", FuncNameF(LocalIfaces), len(ifaces))
	return ifaces
}

func RemoteIfaces(m *Machine, d QueryData) {
	var content []byte

	log.Printf("%s %s ...", FuncNameF(RemoteIfaces), m.PublicIP)

	content, err := Fetch(SelfRequest(m.PublicIP, CONF.Urls.Interfaces))
	if err != nil {
		log.Printf("%s error %v", FuncNameF(GetMachines), err)
		return
	}

	if content == nil {
		log.Printf("%s with empty content: %s", FuncNameF(RemoteIfaces),
			SelfRequest(m.PublicIP, CONF.Urls.Interfaces))
		return
	}
	ret := json.Unmarshal(content, &m.Interfaces)
	if ret != nil {
		log.Println(ret)
		log.Printf("%s ERROR %s with NO ifaces", FuncNameF(RemoteIfaces), m.PublicIP)

		return
	}
	log.Printf("%s %s with %d ifaces", FuncNameF(RemoteIfaces), m.PublicIP, len(m.Interfaces))
}