package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	err := CreateCerts("ca", "localhost", "client")
	if err != nil {
		panic(fmt.Sprintf("generate certs: %v", err))
	}
}

func withServer(fn func(string)) {
	handler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			world := "world"
			if len(r.TLS.PeerCertificates) > 0 {
				if len(r.TLS.PeerCertificates[0].EmailAddresses) > 0 {
					world = r.TLS.PeerCertificates[0].EmailAddresses[0]
				}
			}
			fmt.Fprintf(w, "hello %v", world)
		},
	)

	server := httptest.NewUnstartedServer(handler)
	defer server.Close()

	tlsConfig, err := NewTLSConfig("ca", "localhost")
	if err != nil {
		panic(fmt.Sprintf("create tls config: %v", err))
	}

	server.TLS = tlsConfig
	server.StartTLS()

	fn(server.URL)
}

func TestClientAuth(t *testing.T) {
	withServer(func(url string) {
		client, err := NewTLSClient("ca", "client")
		if err != nil {
			t.Fatalf("create client: %v", err)
		}

		res, err := client.Get(url)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		greeting := string(data)

		if greeting != "hello client@washingmachine" {
			t.Fatalf("Unexpected greeting: %q", greeting)
		}
	})
}

func TestClientNoAuth(t *testing.T) {
	withServer(func(url string) {
		client, err := NewTLSClient("ca", "client")
		if err != nil {
			t.Fatalf("create client: %v", err)
		}

		switch v := client.Transport.(type) {
		case *http.Transport:
			v.TLSClientConfig.Certificates = []tls.Certificate{}
		}

		res, err := client.Get(url)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		greeting := string(data)

		if greeting != "hello world" {
			t.Fatalf("Unexpected greeting: %q", greeting)
		}
	})
}
