package cert

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"time"
)

var (
	pkiPath  = path.Dir(os.Args[0])
	validity = 1 * time.Minute
)

func init() {
	log.SetFlags(log.Lshortfile)
	// log.SetOutput(os.Stderr)
	log.SetOutput(ioutil.Discard)

	err := CreateCerts(pkiPath, "ca", "localhost", "client", validity)
	if err != nil {
		log.Fatalf("create certs: %v", err)
	}
}

func withServer(fn func(string)) {
	server := httptest.NewUnstartedServer(HelloHandler("world"))
	defer server.Close()

	tlsConfig, err := NewTLSConfig(pkiPath, "ca", "localhost")
	if err != nil {
		log.Fatalf("create tls config: %v", err)
	}

	server.TLS = tlsConfig
	server.StartTLS()

	fn(server.URL)
}

func TestClientAuth(t *testing.T) {
	withServer(func(url string) {
		client, err := NewTLSClient(pkiPath, "ca", "client")
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
			t.Fatalf("unexpected greeting: %q", greeting)
		}
	})
}

func TestClientNoAuth(t *testing.T) {
	withServer(func(url string) {
		client, err := NewTLSClient(pkiPath, "ca", "client")
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
			t.Fatalf("unexpected greeting: %q", greeting)
		}
	})
}
