package withdraw

import (
	"context"
	"gofermart/internal/model/auth"
	modelWithdraw "gofermart/internal/model/withdraw"

	"github.com/jmoiron/sqlx"
)

type repositoryWithdraw struct {
	conn *sqlx.DB
}

func New(ctx context.Context, db *sqlx.DB) *repositoryWithdraw {
	repo := &repositoryWithdraw{
		conn: db,
	}

	return repo
}

// Возможные ошибки
//
//	ErrUserHasNotEnoughBalance - недостаточно средств
func (r *repositoryWithdraw) WithdrawFromUserBalance(ctx context.Context, addModel modelWithdraw.AddUserWithdrawModel) error {
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	type userData struct {
		Current  int `db:"loyality_balance_current"`
		Withdraw int `db:"loyality_balance_withdrawn"`
	}
	user := &userData{}
	query := `
	SELECT loyality_balance_current, loyality_balance_withdrawn FROM users WHERE id=$1
	`
	err = tx.QueryRowContext(ctx, query, addModel.UserID).Scan(&user.Current, &user.Withdraw)
	if err != nil {
		return err
	}

	if user.Current < addModel.Sum {
		return modelWithdraw.ErrUserHasNotEnoughBalance
	}

	query = `INSERT INTO withdrawals (order_id, user_id, sum) VALUES ($1, $2, $3)`
	_, err = tx.ExecContext(ctx, query, addModel.OrderID, addModel.UserID, addModel.Sum)
	if err != nil {
		return err
	}

	query = `
	UPDATE users
	SET loyality_balance_withdrawn=$1, loyality_balance_current=$2
	WHERE id=$3;
	`
	newWithdraw := user.Withdraw + addModel.Sum
	newCurrent := user.Current - addModel.Sum
	_, err = tx.ExecContext(ctx, query, newWithdraw, newCurrent, addModel.UserID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *repositoryWithdraw) GetUserWithdrawals(ctx context.Context, userID auth.UserID) ([]modelWithdraw.GetUserWithdrawModel, error) {
	query := `
	SELECT order_id, sum, processed_at FROM withdrawals
	WHERE user_id=$1
	ORDER BY processed_at
	`

	withdrawals := []modelWithdraw.GetUserWithdrawModel{}
	err := r.conn.SelectContext(ctx, &withdrawals, query, userID)
	if err != nil {
		return nil, err
	}

	return withdrawals, nil
}
