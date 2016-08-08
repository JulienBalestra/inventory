package main

import (
	"log"
	"io/ioutil"
	"net/http"
	"encoding/json"
)

type EtcdReply struct {
	Node   EtcdNode
	Nodes  []EtcdNode
	Action string
}

type EtcdNode struct {
	CreatedIndex  int
	ModifiedIndex int
	Value         string
	Key           string

	Nodes         []EtcdNode
	Dir           bool
}

func Fetch(url string) []byte {
	var b []byte

	log.Printf("GET %s ...\n", url)
	r, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return b
	}
	b, err_read := ioutil.ReadAll(r.Body)
	r.Body.Close()
	//r.StatusCode
	if err_read != nil {
		log.Println(err_read)
		return b
	}
	return b
}

func MarshalAndSend(w http.ResponseWriter, i interface{}) {
	b, err_marshal := json.Marshal(i)
	if err_marshal != nil {
		log.Println(err_marshal)
		return
	}
	_, err := w.Write(b)
	if err != nil {
		log.Println(err)
	}
}
