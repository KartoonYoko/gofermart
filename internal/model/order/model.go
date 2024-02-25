package order

import "gofermart/internal/model/auth"

// возможные статусы заказа
//
//	NEW — заказ загружен в систему, но не попал в обработку;
// 	PROCESSING — вознаграждение за заказ рассчитывается;
// 	INVALID — система расчёта вознаграждений отказала в расчёте;
// 	PROCESSED — данные по заказу проверены и информация о расчёте успешно получена.
type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

func (os *OrderStatus) Valid() bool {
	return *os == "NEW" ||
		*os == "PROCESSING" ||
		*os == "INVALID" ||
		*os == "PROCESSED"
}

// модель добавления заказа в хранилище
type AddOrderModel struct {
	UserID  auth.UserID `db:"user_id"`
	OrderID int64       `db:"order_id"`
	Status  OrderStatus `db:"status"`
	Accrual int         `db:"accrual"`
}
