package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"
)

func CreateClientCert(ca, client string) error {
	tlsCA, err := ReadTLSCert(ca, ca)
	if err != nil {
		return fmt.Errorf("read ca cert: %v", err)
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("generate key: %v", err)
	}

	subjectKeyId := sha1.Sum([]byte("CN=" + client + ",O=washingmachine,ST=france,C=EU"))

	cert := x509.Certificate{
		SerialNumber: big.NewInt(12000),
		SubjectKeyId: subjectKeyId[:],
		Subject: pkix.Name{
			CommonName:   client,
			Organization: []string{"washingmachine"},
			Province:     []string{"france"},
			Country:      []string{"EU"},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(24 * time.Hour),
		KeyUsage:    x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		//
		EmailAddresses: []string{client + "@washingmachine"},
	}

	der, err := x509.CreateCertificate(
		rand.Reader,
		&cert,
		tlsCA.Leaf,
		&key.PublicKey,
		tlsCA.PrivateKey,
	)
	if err != nil {
		return fmt.Errorf("encode cert: %v", err)
	}
	block := pem.Block{Type: "CERTIFICATE", Bytes: der}

	err = SaveKey(client, "", key)
	if err != nil {
		return fmt.Errorf("save key: %v", err)
	}

	err = SaveCertBlock(client, &block)
	if err != nil {
		return fmt.Errorf("save cert: %v", err)
	}

	return nil
}
