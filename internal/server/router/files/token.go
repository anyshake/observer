package files

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

func generateToken(filePath string, secretKey []byte, expireAt int64) string {
	mac := hmac.New(sha256.New, secretKey)
	message := fmt.Sprintf("%s:%d", filePath, expireAt)
	mac.Write([]byte(message))
	return fmt.Sprintf("%d:%s", expireAt, hex.EncodeToString(mac.Sum(nil)))
}

func validateToken(filePath, token string, secretKey []byte) bool {
	var expireAt int64
	var sig string
	_, err := fmt.Sscanf(token, "%d:%s", &expireAt, &sig)
	if err != nil || time.Now().Unix() > expireAt {
		return false
	}
	expected := generateToken(filePath, secretKey, expireAt)
	return hmac.Equal([]byte(token), []byte(expected))
}
