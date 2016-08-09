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
	return ifaces
}

func RemoteAddIfaces(ip string, interfaces *[]Iface) {
	var content []byte
	var remote_host []Iface

	content = Fetch(InternalRequest(ip, CONF.Urls.Interfaces))

	if content == nil {
		log.Printf("remote_interfaces with empty content: %s",
			InternalRequest(ip, CONF.Urls.Interfaces))
		return
	}
	ret := json.Unmarshal(content, &remote_host)
	if ret != nil {
		log.Println(ret)
		return
	}
	for _, i := range remote_host {
		*interfaces = append(*interfaces, i)
	}
}

func GetInterfaces(machines []Machine) []Iface {
	var ifaces []Iface

	if machines != nil {
		for _, m := range machines {
			RemoteAddIfaces(m.PublicIP, &ifaces)
		}

	} else {
		ifaces = LocalIfaces()
	}

	return ifaces
}
