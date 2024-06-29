package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strings"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)
	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	return encodedSalt + "$" + encodedHash, nil
}

func ComparePassword(hashedPassword, password string) error {
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 2 {
		return errors.New("invalid hash")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return err
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return err
	}

	if !compareHashes(expectedHash, hash) {
		return errors.New("password mismatch")
	}

	return nil
}

func compareHashes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	result := 0
	for i := 0; i < len(a); i++ {
		result |= int(a[i] ^ b[i])
	}

	return result == 0
}
