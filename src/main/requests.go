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

func Fetch(url string) ([]byte, error) {
	var b []byte

	log.Printf("%s GET %s ...", FuncNameF(Fetch), url)

	r, err := http.DefaultClient.Get(url)
	if err != nil {
		log.Println(err)
		return b, err
	}
	b, err_read := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err_read != nil {
		log.Println(err_read)
		return b, err_read
	}
	return b, nil
}

func MarshalAndSend(w http.ResponseWriter, i interface{}) {
	b, err_marshal := json.Marshal(i)
	if err_marshal != nil {
		log.Println(err_marshal)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(b)
	if err != nil {
		log.Println(err)
	}
}
