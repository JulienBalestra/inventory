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
	log.Printf("%s with %d ifaces", FuncName(LocalIfaces), len(ifaces))
	return ifaces
}

func RemoteAddIfaces(ip string, interfaces *[]Iface) {
	var content []byte

	log.Printf("%s %s ...", FuncNameF(RemoteAddIfaces), ip)
	content = Fetch(InternalRequest(ip, CONF.Urls.Interfaces))

	if content == nil {
		log.Printf("%s with empty content: %s", FuncNameF(RemoteAddIfaces),
			InternalRequest(ip, CONF.Urls.Interfaces))
		return
	}
	ret := json.Unmarshal(content, interfaces)
	if ret != nil {
		log.Println(ret)
		return
	}
	log.Printf("%s %s with %d ifaces", FuncNameF(RemoteAddIfaces), ip, len(*interfaces))
}

func GetInterfaces(machines []Machine) []Iface {
	var ifaces []Iface

	log.Printf("%s machine number: %d", FuncNameF(GetInterfaces), len(machines))
	if machines != nil {
		for _, m := range machines {
			RemoteAddIfaces(m.PublicIP, &ifaces)
		}

	}
	return ifaces
}
