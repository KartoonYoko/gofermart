package auth

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type authRepository struct {
	conn *sqlx.DB
}

func New(ctx context.Context, db *sqlx.DB) (*authRepository, error) {
	repo := &authRepository{
		conn: db,
	}

	return repo, nil
}
