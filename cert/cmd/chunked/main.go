package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"time"
)

var cli = flag.Bool("cli", false, "client")

func main() {
	flag.Parse()
	log.SetFlags(0)

	if *cli {
		log.Fatalf("client: %v", client())
	}
	log.Fatalf("server: %v", serve())
}

func client() error {
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 10*time.Second)

	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8000/", nil)
	req = req.WithContext(ctx)

	var client http.Client
	resp, err := client.Do(req)
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
	mux := http.NewServeMux()
	mux.HandleFunc("/", chunked)
	mux.HandleFunc("/worker.html", handleFile("worker.html"))
	mux.HandleFunc("/worker.js", handleFile("worker.js"))

	mux.HandleFunc("/debug/pprof/", pprof.Index)

	mux.HandleFunc("/favicon.ico", http.NotFound)
	mux.HandleFunc("/favicon.png", http.NotFound)
	mux.HandleFunc("/opensearch.xml", http.NotFound)

	server := http.Server{
		Addr:         "localhost:8000",
		Handler:      mux,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 12 * time.Second,
	}
	return server.ListenAndServe()
}

func handleFile(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, name)
	}
}

func chunked(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		panic("not a http.Flusher")
	}

	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	for {
		select {
		case <-time.After(1 * time.Second):
			fmt.Fprintf(w, "It is now %s\n", time.Now().UTC())
			flusher.Flush()

		case <-r.Context().Done():
			log.Printf("server: %v", r.Context().Err())
			return
		}
	}
}
