package main

import (
	"time"
	"strings"
)

type ConnectStatus struct {
	IPv4    string
	Reach   bool
	Latency int
}

func GetAllIPv4(machines *[]Machine) []string {
	var all_ips []string

	for _, m := range *machines {
		for _, i := range m.Interfaces {
			all_ips = append(all_ips, i.IPv4)
		}
	}
	return all_ips
}

func IsWantedPrefix(ip string) bool {
	for _, prefix := range CONF.prefix {
		if strings.Contains(ip, prefix) {
			return true
		}
	}
	return false

}

func RemoteTangle(m *Machine, d QueryData) {

	var conn ConnectStatus

	local_ifaces := LocalIfaces()
	all_ips := GetAllIPv4(&d.machines)

	var ips []string

	for _, local := range local_ifaces {

		for _, ip := range all_ips {
			if ip != local.IPv4 && IsWantedPrefix(ip) {
				ips = append(ips, ip)
			}
		}
	}

	for _, ip := range ips {
		conn.IPv4 = ip
		start := time.Now().Nanosecond() / 1000
		conn.Reach = IsAlive(ip)
		conn.Latency = time.Now().Nanosecond() / 1000 - start
		m.Connections = append(m.Connections, conn)
	}
}
