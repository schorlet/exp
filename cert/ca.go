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
		Version:               3,
		SerialNumber:          big.NewInt(1),
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
		KeyUsage:  x509.KeyUsageCertSign,
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

	err = SaveCertBlock(cn, &block)
	if err != nil {
		return fmt.Errorf("save cert: %v", err)
	}

	err = SaveKey(cn, cn, key)
	if err != nil {
		return fmt.Errorf("save key: %v", err)
	}

	return nil
}

func SaveCertBlock(cn string, block *pem.Block) error {
	path := filepath.Join(PKI_PATH, cn+".crt")

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("open %s: %v", path, err)
	}

	if err = pem.Encode(file, block); err != nil {
		return fmt.Errorf("write to %s: %v", path, err)
	}

	if err = file.Close(); err != nil {
		return fmt.Errorf("close %s: %v", path, err)
	}

	fmt.Printf("%q certificate saved to %q\n", cn, path)

	return nil
}
