package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func MakeRefreshToken() (string, error) {
	randomData := make([]byte, 32)
	_, err := rand.Read(randomData)
	if err != nil {
		return "", fmt.Errorf(" error making a ref token  %v ", err)
	}

	encodedString := hex.EncodeToString(randomData)
	return encodedString, nil
}
