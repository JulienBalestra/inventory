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

	machine.Hostname = string(Fetch(AppRequest(ip, CONF.Urls.Hostname)))
}