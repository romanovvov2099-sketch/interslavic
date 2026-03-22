package middlewares

import (
	"interslavic/config"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// NewMiddleware return request logging handler
func NewAdminAuthorizationMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() == http.MethodOptions {
			return c.Next()
		}

		key := c.Request().Header.Peek("Authorization")
		if string(key) != cfg.HTTP.AuthorizationKey {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		return c.Next()
	}
}
