package orderaccrual

import (
	"context"
	"gofermart/internal/model/auth"
	model "gofermart/internal/model/order_accrual"

	"github.com/jmoiron/sqlx"
)

type orderAccrualRepository struct {
	conn *sqlx.DB
}

func New(ctx context.Context, db *sqlx.DB) (*orderAccrualRepository, error) {
	repo := &orderAccrualRepository{
		conn: db,
	}

	return repo, nil
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
func (r *orderAccrualRepository) AccrualToOrderAndUser(ctx context.Context, orderID int64, userID auth.UserID, sum int) error {
	r.conn.
}
