package storage

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// RandomString returns a random string of the specified length.
func RandomString(length int) (string, error) {
	// https://www.commandlinefu.com/commands/view/24071/generate-random-text-based-on-length
	if length < 0 {
		return "", fmt.Errorf("invalid length: %d", length)
	}
	raw := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, raw)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(raw)[:length], nil
}
