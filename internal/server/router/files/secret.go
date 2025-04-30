package files

import (
	"crypto/rand"
)

func getSecretKey(length int) ([]byte, error) {
	secretKey := make([]byte, length)
	if _, err := rand.Read(secretKey); err != nil {
		return nil, err
	}

	return secretKey, nil
}
