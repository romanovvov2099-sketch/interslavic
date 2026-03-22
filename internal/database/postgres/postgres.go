package postgres

import (
	"context"
	"database/sql"
	"interslavic/config"
	"interslavic/logging"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/fx"
)

func NewPostgres(cfg *config.Config, logger *logging.ModuleLogger) *sql.DB {
	mlogger := logging.NewModuleLogger("DB", "POSTGRES", logger)
	mlogger.Info("Provide DB Postgres")
	db, err := sql.Open("postgres", cfg.Database.URL)
	if err != nil {
		mlogger.Error(err.Error())
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}

type sqlTransaction struct {
	tx     *sql.Tx
	active bool
}

func (t *sqlTransaction) IsActive() bool {
	return t.active
}

func (t *sqlTransaction) Commit(ctx context.Context) error {
	if err := t.tx.Commit(); err != nil {
		return err
	}
	t.active = false
	return nil
}

func (t *sqlTransaction) Rollback(ctx context.Context) error {
	if err := t.tx.Rollback(); err != nil {
		return err
	}
	t.active = false
	return nil
}

// PostgresModule
var PostgresModule = fx.Module("postgres",
	fx.Provide(
		NewPostgres,
		NewUserRepository,
		NewCourseRepository,
		NewLessonRepository,
		NewTaskRepository,
		NewLessonProgressRepository,
	),
)
