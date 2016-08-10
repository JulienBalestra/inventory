package main

import (
	"net"
	"runtime"
	"reflect"
)

type Config struct {
	Urls        Urls
	Port        string
	Protocol    string
	Bind        string

	EtcdAddress string
	FleetUrl    string

	LogPadding  int
}

type Urls struct {
	Version    string
	Api        string
	Root       string

	Machines   string
	Interfaces string

	Probe string
}

func CreateConfig() Config {
	var c Config

	c.Port = "8080"
	c.Protocol = "http://"
	valid_ip := net.ParseIP("0.0.0.0")
	c.Bind = valid_ip.String()

	// Common use of Etcd and Fleet
	c.EtcdAddress = "http://127.0.0.1:2379/v2/keys"
	c.FleetUrl = "/_coreos.com/fleet/machines"

	c.LogPadding = 22


	// Internal Application //
	c.Urls.Api = "/api"
	c.Urls.Version = "/v0"

	// /api/v0
	c.Urls.Root = c.Urls.Api + c.Urls.Version

	// /api/v0/machines
	c.Urls.Machines = c.Urls.Root + "/machines"

	// /api/v0/interfaces
	c.Urls.Interfaces = c.Urls.Root + "/interfaces"

	// /api/v0/probe
	c.Urls.Probe = c.Urls.Root + "/probe"

	return c
}

func InternalRequest(target string, url string) string {

	// http://1.1.1.1:8080/api/v0/machines
	return CONF.Protocol + target + ":" + CONF.Port + url
}

func FuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func FuncNameF(i interface{}) string {
	name := FuncName(i)
	for i := len(name); i < CONF.LogPadding; i++ {
		name = name + " "
	}
	return name
}