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

	if len(os.Args) < 2 {
		panic("에러: 2개 미만의 argument")
	}
	//	programName := os.Args[0]
	firstArg := os.Args[1]

	var newipset []string
	var oldipset []string
	var k int = 0 //for test
	var m int = 0 // for test

	//	url := "https://app.rainforestqa.com/api/1/vm_stack"
	//testipset, _ := getIPsfromHTTP(url)
	testipset := []string{"107.170.58.26",
		"136.243.34.202", "136.243.61.144", "144.76.106.134", "144.76.184.203", "144.76.222.134", "144.76.39.177",
		"144.76.45.197", "144.76.79.99", "144.76.84.173", "144.76.93.110", "148.251.0.38", "148.251.13.134", "148.251.151.9",
		"148.251.81.139", "159.203.36.96", "176.9.104.78", "176.9.111.175", "176.9.120.98", "176.9.122.86", "176.9.126.173",
		"176.9.142.197", "176.9.144.139", "176.9.148.174", "176.9.150.147", "176.9.151.195", "176.9.156.166", "176.9.18.90",
		"176.9.19.67", "176.9.2.102", "176.9.2.45", "176.9.26.108", "176.9.26.118", "176.9.28.167", "176.9.30.239", "176.9.32.6",
		"176.9.36.105", "176.9.37.10", "176.9.6.149", "176.9.66.71", "176.9.80.103", "176.9.91.135", "178.63.100.16", "178.63.101.7",
		"178.63.14.136", "178.63.22.133", "178.63.70.204", "178.63.83.83", "178.63.9.77", "188.40.127.2", "188.40.134.143", "188.40.32.78",
		"188.40.76.201", "188.40.88.202", "188.40.89.135", "188.40.98.213", "213.239.199.226", "213.239.204.102", "213.239.204.92", "46.4.100.243",
		"46.4.103.15", "46.4.104.74", "46.4.104.76", "46.4.108.17", "46.4.108.83", "46.4.112.228", "46.4.113.72", "46.4.114.46", "46.4.116.168", "46.4.119.176",
		"46.4.119.233", "46.4.122.179", "46.4.123.169", "46.4.21.5", "46.4.23.70", "46.4.37.81", "46.4.40.2", "46.4.63.139", "46.4.66.90", "46.4.69.58", "46.4.79.171",
		"46.4.84.211", "46.4.85.214", "46.4.85.90", "46.4.92.231", "46.4.94.111", "46.4.94.240", "46.4.96.233", "46.4.99.4", "5.9.101.208", "5.9.101.43", "5.9.120.226",
		"5.9.122.81", "5.9.13.153", "5.9.143.37", "5.9.25.67", "5.9.31.130", "5.9.32.241", "5.9.36.100", "5.9.37.172", "5.9.38.134", "5.9.6.68", "5.9.6.73", "5.9.69.166",
		"5.9.70.148", "5.9.79.226", "5.9.80.202", "5.9.85.179", "5.9.93.73", "52.62.170.126", "78.46.128.183", "78.46.35.218", "78.46.36.58", "78.46.38.131", "78.46.38.150",
		"78.46.40.164", "78.46.64.100", "78.46.83.240", "78.47.60.200", "88.198.107.43", "88.198.111.101", "88.198.206.41", "88.198.221.230", "88.198.26.199", "88.198.53.231",
		"88.198.55.178", "88.198.90.99"}
	if firstArg == "nom" {
		//newipset, _ = getIPsfromHTTP(url)
		path := os.Getenv("HOME") + "/addrbook.json"
		newipset, _ = getIPsFromFile(path)

		for {
			//  compare
			res := reflect.DeepEqual(newipset, oldipset)
			if res == false {
				go sendIPS(newipset, string(os.Getenv("MACH1_ADDR")))
				log.Println(len(newipset), " sent")
			} else {
				//log.Println("ipset not changed")
			}
			oldipset = newipset
			if k%30 == 0 {
				oldipset = nil
				newipset = append(newipset, testipset[m])
				k = 0
				m++
			}
			//sleep
			time.Sleep(2000 * time.Millisecond)
			k++
		}
	}

	if firstArg == "test" {
		sendIPS(testipset, string(os.Getenv("MACH1_ADDR")))
		log.Println("test:", len(testipset), " sent")
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
