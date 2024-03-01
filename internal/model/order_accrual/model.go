package orderaccrual

import "gofermart/internal/model/auth"

// GetOrderAccrualFromRemoteModel - модель получения информации о начислении баллов из системы начисления баллов
type GetOrderAccrualFromRemoteModel struct {
	Order string `json:"order"`
	// Возможные значения:
	//
	//	REGISTERED — заказ зарегистрирован, но вознаграждение не рассчитано;
	// 	INVALID — заказ не принят к расчёту, и вознаграждение не будет начислено;
	// 	PROCESSING — расчёт начисления в процессе;
	// 	PROCESSED — расчёт начисления окончен;
	Status  string   `json:"status"`
	Accrual *float64 `json:"accrual"`
}

type GetOrderModel struct {
	ID int64 `db:"id"`
	// Возможные значения:
	//
	//	NEW — заказ загружен в систему, но не попал в обработку;
	// 	PROCESSING — вознаграждение за заказ рассчитывается;
	// 	INVALID — система расчёта вознаграждений отказала в расчёте;
	// 	PROCESSED — данные по заказу проверены и информация о расчёте успешно получена.
	Status string      `db:"status"`
	UserID auth.UserID `db:"user_id"`
}
