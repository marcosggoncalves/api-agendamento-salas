package hashing

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

var Key = []byte("fv@3:@}Gr_BN5W.EkjgYCsD99#55Z8F}2>WFKuo2vF*,c7e{t#ht+WfL>;&lbP,Q")

func GenerateToken(maxDigits int) (string, error) {
	numBytes := (maxDigits * 3) / 4
	tokenBytes := make([]byte, numBytes)

	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}

	hash := sha256.New()
	hash.Write([]byte(Key))
	hash.Write(tokenBytes)
	tokenHash := hash.Sum(nil)

	token := base64.URLEncoding.EncodeToString(tokenHash)
	if len(token) > maxDigits {
		token = token[:maxDigits]
	}

	return token, nil
}
