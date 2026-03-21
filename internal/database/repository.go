package database

import (
	"context"
	"database/sql"
	"interslavic/logging"
)

// [SERVICE INTERFACES]

type TrManager interface {
	Begin(context.Context) (context.Context, error)
	Commit(context.Context) error
	Rollback(context.Context) error
	GetCurrentTr(context.Context) *sql.Tx // WARNING: NOT UNIQUE FOR ALL DRIVERS!
}

type Repository struct {
}

func NewRepository(db *sql.DB, logger *logging.ModuleLogger) *Repository {
	return &Repository{}
}
