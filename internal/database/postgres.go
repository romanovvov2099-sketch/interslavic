package database

import (
	"context"
	"database/sql"
	"errors"
	"interslavic/config"
	"interslavic/logging"
	"time"

	_ "github.com/lib/pq"
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

	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}

// -----------[TRANSACTION MANAGER]-----------

type PostgresTrManager struct {
	db        *sql.DB
	currentTr Transaction
}

func NewPostgresTrManager(db *sql.DB) *PostgresTrManager {
	return &PostgresTrManager{
		db: db,
	}
}

// ...........[service structures]...........

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

// ...........[transaction structures]...........

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

/*
type Settings interface {
	// EnrichBy fills nil properties from external Settings.
	EnrichBy(external Settings) Settings

	// CtxKey returns trm.CtxKey for the trm.Transaction.
	CtxKey() CtxKey
	CtxKeyOrNil() *CtxKey
	SetCtxKey(*CtxKey) Settings

	// Propagation returns trm.Propagation.
		Propagation() Propagation
	PropagationOrNil() *Propagation
	SetPropagation(*Propagation) Settings

	// Cancelable defines that parent trm.Transaction can cancel child trm.Transaction or goroutines.
	Cancelable() bool
	CancelableOrNil() *bool
	SetCancelable(*bool) Settings

	// TimeoutOrNil returns time.Duration of the trm.Transaction.
	TimeoutOrNil() *time.Duration
	SetTimeout(*time.Duration) Settings
}

type СtxManager interface {
	// Default gets Transaction from context.Context by default CtxKey.
	Default(ctx context.Context) Transaction
	// SetDefault sets.Transaction in context.Context by default CtxKey.
	SetDefault(ctx context.Context, t Transaction) context.Context

	// ByKey gets Transaction from context.Context by CtxKey.
	ByKey(ctx context.Context, key CtxKey) Transaction
	// SetByKey sets Transaction in context.Context by.CtxKey.
	SetByKey(ctx context.Context, key CtxKey, t Transaction) context.Context
}
*/
