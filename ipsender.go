package main

import (
	ejson "encoding/json"
	"fmt"
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
		url := "https://app.rainforestqa.com/api/1/vm_stack"
		newipset, _ = getIPsfromHTTP(url)
		//  compare
		res := reflect.DeepEqual(newipset, oldipset)
		if res == false {
			go sendIPS(newipset, string(os.Getenv("MACH1_ADDR")))
			go sendIPS(newipset, string(os.Getenv("MACH2_ADDR")))
			/*
				for i := 1; i < 3; i++ { //temporaily hardcoded
					s := "MACH" + strconv.Itoa(i) + "_ADDR"
					client, err := net.Dial("tcp", string(os.Getenv(s)))
					if err != nil {
						fmt.Println("error : ", err)
					}
					c := jsonrpc.NewClient(client)
					dummy := &Args{nil}
					err = c.Call("IpUpdater.IpUpdateInit", dummy, &res)
					if err != nil {
						log.Fatal("IpUpdateInit error:", err)
					}
					args := &Args{newipset}
					err = c.Call("IpUpdater.IpUpdate", args, &res)
					if err != nil {
						log.Fatal("IpUpdate error:", err)
					}
				}
			*/
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
	//for i := 1; i < 3; i++ { //temporaily hardcoded
	//s := "MACH" + strconv.Itoa(i) + "_ADDR"
	client, err := net.Dial("tcp", machaddr)
	if err != nil {
		fmt.Println("error : ", err)
	}
	c := jsonrpc.NewClient(client)
	dummy := &Args{nil}
	err = c.Call("IpUpdater.IpUpdateInit", dummy, &res)
	if err != nil {
		log.Fatal("IpUpdateInit error:", err)
	}
	args := &Args{newipset}
	err = c.Call("IpUpdater.IpUpdate", args, &res)
	if err != nil {
		log.Fatal("IpUpdate error:", err)
	}
	//}
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
