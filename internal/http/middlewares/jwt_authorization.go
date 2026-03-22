// internal/http/middlewares/jwt_auth.go
package middlewares

import (
	"strings"
	"interslavic/internal/auth"
	"github.com/gofiber/fiber/v2"
)

func JWTAuthMiddleware(jwtCfg *auth.JWTConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() == fiber.MethodOptions {
			return c.Next()
		}

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization header format",
			})
		}

		tokenString := parts[1]

		claims, err := jwtCfg.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("login", claims.Login)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

func AdminOnlyMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok || role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "admin access required",
			})
		}
		return c.Next()
	}
}

func OptionalJWTAuthMiddleware(jwtCfg *auth.JWTConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Next()
		}

		tokenString := parts[1]

		claims, err := jwtCfg.ValidateToken(tokenString)
		if err == nil {
			c.Locals("user_id", claims.UserID)
			c.Locals("login", claims.Login)
			c.Locals("role", claims.Role)
		}

		return c.Next()
	}
}