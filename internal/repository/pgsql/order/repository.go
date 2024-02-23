package order

import (
	"context"
	model "gofermart/internal/model/order"

	"github.com/jmoiron/sqlx"
)

type orderRepository struct {
	conn *sqlx.DB
}

func New(ctx context.Context, db *sqlx.DB) (*orderRepository, error) {
	repo := &orderRepository{
		conn: db,
	}

	return repo, nil
}

func (r *orderRepository) AddOrder(ctx context.Context, addModel *model.AddOrderModel) error {
	query := `
	INSERT INTO orders (id, status, accrual, user_id) 
	VALUES (:order_id, :status, :accrual, :user_id)
	`
	_, err := r.conn.NamedExec(query, addModel)
	if err != nil {
		// TODO обрабатывать ошибку, если номер заказа уже существует от текущего пользователя
		// TODO обрабатывать ошибку, если номер заказа уже существует от другого пользователя
		return err
	}

	return nil
}