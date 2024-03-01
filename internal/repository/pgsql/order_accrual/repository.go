package orderaccrual

import (
	"context"
	"errors"
	"gofermart/internal/model/auth"
	modelOrder "gofermart/internal/model/order"
	model "gofermart/internal/model/order_accrual"

	"github.com/jmoiron/sqlx"
)

type orderAccrualRepository struct {
	conn *sqlx.DB
}

func New(ctx context.Context, db *sqlx.DB) *orderAccrualRepository {
	repo := &orderAccrualRepository{
		conn: db,
	}

	return repo
}

func (r *orderAccrualRepository) GetUnhandledOrders(ctx context.Context) ([]model.GetOrderModel, error) {
	query := `
		SELECT id, status, user_id FROM orders
		WHERE status!='INVALID' AND status!='PROCESSED'
	`

	orders := []model.GetOrderModel{}
	err := r.conn.SelectContext(ctx, &orders, query)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

// AccrualToOrderAndUser - начислит баллы в заказ, также прибавит эти баллы пользователю
func (r *orderAccrualRepository) AccrualToOrderAndUser(ctx context.Context, orderID int64, userID auth.UserID, sum int, orderStatus modelOrder.OrderStatus) error {
	if !orderStatus.Valid() {
		return errors.New("not valid order status")
	}

	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	type getUserModel struct {
		ID                     auth.UserID
		CurrentLoyalityBalance int
	}
	var user getUserModel
	query := `SELECT id, loyality_balance_current FROM users WHERE id=$1`
	err = tx.QueryRowContext(ctx, query, userID).Scan(&user.ID, &user.CurrentLoyalityBalance)
	if err != nil {
		return err
	}

	query = `UPDATE users SET loyality_balance_current = $1 WHERE id=$2`
	_, err = tx.ExecContext(ctx, query, user.CurrentLoyalityBalance+sum, user.ID)
	if err != nil {
		return err
	}

	query = `UPDATE orders SET status=$1, accrual=$2 WHERE id=$3`
	_, err = tx.ExecContext(ctx, query, orderStatus, sum, orderID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
