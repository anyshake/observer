package files

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestValidateToken(t *testing.T) {
	t.Parallel()

	path := "archive/station/HNZ.mseed"
	secret := []byte("download-test-secret")
	expires := time.Now().Add(24 * time.Hour).UnixMilli()
	token := generateToken(path, secret, expires)
	signature := strings.SplitN(token, ":", 2)[1]
	tamperedSignature := "0" + signature[1:]
	if signature[0] == '0' {
		tamperedSignature = "1" + signature[1:]
	}
	tests := []struct {
		name   string
		path   string
		token  string
		secret []byte
		want   bool
	}{
		{"valid", path, token, secret, true},
		{"different file", "archive/station/HNE.mseed", token, secret, false},
		{"different secret", path, token, []byte("another-secret"), false},
		{"expired", path, generateToken(path, secret, time.Now().Add(-time.Hour).UnixMilli()), secret, false},
		{"changed expiration", path, strconv.FormatInt(expires+1, 10) + ":" + signature, secret, false},
		{"changed signature", path, strconv.FormatInt(expires, 10) + ":" + tamperedSignature, secret, false},
		{"empty", path, "", secret, false},
		{"missing signature", path, strconv.FormatInt(expires, 10) + ":", secret, false},
		{"invalid expiration", path, "invalid:" + signature, secret, false},
		{"overflowing expiration", path, "9223372036854775808:" + signature, secret, false},
		{"trailing data", path, token + ":extra", secret, false},
		{"leading whitespace", path, " " + token, secret, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateToken(tt.path, tt.token, tt.secret); got != tt.want {
				t.Fatalf("validateToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
