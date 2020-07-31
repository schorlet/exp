package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/schorlet/exp/cert"
)

var (
	pkiPath  = flag.String("pki", os.TempDir(), "path to read/write certificates and keys")
	validity = flag.Duration("validity", 24*time.Hour, "validity lifetime in hours")
	verbose  = flag.Bool("v", false, "print log messages")
)

func init() {
	flag.Parse()

	log.SetOutput(ioutil.Discard)
	if *verbose {
		log.SetFlags(0)
		log.SetPrefix(os.Args[0] + ": ")
		log.SetOutput(os.Stderr)
	}
}

func main() {
	if err := cert.CreateCerts(*pkiPath, "ca", "localhost", "client", *validity); err != nil {
		log.Fatalf("create certs: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", cert.HelloHandler("world"))
	mux.HandleFunc("/echo", cert.EchoHandler)

	tlsServer, err := cert.NewTLSServer(*pkiPath, "ca", "localhost", ":8443", mux)
	if err != nil {
		log.Fatalf("create tls server: %v", err)
	}

	log.Println("starting server ...")
	if err := tlsServer.ListenAndServeTLS("", ""); err != nil {
		log.Fatal(err)
	}
}
