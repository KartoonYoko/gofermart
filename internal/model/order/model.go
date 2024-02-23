package order

import "gofermart/internal/model/auth"

// возможные статусы заказа
type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

// модель добавления заказа в хранилище
type AddOrderModel struct {
	UserID  auth.UserID `db:"user_id"`
	OrderID int         `db:"order_id"`
	Status  OrderStatus `db:"status"`
	Accrual int         `db:"accrual"`
}
