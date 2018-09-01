package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var cli = flag.Bool("cli", false, "client")

func main() {
	flag.Parse()
	log.SetFlags(0)

	if *cli {
		log.Fatal(client())
	}
	log.Fatal(serve())
}

func client() error {
	resp, err := http.Get("http://localhost:8000/")
	if err != nil {
		return err
	}
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return err
		}
		log.Print(string(line))
	}
}

func serve() error {
	http.HandleFunc("/", chunked)
	return http.ListenAndServe("localhost:8000", nil)
}

func chunked(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		panic("not a http.Flusher")
	}

	w.Header().Set("Connection", "Keep-Alive")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	for {
		time.Sleep(1 * time.Second)
		fmt.Fprintf(w, "It is now %s\n", time.Now().UTC())
		flusher.Flush()
	}
}
