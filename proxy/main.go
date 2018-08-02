package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

var (
	target = flag.String("target", "", "target URL (local only)")
	addr   = flag.String("addr", ":8888", "listen address")
)

func main() {
	flag.Parse()
	log.SetPrefix("proxy: ")

	targetURL, err := url.Parse(*target)
	if err != nil {
		log.Fatal(err)
	}

	proxy := NewProxy(targetURL)
	err = proxy.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

type Proxy struct {
	*http.Server
	URL string
}

func NewProxy(targetURL *url.URL) *Proxy {
	rp := httputil.NewSingleHostReverseProxy(targetURL)

	rp.ModifyResponse = func(resp *http.Response) error {
		dump, err := httputil.DumpResponse(resp, false)
		if err == nil {
			log.Printf("%s", dump)
		}

		resp.Header.Set("Access-Control-Allow-Origin", "*")
		return nil
	}

	address := *addr
	if address[0] == ':' {
		address = "localhost" + address
	}

	return &Proxy{
		Server: &http.Server{
			Addr:         *addr,
			Handler:      rp,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		URL: "http://" + address,
	}
}
