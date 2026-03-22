package http

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"interslavic/config"

	"github.com/gofiber/contrib/swagger"
	"go.uber.org/fx"
)

// New returns new *fiber.App with fx lifecycle hooks
func New(cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		//ReadTimeout:           10 * time.Second,
		DisableStartupMessage:        true,
		DisablePreParseMultipartForm: true,
		BodyLimit:                    4 * 1024 * 1024,
	})

	// swagger
	swagCfg := swagger.Config{
		BasePath: "/",
		FilePath: "./docs/swagger.json",
		Path:     "swagger",
		Title:    "Swagger API Docs",
	}

	app.Use(swagger.New(swagCfg))

	// cors
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     fmt.Sprintf("http://127.0.0.1:%s", cfg.HTTP.Port),
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders:     "Accept, Authorization, Content-Type, X-CSRF-Token",
		ExposeHeaders:    "Content-Type",
	}))

	return app
}

// Listen starts *fiber.App and appends fx.Lifecycle hooks for graceful shutdown
func Listen(lc fx.Lifecycle, app *fiber.App, cfg *config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := app.Listen(":" + cfg.HTTP.Port); err != nil {
					log.Fatal(err)
				}
			}()
			return nil
		},

		OnStop: func(ctx context.Context) error {
			// Shutdown gracefully shuts down the server without interrupting any active connections.
			// Shutdown works by first closing all open listeners and then
			// waits indefinitely for all connections to return to idle before shutting down.
			return app.ShutdownWithContext(ctx)
		},
	})
}

// Module provides fiber.App and adds lifecycle hooks for graceful shutdown
var Module = fx.Module("fiber",
	fx.Provide(
		New,
	),
	fx.Invoke(
		Listen,
	),

	APIModule,
)
