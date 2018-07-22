package main

import (
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
		func(w http.ResponseWriter, req *http.Request) {
			fmt.Fprintf(w, "hello %v", req.TLS.PeerCertificates[0].EmailAddresses[0])
			fmt.Fprintln(w)
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

		greeting, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("%s", greeting)
	})
}
