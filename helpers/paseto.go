package helpers

import (
	"api-gateway/models"
	"os"
	"time"

	"github.com/o1egl/paseto"
)

func loadPrivateKey(privKeyFile string) ([]byte, error) {
	privateKeyBytes, err := os.ReadFile(privKeyFile)
	if err != nil {
		return nil, err
	}
	return privateKeyBytes, nil
}

func CreateToken(data models.UserRegister) (string, string, error) {

	// publicKey := os.Getenv("PASETO_PUBLIC_KEY")
	privateKey, err := loadPrivateKey("private_key.pem")
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
	refreshToken, err := v2.Sign(privateKey, jsonData, &paseto.JSONToken{
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
