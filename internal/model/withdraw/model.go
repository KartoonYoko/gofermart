package withdraw

import "gofermart/internal/model/auth"


type AddUserWithdrawModel struct {
	UserID  auth.UserID `db:"user_id"`
	OrderID int64       `db:"order_id"`
	Sum     int         `db:"sum"`
}
