package main

import (
	ejson "encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/rpc/jsonrpc"
	"os"
	"time"
)

type Args struct{ IPS []string }
type Result bool

func main() {

	if len(os.Args) < 2 {
		panic("error: need argument")
	}

	firstArg := os.Args[1]
	rand.Seed(time.Now().UTC().UnixNano())

	testipset := []string{"107.170.58.26:12345",
		"136.243.34.202:12345", "136.243.61.144:12345", "144.76.106.134:12345", "144.76.184.203:12345", "144.76.222.134:12345", "144.76.39.177:12345",
		"144.76.45.197:12345", "144.76.79.99:12345", "144.76.84.173:12345", "144.76.93.110:12345", "148.251.0.38:12345", "148.251.13.134:12345", "148.251.151.9:12345",
		"148.251.81.139", "159.203.36.96", "176.9.104.78", "176.9.111.175", "176.9.120.98", "176.9.122.86", "176.9.126.173",
		"176.9.142.197", "176.9.144.139", "176.9.148.174", "176.9.150.147", "176.9.151.195", "176.9.156.166", "176.9.18.90",
		"176.9.19.67", "176.9.2.102", "176.9.2.45", "176.9.26.108", "176.9.26.118", "176.9.28.167", "176.9.30.239", "176.9.32.6",
		"176.9.36.105", "176.9.37.10", "176.9.6.149:12345", "176.9.66.71", "176.9.80.103", "176.9.91.135", "178.63.100.16", "178.63.101.7",
		"178.63.14.136", "178.63.22.133", "178.63.70.204", "178.63.83.83", "178.63.9.77", "188.40.127.2", "188.40.134.143", "188.40.32.78",
		"188.40.76.201", "188.40.88.202:12345", "188.40.89.135:12345", "188.40.98.213", "213.239.199.226", "213.239.204.102", "213.239.204.92", "46.4.100.243",
		"46.4.103.15", "46.4.104.74", "46.4.104.76", "46.4.108.17", "46.4.108.83", "46.4.112.228", "46.4.113.72", "46.4.114.46", "46.4.116.168", "46.4.119.176",
		"46.4.119.233", "46.4.122.179", "46.4.123.169", "46.4.21.5", "46.4.23.70:12345", "46.4.37.81", "46.4.40.2", "46.4.63.139", "46.4.66.90", "46.4.69.58", "46.4.79.171",
		"46.4.84.211", "46.4.85.214", "46.4.85.90", "46.4.92.231", "46.4.94.111", "46.4.94.240", "46.4.96.233", "46.4.99.4", "5.9.101.208", "5.9.101.43", "5.9.120.226",
		"5.9.122.81", "5.9.13.153", "5.9.143.37", "5.9.25.67:12345", "5.9.31.130:12345", "5.9.32.241", "5.9.36.100", "5.9.37.172", "5.9.38.134", "5.9.6.68", "5.9.6.73", "5.9.69.166",
		"5.9.70.148", "5.9.79.226", "5.9.80.202", "5.9.85.179", "5.9.93.73:12345", "52.62.170.126", "78.46.128.183", "78.46.35.218", "78.46.36.58", "78.46.38.131", "78.46.38.150",
		"78.46.40.164", "78.46.64.100", "78.46.83.240", "78.47.60.200", "88.198.107.43", "88.198.111.101", "88.198.206.41", "88.198.221.230", "88.198.26.199", "88.198.53.231",
		"88.198.55.178", "88.198.90.99"}

	if firstArg == "nor" {
		//newipset, _ = getIPsfromHTTP(url)
		path := os.Getenv("HOME") + "/addrbook.json"
		newipset, _ := getIPsFromFile(path)

		sendIPS(newipset, string(os.Getenv("MACH1_ADDR")))
		log.Println("normal:", len(newipset), "sent")

	}

	if firstArg == "ran" {
		var k int = 0
		for {
			if k%2 == 0 {
				//log.Println(randomipset)
				//log.Println(len(randomipset), " sent")
				var sw bool
				var randomipset []string
				for i := 1; i < len(testipset); i++ {
					rand.Seed(time.Now().UTC().UnixNano())
					r := rand.Intn(2)
					//log.Println("first r =", r)
					if r == 0 {
						rand.Seed(time.Now().UTC().UnixNano())
						r := rand.Intn(len(testipset))
						//log.Println("se r =", r)
						if r == 0 {
							r = 1
						}
						if len(randomipset) == 0 {
							randomipset = append(randomipset, testipset[r])
						}
						sw = true
						for j := 0; j < len(randomipset); j++ {
							if randomipset[j] == testipset[r] {
								sw = false
								break
							}
						}
						if sw == true {
							//log.Println("randomipset length=", len(randomipset))
							randomipset = append(randomipset, testipset[r])
						}
					}
				}
				log.Println("sending ", len(randomipset))
				time.Sleep(1 * time.Second)
				sendIPS(randomipset, string(os.Getenv("MACH1_ADDR")))
				log.Println(len(randomipset), "sent")

				k = 0
			}
			//sleep
			time.Sleep(200 * time.Millisecond)
			k++
		}
	}

	if firstArg == "all" {
		sendIPS(testipset, string(os.Getenv("MACH1_ADDR")))
		log.Println("test:", len(testipset), "sent")
	}

	if firstArg == "zero" {
		sendIPS(nil, string(os.Getenv("MACH1_ADDR")))
		log.Println("test: 0 sent")
	}
}

func Shuffle(vals []string) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	// We start at the end of the slice, inserting our random
	// values one at a time.
	for n := len(vals); n > 0; n-- {
		randIndex := r.Intn(n)
		// We swap the value at index n-1 and the random index
		// to move our randomly chosen value to the end of the
		// slice, and to move the value that was at n-1 into our
		// unshuffled portion of the slice.
		vals[n-1], vals[randIndex] = vals[randIndex], vals[n-1]
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
