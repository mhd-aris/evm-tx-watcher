package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type UnitOfWork interface {
	WithTransaction(ctx context.Context, fn func(tx *sqlx.Tx) error) (err error)
}

type unitOfWork struct {
	db *sqlx.DB
}

func NewUnitOfWork(db *sqlx.DB) UnitOfWork {
	return &unitOfWork{db: db}
}

func (u *unitOfWork) WithTransaction(ctx context.Context, fn func(tx *sqlx.Tx) error) (err error) {
	tx, err := u.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}
