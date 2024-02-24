package orderaccrual

type GetOrderAccrualAPIModel struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual *int   `json:"accrual"`
}

type GetOrderModel struct {
	ID     int64  `db:"id"`
	Status string `db:"status"`
}
