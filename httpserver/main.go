package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var host string
var port uint

func main() {
	flag.StringVar(&host, "host", "0.0.0.0", "bind host")
	flag.UintVar(&port, "port", 8080, "listening port number")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("aaa")
		_, _ = w.Write([]byte("success"))
	})

	addr := fmt.Sprintf("%s:%d", host, port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
