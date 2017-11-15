package main

import (
	ejson "encoding/json"
	"fmt"
	"log"
	//	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"reflect"
	"time"
)

type ipUpdater struct {
	client *rpc.Client
}

type Args struct {
	ip string
}

type Result bool

func (i *ipUpdater) ipupdate(args *Args, result *Result) error {
	err := i.client.Call("ipupdater", args, &result)
	return err
}

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
			log.Println("send ip set", newipset)
			client, err := net.Dial("tcp", string(os.Getenv("MACH1_ADDR")))
			if err != nil {
				fmt.Println("error : ", err)
			}
			//			ipupdater := &ipUpdater{client: client}
			//		ipupdater.ipupdate(&Args{newipset}, &response)

			//	var res Result
			c := jsonrpc.NewClient(client)
			for i := 0; i < len(newipset); i++ {
				args := &Args{newipset[i]}
				err = c.Call("IpUpdater.IpUpdate", args, &res)
				if err != nil {
					log.Fatal("arith error:", err)
				}
			}
			/*
				err = c.Call("IpUpdater.IpUpdate", args, &res)
				if err != nil {
					log.Fatal("arith error:", err)
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

/*
func RandomIP() *ipSet {
	var b, c, d, e int
	totalNum := rand.Intn(10)
	ipset := NewipSet(totalNum)

	for i := 0; i < totalNum; i++ {
		b, c, d, e = rand.Intn(500), rand.Intn(500), rand.Intn(500), rand.Intn(500)
		ipstr := strconv.Itoa(b) + "." + strconv.Itoa(c) + "." + strconv.Itoa(d) + "." + strconv.Itoa(e)
		ipset.ip[i] = ipstr
		//fmt.Println(ipstr)
	}

	return ipset
}

func main() {
	//newipset := ipSet{}
	var newipset []string
	//oldipset := ipSet{}
	var oldipset []string
	var homepath, path, cmd string
	filename := "/addrbook.json"
	for {
		rand.Seed(time.Now().UTC().UnixNano())
		//newipset = *RandomIP()
		url := "https://app.rainforestqa.com/api/1/vm_stack"
		newipset, _ = getIPsfromHTTP(url)
		//  compare
		res := reflect.DeepEqual(newipset, oldipset)
		if res == false {
			bolB, _ := ejson.Marshal(newipset)
			home, _ := exec.Command("echo", os.Getenv("HOME")).Output()
			homepath = strings.TrimSuffix(string(home), "\n")
			path = homepath + filename
			fmt.Println(path)
			err := ioutil.WriteFile(path, bolB, 0644)
			if err != nil {
				errByte := []byte(err.Error())
				ioutil.WriteFile("error.json", errByte, 0644)
			}

			// copy file to remote proxies
			dest, _ := exec.Command("echo", os.Getenv("MACH1_ADDR")).Output()
			cmd = "scp"
			homepath = homepath + filename
			result, err := exec.Command(cmd, homepath, string(dest)).Output()

			fmt.Println("dest: ", dest, "result: ", result)
			if err != nil {
				fmt.Println("error : ", err)
			}
			fmt.Println(string(bolB))
		}
		oldipset = newipset
		//sleep
		time.Sleep(2000 * time.Millisecond)
	}
}

*/
