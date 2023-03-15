package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

var host string
var timeout string

func main() {
	flag.StringVar(&host, "host", "127.0.0.1", "scan host")
	flag.StringVar(&timeout, "timeout", "3s", "scan request timeout, for example: 1s、10s、1500ms")
	flag.Parse()

	t, err := time.ParseDuration(timeout)
	if err != nil {
		log.Fatal("timeout parameter error")
	}

	wg1 := &sync.WaitGroup{}
	for i := 1; i < 65536/2; i++ {
		wg1.Add(1)
		go func(port int) {
			defer wg1.Done()
			addr := fmt.Sprintf("%s:%d", host, port)
			conn, err := net.DialTimeout("tcp", addr, t)
			if err != nil {
				return
			}
			defer conn.Close()
			fmt.Printf("%s is ok\n", addr)
		}(i)
	}
	wg1.Wait()

	wg2 := &sync.WaitGroup{}
	for i := 65536 / 2; i <= 65535; i++ {
		wg2.Add(1)
		go func(port int) {
			defer wg2.Done()
			addr := fmt.Sprintf("%s:%d", host, port)
			conn, err := net.DialTimeout("tcp", addr, t)
			if err != nil {
				return
			}
			defer conn.Close()
			fmt.Printf("%s is ok\n", addr)
		}(i)
	}
	wg2.Wait()
}
