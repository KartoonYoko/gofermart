package balance

type GetUserBalanceModel struct {
	Current   int `db:"loyality_balance_current"`
	Withdrawn int `db:"loyality_balance_withdrawn"`
}

type GetUserBalanceAPIModel struct {
	Current   int `json:"current"`
	Withdrawn int `json:"withdrawn"`
}
