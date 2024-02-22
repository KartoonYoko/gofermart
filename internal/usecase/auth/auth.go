package auth

import (
	"gofermart/config"

	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type authUsecase struct {
	conf *config.JWTConfig
	repo authRepository
}

type authRepository interface {
	AddUser(login string, password string) error
}

// Register регистрирует пользователя
//
// При валидации может вернуть следующие ошибки:
//   - ErrLoginIsOccupiedByAnotherUser логин уже занят
//   - ErrWrongDataFormat неверный формат данных
func (uc *authUsecase) RegisterAndGetUserJWT(login string, password string) error {
	if login == "" || password == "" {
		return ErrWrongDataFormat
	}

	err := uc.repo.AddUser(login, password)
	if err != nil {
		// TODO проверить на ошибку существования пользователя
		return err
	}

	return errors.New("not implemented")
}

func (uc *authUsecase) ValidateJWT() error {
	return errors.New("not implemented")
}

func (uc *authUsecase) Login() error {
	return errors.New("not implemented")
}

// Claims — структура утверждений, которая включает стандартные утверждения
// и одно пользовательское — UserID
type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

// buildJWTString создаёт токен и возвращает его в виде строки.
func (uc *authUsecase) buildJWTString(userID string) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID,
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(uc.conf.SecretJWTKey))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}
