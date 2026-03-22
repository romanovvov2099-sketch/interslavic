package middlewares

import (
	"encoding/json"
	"fmt"
	"interslavic/logging"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type RequestInfo struct {
	Route   string              `json:"route"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
}

type ResponseInfo struct {
	Status int    `json:"status"`
	Body   string `json:"body"`
}

// NewMiddleware return request logging handler
func NewLoggingMiddleware(logger *logging.ModuleLogger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() == http.MethodOptions {
			return c.Next()
		}

		reqInfoJson, err := json.Marshal(&RequestInfo{
			Route:   c.OriginalURL(),
			Headers: c.GetReqHeaders(),
			Body:    string(c.Body()),
		})
		if err != nil {
			logger.Error(fmt.Sprintf("error convert json req info: %s", err))
			return c.Next()
		}

		logger.Info("REQUEST " + string(reqInfoJson))

		return c.Next()
	}
}
