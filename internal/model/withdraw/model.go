package withdraw

import "gofermart/internal/model/auth"

type AddUserWithdrawModel struct {
	UserID  auth.UserID `db:"user_id"`
	OrderID int64       `db:"order_id"`
	Sum     int         `db:"sum"`
}

type GetUserWithdrawModel struct {
	UserID      auth.UserID `db:"user_id" json:"-"`
	OrderID     int64       `db:"order_id" json:"order_id"`
	ProcessedAt string      `db:"processed_at" json:"processed_at"`
	Sum         int         `db:"sum" json:"sum"`
}
