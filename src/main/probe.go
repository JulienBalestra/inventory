package main

import (
	"log"
	"encoding/json"
)

type  ProbeResponse struct {
	Probe bool
	Hostname string
}

func Probe() []byte {
	var probe ProbeResponse
	var ret []byte

	probe.Hostname = LocalHostname()
	probe.Probe = true

	ret, err := json.Marshal(probe)
	if err != nil {
		log.Printf("%s %s", FuncNameF(Probe), err)
	}

	return ret
}

func is_alive(ip string) bool {
	log.Printf("%s %s ...", FuncNameF(is_alive), ip)
	uri := CONF.Protocol + ip + ":" + CONF.Port + CONF.Urls.Probe
	content, _ := Fetch(uri)

	var probe ProbeResponse

	json.Unmarshal(content, &probe)
	if probe.Probe == true {
		return true
	} else {
		log.Printf("%s %s DEAD", FuncNameF(is_alive), ip)
		return false
	}

}
