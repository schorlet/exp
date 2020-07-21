package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

func CreateCACert(cn string) error {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("generate key: %v", err)
	}

	subjectKeyId := sha1.Sum([]byte("CN=" + cn + ",O=washingmachine,ST=france,C=EU"))

	cert := x509.Certificate{
		SerialNumber:          big.NewInt(10000),
		IsCA:                  true,
		MaxPathLenZero:        true,
		BasicConstraintsValid: true,
		SubjectKeyId:          subjectKeyId[:],
		Subject: pkix.Name{
			CommonName:   cn,
			Organization: []string{"washingmachine"},
			Province:     []string{"france"},
			Country:      []string{"EU"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(24 * time.Hour),
		KeyUsage:  x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
	}

	der, err := x509.CreateCertificate(
		rand.Reader,
		&cert,
		&cert,
		&key.PublicKey,
		key,
	)
	if err != nil {
		return fmt.Errorf("encode cert: %v", err)
	}
	block := pem.Block{Type: "CERTIFICATE", Bytes: der}

	err = SaveKey(cn, cn, key)
	if err != nil {
		return fmt.Errorf("save key: %v", err)
	}

	err = SaveCertBlock(cn, &block)
	if err != nil {
		return fmt.Errorf("save cert: %v", err)
	}

	return nil
}

func SaveCertBlock(cn string, block *pem.Block) error {
	path := filepath.Join(*pkiPath, cn+".crt")
	log.Printf("writing %q certificate to %q\n", cn, path)

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	if err = pem.Encode(file, block); err != nil {
		return fmt.Errorf("write cert: %v", err)
	}

	if err = file.Close(); err != nil {
		return err
	}

	return nil
}
