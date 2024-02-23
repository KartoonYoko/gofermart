package auth

type GetUserByLoginAndPasswordModel struct {
	ID       int64  `db:"id"`
	Login    string `db:"login"`
	Password string `db:"password"`
}
