package storage

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

func RandomString(length int) (string, error) {
	// https://www.commandlinefu.com/commands/view/24071/generate-random-text-based-on-length
	raw := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, raw)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(raw)[:length], nil
}
