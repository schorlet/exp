package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

var (
	pkiPath = flag.String("pki", os.TempDir(), "path to read/write certificates and keys")
	stderr  = log.New(os.Stderr, "", log.LstdFlags)
)

func init() {
	verbose := flag.Bool("v", false, "print log messages")
	flag.Parse()

	log.SetOutput(ioutil.Discard)
	if *verbose {
		log.SetOutput(os.Stderr)
	}
}

func main() {
	stderr := log.New(os.Stderr, "", log.LstdFlags)

	if err := CreateCerts("ca", "localhost", "client"); err != nil {
		stderr.Fatalf("create certs: %v", err)
	}

	handler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			dump, err := httputil.DumpRequest(r, false)
			if err != nil {
				http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
				return
			}
			log.Printf("%s", dump)

			world := "world"
			if len(r.TLS.PeerCertificates) > 0 {
				if len(r.TLS.PeerCertificates[0].EmailAddresses) > 0 {
					world = r.TLS.PeerCertificates[0].EmailAddresses[0]
				}
			}
			fmt.Fprintf(w, "hello %v", world)
			fmt.Fprintln(w)
		},
	)

	tlsServer, err := NewTLSServer("ca", "localhost", ":8443", handler)
	if err != nil {
		stderr.Fatalf("create tls server: %v", err)
	}

	log.Println("Starting server ...")
	if err := tlsServer.ListenAndServeTLS("", ""); err != nil {
		stderr.Fatal(err)
	}
}

func CreateCerts(ca, server, client string) error {
	if err := CreateCACert(ca); err != nil {
		return fmt.Errorf("create ca cert: %v", err)
	}

	if err := CreateServerCert(ca, server); err != nil {
		return fmt.Errorf("create server cert: %v", err)
	}

	if err := CreateClientCert(ca, client); err != nil {
		return fmt.Errorf("create client cert: %v", err)
	}

	return nil
}
