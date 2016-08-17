package main

import (
	"os"
	"log"
)

func LocalHostname() string {

	hostname, err := os.Hostname()
	if err != nil {
		log.Println(err)
	}
	log.Printf("%s %s", FuncNameF(LocalHostname), hostname)
	return hostname
}

func RemoteHostname(ip string, machine *Machine) {

	c, err := Fetch(AppRequest(ip, CONF.Urls.Hostname))
	if err != nil {
		log.Printf("%s error %v", FuncNameF(RemoteHostname), err)
		return
	}

	machine.Hostname = string(c)
}