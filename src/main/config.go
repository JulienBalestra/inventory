package main

import (
	"net"
	"runtime"
	"reflect"
	"time"
	"strings"
)

type Config struct {
	Urls              Urls
	Port              string
	Protocol          string
	Bind              string

	EtcdAddress       string
	FleetMachineUrl   string

	LogPadding        int
	HttpClientTimeout time.Duration
	GoRoutineTimeout  time.Duration
	GoRoutineSleep    time.Duration

	prefix            []string
}

type Urls struct {
	Version    string
	Api        string
	Root       string

	Machines   string
	Interfaces string
	Hostname   string
	Tangle     string

	Probe      string
	Help       string
	Ui         string
}

func CreateConfig() Config {
	var c Config

	c.Port = "5000"
	c.Protocol = "http://"
	valid_ip := net.ParseIP("0.0.0.0")
	c.Bind = valid_ip.String()
	c.HttpClientTimeout = time.Millisecond * 500
	c.GoRoutineSleep = time.Millisecond * 100
	c.GoRoutineTimeout = time.Second * 5

	// Common use of Etcd and Fleet
	c.EtcdAddress = "http://127.0.0.1:2379/v2/keys"
	c.FleetMachineUrl = "/_coreos.com/fleet/machines"

	c.LogPadding = 15

	c.prefix = append(c.prefix, "192.168")
	c.prefix = append(c.prefix, "10.1.")


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

	// /api/v0/hostname
	c.Urls.Hostname = c.Urls.Root + "/hostname"

	// /api/v0/tangle
	c.Urls.Tangle = c.Urls.Root + "/tangle"


	// /help
	c.Urls.Help = "/help"

	c.Urls.Ui = "/"

	return c
}

func SelfRequest(target string, url string) string {

	// http://1.1.1.1:8080/api/v0/machines
	return CONF.Protocol + target + ":" + CONF.Port + url
}

func FuncName(i interface{}) string {
	return strings.Trim(runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name(), "main.")
}

func FuncNameF(i interface{}) string {
	name := FuncName(i)
	for i := len(name); i < CONF.LogPadding; i++ {
		name = name + " "
	}
	return name
}