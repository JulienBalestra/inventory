package main

import (
	"time"
	"strings"
	"log"
	"net/http"
	"encoding/json"

	"bytes"
	"io/ioutil"
	"math"
)

type ConnectStatus struct {
	IPv4      string
	Reach     bool
	LatencyMs float64
}

func GetSomeIPv4(machines *[]Machine, iswanted func(ip string) bool) []string {
	var all_ips []string

	for _, m := range *machines {
		for _, i := range m.Interfaces {
			if iswanted(i.IPv4) {
				all_ips = append(all_ips, i.IPv4)
			}
		}
	}
	log.Printf("%s return %d IPs", FuncNameF(GetSomeIPv4), len(all_ips))
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

func GetPostData(r *http.Request) []string {
	var given_ips []string

	p := make([]byte, r.ContentLength)
	n, err := r.Body.Read(p)
	if n == 0 {
		log.Printf("%s ContentLen(%d) Read(%d)",
			FuncNameF(NetRo), r.ContentLength, n)
	} else if int64(n) != r.ContentLength && err != nil {
		log.Println(FuncNameF(HTangle), "ERROR", err)
		return given_ips
	}

	json_error := json.Unmarshal(p, &given_ips)
	if json_error != nil {
		log.Println(FuncNameF(HTangle), "ERROR", json_error)
		return given_ips
	}
	return given_ips
}

func TimeNowMillisecond() float64 {
	return float64(time.Now().Nanosecond()) / 1000000
}

func MakeNet(ch chan ConnectStatus, ip string) {
	var conn ConnectStatus
	var reach bool

	log.Printf("%s starting", FuncNameF(MakeNet))

	conn.IPv4 = ip
	start := TimeNowMillisecond()
	reach = IsAlive(ip)
	if reach {
		conn.Reach = true
		conn.LatencyMs = math.Abs(TimeNowMillisecond() - start) // abs: sometimes get - 0.00999
	}
	ch <- conn
	log.Printf("%s finished", FuncNameF(MakeNet))
}

func NetRo(r *http.Request) []byte {

	var ips []string
	var skip bool

	for _, ip := range GetPostData(r) {
		skip = false
		for _, local := range LocalIfaces() {
			if ip == local.IPv4 {
				log.Printf("%s skip local %s", FuncNameF(NetRo), ip)
				skip = true
				break
			}
		}
		if skip == false {
			ips = append(ips, ip)
		}
	}
	log.Printf("%s to query IP[%d]", FuncNameF(NetRo), len(ips))

	ch := make(chan ConnectStatus, len(ips))
	for i, ip := range ips {
		log.Printf("%s starting %d/%d", FuncNameF(NetRo), i + 1, len(ips))
		go MakeNet(ch, ip)
	}
	var conns []ConnectStatus
	for i := range ips {
		log.Printf("%s waiting %d/%d", FuncNameF(NetRo), i + 1, len(ips))
		conns = append(conns, <-ch)
	}
	close(ch)
	ret, _ := json.Marshal(conns)

	return ret
}

func RemoteNet(m *Machine, d QueryData) {

	serial_ips, err := json.Marshal(d.all_ips)
	if err != nil {
		log.Println("ERROR", err)
	}
	uri := SelfRequest(m.PublicIP, CONF.Urls.Network)

	log.Printf("%s POST %s IP[%d] %s", FuncNameF(RemoteNet), uri, len(d.all_ips), serial_ips)
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(serial_ips))
	if err != nil {
		log.Println("ERROR", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("ERROR", err)
		return
	}
	defer resp.Body.Close()

	body, erro := ioutil.ReadAll(resp.Body)
	if erro != nil {
		log.Println("ERROR", erro)
		return
	}

	json.Unmarshal(body, &m.Connections)
}
