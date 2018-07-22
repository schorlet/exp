package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"
)

func ReadTLSCert(cn, password string) (tls.Certificate, error) {
	key, err := ReadKey(cn, password)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("read key: %v", err)
	}

	path := filepath.Join(PKI_PATH, cn+".crt")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return tls.Certificate{}, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return tls.Certificate{}, fmt.Errorf("no pem data found")
	}

	if block.Type != "CERTIFICATE" {
		return tls.Certificate{}, fmt.Errorf("invalid type: %v", block.Type)
	}

	leaf, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("decode cert: %v", err)
	}

	return tls.Certificate{
		Certificate: [][]byte{block.Bytes},
		PrivateKey:  key,
		Leaf:        leaf,
	}, nil
}

func NewTLSConfig(ca, server string) (*tls.Config, error) {
	clientCAs, err := NewCertPool(ca, false)
	if err != nil {
		return nil, fmt.Errorf("clientCAs pool: %v", err)
	}

	tlsServer, err := ReadTLSCert(server, "")
	if err != nil {
		return nil, fmt.Errorf("read server cert: %v", err)
	}

	// https://blog.gopheracademy.com/advent-2016/exposing-go-on-the-internet/
	return &tls.Config{
		PreferServerCipherSuites: true,
		SessionTicketsDisabled:   true,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
		Certificates: []tls.Certificate{tlsServer},
		NextProtos:   []string{"h2"},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCAs,
	}, nil
}

func NewTLSClient(ca, client string) (*http.Client, error) {
	rootCAs, err := NewCertPool(ca, true)
	if err != nil {
		return nil, fmt.Errorf("rootCAs pool: %v", err)
	}

	tlsClient, err := ReadTLSCert(client, "")
	if err != nil {
		return nil, fmt.Errorf("read server cert: %v", err)
	}

	return &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			IdleConnTimeout:       60 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			//
			TLSHandshakeTimeout: 3 * time.Second,
			TLSClientConfig: &tls.Config{
				RootCAs:                  rootCAs,
				MinVersion:               tls.VersionTLS12,
				SessionTicketsDisabled:   true,
				PreferServerCipherSuites: true,
				Certificates:             []tls.Certificate{tlsClient},
			},
		},
	}, nil
}

func NewCertPool(ca string, system bool) (*x509.CertPool, error) {
	path := filepath.Join(PKI_PATH, ca+".crt")

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read ca cert: %v", err)
	}

	var pool *x509.CertPool
	if system {
		pool, _ = x509.SystemCertPool()
	}
	if pool == nil {
		pool = x509.NewCertPool()
	}

	if !pool.AppendCertsFromPEM(data) {
		return nil, fmt.Errorf("append ca cert to pool")
	}

	return pool, nil
}

func NewTLSServer(ca, server, addr string, handler http.Handler) (*http.Server, error) {
	tlsConfig, err := NewTLSConfig("ca", "localhost")
	if err != nil {
		return nil, fmt.Errorf("create tls config: %v", err)
	}

	return &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
		TLSConfig:    tlsConfig,
	}, nil
}
