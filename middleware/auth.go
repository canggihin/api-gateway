package middleware

import (
	"api-gateway/helpers"
	"encoding/json"
	"os"

	"aidanwoods.dev/go-paseto"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(roleParams ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("x-authorization")
		if tokenString == "" {
			return helpers.ErrorHandler(c, &helpers.UnauthorizedError{Message: "Token is required"})
		}

		var token *paseto.Token
		var pubKey paseto.V4AsymmetricPublicKey
		var err error
		var result map[string]interface{}
		pubKey, err = paseto.NewV4AsymmetricPublicKeyFromHex(os.Getenv("PUBLIC_KEY_PASSETO"))
		if err != nil {
			return helpers.ErrorHandler(c, &helpers.InternalServerError{Message: "Failed to get public key", MessageDev: err.Error()})
		}

		parser := paseto.NewParser()
		token, err = parser.ParseV4Public(pubKey, tokenString, nil)
		if err != nil {
			return helpers.ErrorHandler(c, &helpers.UnauthorizedError{Message: "Invalid Token", MessageDev: err.Error()})
		}
		if err := json.Unmarshal(token.ClaimsJSON(), &result); err != nil {
			return helpers.ErrorHandler(c, &helpers.UnauthorizedError{Message: "Invalid Token", MessageDev: err.Error()})
		}
		for _, role := range roleParams {
			if role == result["data"].(map[string]interface{})["role"] {
				return c.Next()
			}
		}
		return c.Next()
	}
}
