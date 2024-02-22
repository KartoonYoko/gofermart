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

// AddUser добавит пользователя и вернёт его ID
func (r *authRepository) AddUser(ctx context.Context, login string, password string) (int64, error) {
	query := `
		INSERT INTO users (login, password, loyality_balance_current, loyality_balance_withdrawn)
		VALUES ($1, $2, 0, 0)
		RETURNING id
	`
	res, err := r.conn.ExecContext(ctx, query, login, password)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}
