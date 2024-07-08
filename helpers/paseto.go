package helpers

import (
	"api-gateway/models"
	"os"
	"time"

	"aidanwoods.dev/go-paseto"
)

// GenerateKey generates a new asymmetric key pair and returns the private and public keys.
func GenerateKey() (privateKey, publicKey string) {
	secretKey := paseto.NewV4AsymmetricSecretKey() // don't share this!!!
	publicKey = secretKey.Public().ExportHex()     // DO share this one
	privateKey = secretKey.ExportHex()
	return privateKey, publicKey
}

// EncodeWithStruct encodes data into a PASETO token with additional refresh token.
func EncodeWithStruct(data *models.UserRegister) (string, string, error) {
	token := paseto.NewToken()
	now := time.Now()
	expiration := now.Add(2 * time.Hour)

	// Set claims for access token
	token.SetIssuedAt(now)
	token.SetNotBefore(now)
	token.SetExpiration(expiration)
	token.SetString("id", data.Username)

	jsonData := map[string]interface{}{
		"status":   data.Status,
		"sub":      data.Subscription,
		"exp":      data.ExpSubs,
		"username": data.Username,
	}

	err := token.Set("data", jsonData)
	if err != nil {
		return "", "", err
	}

	secretKey, err := paseto.NewV4AsymmetricSecretKeyFromHex(os.Getenv("PRIVATE_KEY_PASSETO"))
	if err != nil {
		return "", "", err
	}

	// Sign the access token
	accessToken := token.V4Sign(secretKey, nil)

	// Generate refresh token
	refreshToken := paseto.NewToken()
	refreshToken.SetIssuedAt(now)
	refreshToken.SetNotBefore(now)
	refreshToken.SetExpiration(expiration.Add(24 * time.Hour)) // Refresh token expires in 24 hours

	// Optionally, you can set additional claims for refresh token here

	// Sign the refresh token
	refreshTokenString := refreshToken.V4Sign(secretKey, nil)

	return accessToken, refreshTokenString, nil
}
