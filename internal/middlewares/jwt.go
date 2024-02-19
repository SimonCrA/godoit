package middlewares

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
)

// middleware JWT function
func NewAuthMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" || !strings.Contains(authHeader, "Bearer") {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		// Extract token string from Authorization header
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		// Check if the token is valid
		if !token.Valid {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
		}

		// Extract claims from the token
		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to extract token claims")
		}

		// Set claims in the context for later use
		c.Locals("idUser", claims["ID"])

		return c.Next()
	}
}
