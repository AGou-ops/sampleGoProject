package main

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

func ping2(addr string) (ok bool) {
	// 方法一
	// Command := fmt.Sprintf("ping -c 1 %s  > /dev/null && echo true || echo false", addr)
	// output, err := exec.Command("/bin/bash", "-c", Command).Output()
	// if err != nil {
	// 	os.Exit(-1)
	// }
	// if string(output) == "true" {
	// 	return true
	// }
	// return false

	// 方法二
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	Command := fmt.Sprintf("ping -c 1 -t 1 %s  > /dev/null", addr)
	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", Command)
	err := cmd.Run()
	return err == nil
}

const (
	IPNET string = "192.168.20"
)

func main() {
	wg := sync.WaitGroup{}
	ch := make(chan bool, 10)

	startTime := time.Now()

	// START
	i := 1
	tmp := 0
	end := 100
	for {
		if i >= end {
			break
		}
		addr := fmt.Sprintf(IPNET + "." + strconv.Itoa(i))
		i++

		ch <- true
		wg.Add(1)
		go func(a string) {
			defer wg.Done()
			defer func() {
				<-ch
			}()

			ok := ping2(a)
			if ok {
				fmt.Printf("[+] %s\n", a)
				tmp++
				return
			}
			fmt.Printf("[-] %s\n", a)
		}(addr)
	}
	wg.Wait()

	spendTime := time.Since(startTime)
	fmt.Println("耗时：", spendTime)
	fmt.Printf("%s 汇总（%d个）", strings.Repeat("-", 10), end)
	fmt.Printf(`
	存活：%d个
	死亡：%d个
	`, tmp, end-tmp)
}
