package auth

import (
	"context"
	"database/sql"
	"errors"
	model "gofermart/internal/model/auth"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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
//
// Может вернуть следующие ошибки:
//   - ErrLoginAlreadyExists - логин уже занят
func (r *authRepository) AddUser(ctx context.Context, login string, password string) (int64, error) {
	query := `
		INSERT INTO users (login, password, loyality_balance_current, loyality_balance_withdrawn)
		VALUES ($1, $2, 0, 0)
		RETURNING id
	`
	res, err := r.conn.ExecContext(ctx, query, login, password)
	if err != nil {
		var pgErr *pgconn.PgError
		// если такой пользователь уже существует
		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
			return -1, NewErrLoginAlreadyExists(login, err)
		}
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}

// GetUserByLoginAndPassword вернёт информацию о пользователе по совпадению его логина и пароля
//
// Может вернуть следующие ошибки:
//   - ErrUserNotFound - пользователь не найден
func (r *authRepository) GetUserByLoginAndPassword(ctx context.Context, login string, password string) (*model.GetUserByLoginAndPasswordModel, error) {
	query := `
	SELECT id, login, password FROM users 
	WHERE login=$1 AND password=$2
	`

	user := model.GetUserByLoginAndPasswordModel{}
	err := r.conn.GetContext(ctx, &user, query, login, password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
