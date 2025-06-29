package postgresql

import (
	"context"
	"database/sql"
	"fmt"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db}
}

func runInTx(ctx context.Context, db *sql.DB, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
