package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"time"
)

func CreateServerCert(pkiPath, ca, server string, validity time.Duration) error {
	tlsCA, err := ReadTLSCert(pkiPath, ca, ca)
	if err != nil {
		return fmt.Errorf("read ca cert: %v", err)
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("generate key: %v", err)
	}

	subjectKeyId := sha1.Sum([]byte("CN=" + server + ",O=washingmachine,ST=france,C=EU"))
	ipAddresses, _ := net.LookupIP(server)

	cert := x509.Certificate{
		SerialNumber: big.NewInt(11000),
		SubjectKeyId: subjectKeyId[:],
		Subject: pkix.Name{
			CommonName:   server,
			Organization: []string{"washingmachine"},
			Province:     []string{"france"},
			Country:      []string{"EU"},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(validity),
		KeyUsage:    x509.KeyUsageKeyAgreement | x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		//
		DNSNames:    []string{server},
		IPAddresses: ipAddresses,
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

	err = SaveKey(pkiPath, server, "", key)
	if err != nil {
		return fmt.Errorf("save key: %v", err)
	}

	err = SaveCertBlock(pkiPath, server, &block)
	if err != nil {
		return fmt.Errorf("save cert: %v", err)
	}

	return nil
}
