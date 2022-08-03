package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main1() {
	start := time.Now()
	ip := "192.168.20."
	wg.Add(254)
	for i := 1; i <= 254; i++ {
		//fmt.Println(ip + strconv.Itoa(i))
		trueIP := ip + strconv.Itoa(i)
		go ping(trueIP)
	}
	wg.Wait()
	cost := time.Since(start)
	fmt.Println("执行时间:", cost)
}

func ping(ip string) {
	var beaf = "false"
	Command := fmt.Sprintf("ping -c 1 %s  > /dev/null && echo true || echo false", ip)
	output, err := exec.Command("/bin/bash", "-c", Command).Output()
	if err != nil {
		fmt.Println(err)
		return
	}
	real_ip := strings.TrimSpace(string(output))

	if real_ip == beaf {
		fmt.Printf("IP: %s  失败\n", ip)
	} else {
		fmt.Printf("IP: %s  成功 ping通\n", ip)
	}
	wg.Done()
}
