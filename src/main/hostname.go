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

func RemoteHostname(m *Machine, d QueryData) {

	c, err := Fetch(AppRequest(m.PublicIP, CONF.Urls.Hostname))
	if err != nil {
		log.Printf("%s error %v", FuncNameF(RemoteHostname), err)
		return
	}

	m.Hostname = string(c)
}