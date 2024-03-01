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

func New(ctx context.Context, db *sqlx.DB) *authRepository {
	repo := &authRepository{
		conn: db,
	}

	return repo
}

// AddUser добавит пользователя и вернёт его ID
//
// Может вернуть следующие ошибки:
//   - ErrLoginAlreadyExists - логин уже занят
func (r *authRepository) AddUser(ctx context.Context, login string, password string) (model.UserID, error) {
	query := `
		INSERT INTO users (login, password, loyality_balance_current, loyality_balance_withdrawn)
		VALUES ($1, $2, 0, 0)
		RETURNING id
	`
	var id int
	err := r.conn.QueryRowContext(ctx, query, login, password).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		// если такой пользователь уже существует
		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
			return -1, NewErrLoginAlreadyExists(login, err)
		}
		return -1, err
	}

	return model.UserID(id), nil
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
