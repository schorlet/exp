package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

var PKI_PATH = os.TempDir()

func main() {
	// if err := CreateCerts("ca", "localhost", "client"); err != nil {
	// log.Fatalf("generate certs: %v", err)
	// }

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
		log.Fatalf("create tls server: %v", err)
	}

	fmt.Println("Starting server ...")
	if err := tlsServer.ListenAndServeTLS("", ""); err != nil {
		log.Fatal(err)
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
