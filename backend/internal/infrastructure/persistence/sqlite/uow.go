package sqlite

import (
	"context"
	"fmt"

	"example.com/m/internal/usecase"
	"github.com/jmoiron/sqlx"
)

type ctxKey struct{}

type Transactor struct {
	db *sqlx.DB
}

func NewTransactor(db *sqlx.DB) usecase.UnitOfWork {
	return &Transactor{db: db}
}

func (t *Transactor) Do(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	tx, err := t.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	txCtx := context.WithValue(ctx, ctxKey{}, tx)

	if err := fn(txCtx); err != nil {
		if rollBackErr := tx.Rollback(); rollBackErr != nil {
			return fmt.Errorf(
				"rollback failed: %v (original error: %w)",
				rollBackErr,
				err,
			)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// helper function to get the correct DB or Tx from context
func GetExt(ctx context.Context, db *sqlx.DB) DB {
	if tx, ok := ctx.Value(ctxKey{}).(*sqlx.Tx); ok {
		return tx
	}
	return db
}
