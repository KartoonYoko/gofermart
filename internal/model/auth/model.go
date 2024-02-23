package auth

import "github.com/golang-jwt/jwt/v5"

type UserID int64

// Claims — структура утверждений, которая включает стандартные утверждения
// и одно пользовательское — UserID
type Claims struct {
	jwt.RegisteredClaims
	UserID UserID
}

type GetUserByLoginAndPasswordModel struct {
	ID       UserID `db:"id"`
	Login    string `db:"login"`
	Password string `db:"password"`
}
