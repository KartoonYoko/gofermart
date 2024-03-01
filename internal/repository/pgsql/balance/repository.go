package balance

import (
	"context"
	"gofermart/internal/model/auth"
	modelBalance "gofermart/internal/model/balance"

	"github.com/jmoiron/sqlx"
)

type repositoryBalance struct {
	conn *sqlx.DB
}

func New(ctx context.Context, db *sqlx.DB) *repositoryBalance {
	repo := &repositoryBalance{
		conn: db,
	}

	return repo
}

func (r *repositoryBalance) GetUserBalance(ctx context.Context, userID auth.UserID) (*modelBalance.GetUserBalanceModel, error) {
	query := `SELECT loyality_balance_current, loyality_balance_withdrawn FROM users WHERE id=$1`
	res := &modelBalance.GetUserBalanceModel{}
	err := r.conn.GetContext(ctx, res, query, userID)
	if err != nil {
		return nil, err
	}

	return res, nil
}
