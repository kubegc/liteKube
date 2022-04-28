package token

import (
	"crypto/rand"
	"encoding/hex"
)

func Random(size int) (string, error) {
	token := make([]byte, (size+1)/2)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token)[:size], err
}
