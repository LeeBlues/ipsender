package main

import (
	ejson "encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc/jsonrpc"
	"os"
	"reflect"
	"time"
)

type Args struct{ IPS []string }
type Result bool

func main() {
	var newipset []string
	var oldipset []string
	for {
		//rand.Seed(time.Now().UTC().UnixNano())
		//url := "https://app.rainforestqa.com/api/1/vm_stack"
		//newipset, _ = getIPsfromHTTP(url)
		path := os.Getenv("HOME") + "/addrbook.json"
		newipset, _ = getIPsFromFile(path)
		//  compare
		res := reflect.DeepEqual(newipset, oldipset)
		if res == false {
			go sendIPS(newipset, string(os.Getenv("MACH1_ADDR")))
		} else {
			log.Println("ipset not changed")
		}
		oldipset = newipset
		//sleep
		time.Sleep(2000 * time.Millisecond)
	}
}

func sendIPS(newipset []string, machaddr string) {
	var res Result
	client, err := net.Dial("tcp", machaddr)
	if err != nil {
		log.Println("Dial error : ", err)
		return
	}
	c := jsonrpc.NewClient(client)
	args := &Args{newipset}
	err = c.Call("IpUpdater.IpUpdate", args, &res)
	if err != nil {
		log.Println("IpUpdate error:", err)
		return
	}
}

func getIPsFromFile(path string) ([]string, bool) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("Unable to get json IPs from file. ", err.Error())
		return nil, false
	}
	//Unmarshaling
	var data []string
	err = ejson.Unmarshal(raw, &data)
	if err != nil {
		log.Println("Unable to Unmarshal IPs from json. ", err.Error())
		return nil, false
	}
	return data, true
}

func getIPsfromHTTP(url string) ([]string, bool) {
	//url := "https://app.rainforestqa.com/api/1/vm_stack"
	res, err := http.Get(url)
	if err != nil {
		log.Println("Unable to get IPs from HTTP. ", err.Error())
		return nil, false
	}
	defer res.Body.Close()

	decoder := ejson.NewDecoder(res.Body)
	var data []string
	err = decoder.Decode(&data)
	if err != nil {
		log.Println("Error on  decoding json", err.Error())
		return nil, false
	}
	return data, true
}
