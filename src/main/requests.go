package main

import (
	"log"
	"io/ioutil"
	"net/http"
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

func fetch(url string) []byte {
	log.Printf("GET %s ...\n", url)
	r, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	b, err_read := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err_read != nil {
		log.Println(err_read)
		return nil
	}
	return b
}
