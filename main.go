package main

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
)

func main() {
	log.Println("port forwarding will start in 5 seconds...")
	wslIp := getWslIp()
	forwardList := autoGetList()
	autoForward(wslIp, forwardList)
}

// fetch ip in wsl2
func getWslIp() string {
	result, err := exec.Command("wsl", "ifconfig", "eth0").Output()
	if err != nil {
		log.Fatalln(err.Error())
	}
	reg := regexp.MustCompile("inet ([^ ]*)")
	match := reg.FindStringSubmatch(string(result))
	if match == nil {
		log.Fatalln("getWslIp: failed to get ip in wsl")
	}
	return match[1]
}

// get ports opened in wsl2
func autoGetList() []string {
	result, err := exec.Command("wsl", "netstat", "-tnlp").Output()
	if err != nil {
		log.Fatalln(err.Error())
	}
	reg := regexp.MustCompile("tcp.*?0\\.0\\.0\\.0\\:(\\d{1,5})")
	match := reg.FindAllStringSubmatch(string(result), -1)
	if match == nil {
		log.Fatalln("autoGetList: failed to get listend ports in wsl")
	}
	var res []string
	for _, port := range match {
		exist := false
		for _, existPort := range res {
			if port[1] == existPort {
				exist = true
			}
		}
		if !exist {
			res = append(res, port[1])
		}
	}
	return res
}

func autoForward(wslIp string, list []string) {
	for _, port := range list {
		if err := exec.Command("netsh", "interface", "portproxy", "add", "v4tov4", fmt.Sprintf("listenport=%s", port), "listenaddress=*", fmt.Sprintf("connectport=%s", port), fmt.Sprintf("connectaddress=%s", wslIp), "protocol=tcp").Run(); err != nil {
			log.Fatalln("autoForward: " + err.Error())
		} else {
			log.Println(fmt.Sprintf("tcp/%s...ok", port))
		}
	}
}
