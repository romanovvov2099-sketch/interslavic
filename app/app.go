package app

import (
	"interslavic/config"
	"interslavic/internal/auth"
	"interslavic/internal/database/postgres"
	"interslavic/internal/http"
	usecase "interslavic/internal/usecases"
	"interslavic/logging"

	"go.uber.org/fx"
)

func New() *fx.App {
	return fx.New(

		// base setup application
		fx.Provide(
			config.New,
			logging.ProvideFXMLoggers,
		),
		fx.Invoke(
			logging.InvokeBaseLogger,
		),

		// JWT config
		fx.Provide(func(cfg *config.Config) *auth.JWTConfig {
			return auth.NewJWTConfig(cfg.JWT.SecretKey)
		}),

		// Usecases module
		fx.Provide(
			usecase.NewAuthUseCase,
			usecase.NewCourseUseCase,
			usecase.NewLessonUseCase,
			usecase.NewTaskUseCase,
			usecase.NewProgressUseCase,
		),

		postgres.PostgresModule,
		http.Module,
	)
}
