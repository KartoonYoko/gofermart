package orderaccrual

import (
	"context"
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
		SELECT * FROM orders
		WHERE status!='INVALID' AND status!='PROCESSED'
	`

	orders := []model.GetOrderModel{}
	err := r.conn.SelectContext(ctx, &orders, query)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
