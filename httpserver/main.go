package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var host string
var port uint

func main() {
	flag.StringVar(&host, "host", "0.0.0.0", "bind host")
	flag.UintVar(&port, "port", 8080, "listening port number")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		date := time.Now().Format("2006-01-02 15:04:05.000")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("%s | %7s | %s | error: %s\n", date, r.Method, r.URL.Path, err)
			return
		}
		defer r.Body.Close()
		fmt.Printf("%s | %7s | %s | %s\n", date, r.Method, r.RequestURI, string(body))
		_, _ = w.Write([]byte("success"))
	})

	addr := fmt.Sprintf("%s:%d", host, port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
