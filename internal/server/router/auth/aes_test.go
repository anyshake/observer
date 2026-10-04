package auth

import (
	"bytes"
	"testing"
)

func TestAES256GCMDecrypt(t *testing.T) {
	t.Parallel()

	secret := []byte("login-test-secret")
	encryptor := newAES256GCM(secret)
	if encryptor == nil {
		t.Fatal("newAES256GCM() returned nil")
	}
	nonce := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	aad := []byte("session-id")
	plaintext := []byte(`{"username":"observer","password":"test-password"}`)
	ciphertext := encryptor.gcm.Seal(nil, nonce, plaintext, aad)
	packet := append(bytes.Clone(nonce), ciphertext...)
	tampered := bytes.Clone(packet)
	tampered[len(tampered)-1] ^= 1
	changedNonce := bytes.Clone(packet)
	changedNonce[0] ^= 1
	tests := []struct {
		name    string
		secret  []byte
		data    []byte
		aad     []byte
		wantErr bool
	}{
		{"valid", secret, packet, aad, false},
		{"wrong key", []byte("wrong-secret"), packet, aad, true},
		{"wrong session", secret, packet, []byte("another-session"), true},
		{"missing session", secret, packet, nil, true},
		{"tampered tag", secret, tampered, aad, true},
		{"tampered nonce", secret, changedNonce, aad, true},
		{"empty", secret, nil, aad, true},
		{"short nonce", secret, nonce[:len(nonce)-1], aad, true},
		{"nonce only", secret, nonce, aad, true},
		{"truncated tag", secret, packet[:len(packet)-1], aad, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decryptor := newAES256GCM(tt.secret)
			if decryptor == nil {
				t.Fatal("newAES256GCM() returned nil")
			}
			got, err := decryptor.decrypt(tt.data, tt.aad)
			if (err != nil) != tt.wantErr {
				t.Fatalf("decrypt() error = %v, want error = %v", err, tt.wantErr)
			}
			if !tt.wantErr && !bytes.Equal(got, plaintext) {
				t.Fatalf("decrypt() = %q, want %q", got, plaintext)
			}
			if tt.wantErr && len(got) != 0 {
				t.Fatalf("failed decryption returned plaintext: %q", got)
			}
		})
	}
}
