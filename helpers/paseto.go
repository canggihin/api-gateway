package helpers

import (
	"api-gateway/models"
	"crypto/ed25519"
	"encoding/base64"
	"os"
	"time"

	"github.com/o1egl/paseto"
)

func LoadPrivateKey() (ed25519.PrivateKey, error) {
	privateKeyData, err := base64.StdEncoding.DecodeString(os.Getenv("PASETO_PRIVATE_KEY"))
	if err != nil {
		return nil, err
	}
	return ed25519.PrivateKey(privateKeyData), nil
}

func CreateToken(data models.UserRegister) (string, string, error) {

	// publicKey := os.Getenv("PASETO_PUBLIC_KEY")
	privateKey, err := LoadPrivateKey()
	if err != nil {
		return "", "", err
	}

	v2 := paseto.NewV2()

	jsonData := map[string]interface{}{
		"username":     data.Username,
		"role":         data.Status,
		"subscription": data.Subscription,
		"exp_subs":     data.ExpSubs,
	}

	now := time.Now().UTC()

	token, err := v2.Sign(privateKey, jsonData, &paseto.JSONToken{
		Audience:   "service",
		Issuer:     "api-gateway",
		Subject:    "login",
		Expiration: now.Add(5 * time.Minute),
		IssuedAt:   now,
		NotBefore:  now,
	})
	if err != nil {
		return "", "", err
	}
	refreshToken, err := v2.Sign([]byte(privateKey), jsonData, &paseto.JSONToken{
		Audience:   "service",
		Issuer:     "api-gateway",
		Subject:    "refresh",
		Expiration: now.Add(10 * time.Minute), // Refresh token valid for 24 hours
		IssuedAt:   now,
		NotBefore:  now,
	})
	if err != nil {
		return "", "", err
	}
	return token, refreshToken, nil
}
