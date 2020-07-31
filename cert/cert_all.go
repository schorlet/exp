package cert

import (
	"fmt"
	"time"
)

func CreateCerts(pkiPath, ca, server, client string, validity time.Duration) error {
	if err := CreateCACert(pkiPath, ca, validity); err != nil {
		return fmt.Errorf("create ca cert: %v", err)
	}

	if err := CreateServerCert(pkiPath, ca, server, validity); err != nil {
		return fmt.Errorf("create server cert: %v", err)
	}

	if err := CreateClientCert(pkiPath, ca, client, validity); err != nil {
		return fmt.Errorf("create client cert: %v", err)
	}

	return nil
}
