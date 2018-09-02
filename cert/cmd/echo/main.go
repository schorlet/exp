package main

// https://posener.github.io/http2/#full-duplex-communication

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/http2"
)

var (
	url         = "https://localhost:8443/echo"
	httpVersion = flag.Int("version", 2, "HTTP version")
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	// Create a pipe - an object that implements `io.Reader` and `io.Writer`.
	// Whatever is written to the writer part will be read by the reader part.
	pr, pw := io.Pipe()

	// Create an `http.Request` and set its body as the reader part of the
	// pipe - after sending the request, whatever will be written to the pipe,
	// will be sent as the request body.
	// This makes the request content dynamic, so we don't need to define it
	// before sending the request.
	req, err := http.NewRequest(http.MethodPut, url, ioutil.NopCloser(pr))
	if err != nil {
		log.Fatalf("create http request: %v", err)
	}

	// Create TLS config
	data, err := ioutil.ReadFile("/tmp/ca.crt")
	if err != nil {
		log.Fatalf("read ca cert: %v", err)
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(data)

	tlsConfig := tls.Config{
		RootCAs: pool,
	}

	// Use the proper transport in the client
	var client http.Client
	switch *httpVersion {
	case 1:
		client.Transport = &http.Transport{
			TLSClientConfig: &tlsConfig,
		}
	case 2:
		client.Transport = &http2.Transport{
			TLSClientConfig: &tlsConfig,
		}
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("send request: %v", err)
	}
	log.Printf("response status: %s", resp.Status)

	// Run a loop which writes every second to the writer part of the pipe
	// the current time.
	go func() {
		for {
			time.Sleep(1 * time.Second)
			fmt.Fprintf(pw, "It is now %v\n", time.Now().UTC())
		}
	}()

	// Copy the server's response to stdout.
	// Will fail on server WriteTimeout.
	_, err = io.Copy(os.Stdout, resp.Body)
	log.Fatalf("reading from response body: %v", err)
}
