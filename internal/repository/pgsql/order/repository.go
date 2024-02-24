package order

import (
	"context"
	"errors"
	model "gofermart/internal/model/order"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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
	_, err := r.conn.NamedExecContext(ctx, query, addModel)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
			// проверим кому принадлежит заказ, чтобы вернуть верную ошибку
			type GetOrderModel struct {
				ID     int `db:"id"`
				UserID int `db:"user_id"`
			}
			query = `
			SELECT id, user_id FROM orders WHERE id=$1
			`
			var getOrderModel GetOrderModel
			err = r.conn.GetContext(ctx, &getOrderModel, query, addModel.OrderID)
			if err != nil {
				return err
			}

			if getOrderModel.UserID == int(addModel.UserID) {
				return ErrOrderAlreadyExists
			} else {
				return ErrOrderBelongsToAnotherUser
			}
		}

		return err
	}

	return nil
}
