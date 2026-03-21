package app

import (
	"interslavic/config"
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

		// fx.Provide(
		// 	// DB module
		// 	database.NewPostgres,
		// 	database.NewRepository,

		// 	// Usecases module
		// 	client.NewUseCases,
		// 	drivers.NewUseCases,
		// 	logs.NewUseCases,
		// 	license.NewUseCases,
		// 	versions.NewUseCases,
		// 	configs.NewUseCases,
		// ),

		// fx.Invoke(
		// 	cron.UpdateVersions,
		// ),

		// // http module
		// http.Module,
	)
}
