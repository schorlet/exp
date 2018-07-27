package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func ReadKey(cn, password string) (*rsa.PrivateKey, error) {
	path := filepath.Join(*pkiPath, cn+".key")
	log.Printf("reading %q key from %q\n", cn, path)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("no pem data found")
	}

	if block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("invalid type: %v", block.Type)
	}

	der := block.Bytes
	if password != "" {
		der, err = x509.DecryptPEMBlock(block, []byte(password))
		if err != nil {
			return nil, fmt.Errorf("decrypt key: %v", err)
		}
	}

	key, err := x509.ParsePKCS1PrivateKey(der)
	if err != nil {
		return nil, fmt.Errorf("decode key: %v", err)
	}

	return key, nil
}

func SaveKey(cn, password string, key *rsa.PrivateKey) error {
	block := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	if password == "" {
		return SaveKeyBlock(cn, &block)
	}

	encryptedBlock, err := x509.EncryptPEMBlock(
		rand.Reader,
		block.Type,
		block.Bytes,
		[]byte(password),
		x509.PEMCipherAES256,
	)
	if err != nil {
		return fmt.Errorf("encrypt key: %v", err)
	}

	return SaveKeyBlock(cn, encryptedBlock)
}

func SaveKeyBlock(cn string, block *pem.Block) error {
	path := filepath.Join(*pkiPath, cn+".key")
	log.Printf("writing %q key to %q\n", cn, path)

	// O_EXCL: file must not exist
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0400)
	if err != nil {
		return err
	}

	err = pem.Encode(file, block)
	if err != nil {
		return fmt.Errorf("write key: %v", err)
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}
