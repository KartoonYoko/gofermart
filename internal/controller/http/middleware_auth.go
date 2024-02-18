package http

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

// тип ключа контекста для middleware аутентификации
type MiddlewareAuthKey int

const (
	keyUserID MiddlewareAuthKey = iota // ключ для ID пользователя
)

// const TOKEN_EXP = time.Hour * 3
const SecretKey = "supersecretkey"

// authJWTCookieMiddleware авторизует пользователя из куки
func (c *HttpController) authJWTCookieMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO получать куки авторизации;
		// если куки нет или токен из куки не проходит валидацию - возвращать 401

		// var userID string
		// ctx := r.Context()
		// cookie, err := r.Cookie("Authorization")
		// if err != nil {
		// }

		// ctx = context.WithValue(r.Context(), keyUserID, userID)
		// r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// buildJWTString создаёт токен и возвращает его в виде строки.
func buildJWTString(userID string) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID,
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}

// validateAndGetUserID валидирует токен и получает из него UserID
func validateAndGetUserID(tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SecretKey), nil
		})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("token is not valid")
	}

	return claims.UserID, nil
}
