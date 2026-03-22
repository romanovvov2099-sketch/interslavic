package postgres

import (
	"context"
	"database/sql"
	"errors"
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

type PostgresTrManager struct {
	db        *sql.DB
	currentTr Transaction
}

func NewPostgresTrManager(db *sql.DB) *PostgresTrManager {
	return &PostgresTrManager{
		db: db,
	}
}

type txKey struct{}

func (m *PostgresTrManager) Begin(ctx context.Context) (context.Context, error) {
	tx, err := m.db.Begin()
	if err != nil {
		return nil, err
	}

	if m.currentTr != nil {
		return nil, errors.New("previous transaction still exists")
	}

	m.currentTr = &sqlTransaction{tx: tx, active: true}

	return context.WithValue(ctx, txKey{}, tx), nil
}

func (m *PostgresTrManager) Commit(ctx context.Context) error {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx.Commit()
	}

	m.currentTr = nil

	return nil
}

func (m *PostgresTrManager) Rollback(ctx context.Context) error {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx.Rollback()
	}

	m.currentTr = nil

	return nil
}

func (m *PostgresTrManager) GetCurrentTr(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}

	m.currentTr = nil

	return nil
}

type Transaction interface {
	IsActive() bool
	Commit(context.Context) error
	Rollback(context.Context) error
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