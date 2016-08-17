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

func RemoteIfaces(ip string, interfaces *[]Iface) {
	var content []byte

	log.Printf("%s %s ...", FuncNameF(RemoteIfaces), ip)

	content, err := Fetch(AppRequest(ip, CONF.Urls.Interfaces))
	if err != nil {
		log.Printf("%s error %v", FuncNameF(GetMachines), err)
		return
	}

	if content == nil {
		log.Printf("%s with empty content: %s", FuncNameF(RemoteIfaces),
			AppRequest(ip, CONF.Urls.Interfaces))
		return
	}
	ret := json.Unmarshal(content, interfaces)
	if ret != nil {
		log.Println(ret)
		return
	}
	log.Printf("%s %s with %d ifaces", FuncNameF(RemoteIfaces), ip, len(*interfaces))
}