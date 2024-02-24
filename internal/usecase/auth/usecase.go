package auth

import (
	"context"
	"fmt"
	"gofermart/config"
	"gofermart/internal/logger"
	model "gofermart/internal/model/auth"
	repoAuth "gofermart/internal/repository/pgsql/auth"

	"errors"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type authUsecase struct {
	confJWT  *config.JWTConfig
	confAuth *config.AuthConfig

	repo           AuthRepository
	passwordHasher PasswordHasher
}

type AuthRepository interface {
	AddUser(ctx context.Context, login string, password string) (model.UserID, error)
	GetUserByLoginAndPassword(ctx context.Context, login string, password string) (*model.GetUserByLoginAndPasswordModel, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
}

func New(confJWT *config.JWTConfig, confAuth *config.AuthConfig, repo AuthRepository, passwordHasher PasswordHasher) *authUsecase {
	uc := &authUsecase{
		confJWT:        confJWT,
		confAuth:       confAuth,
		repo:           repo,
		passwordHasher: passwordHasher,
	}

	return uc
}

// Register регистрирует пользователя
//
// Может вернуть следующие ошибки:
//   - ErrLoginIsOccupiedByAnotherUser логин уже занят
//   - ErrWrongDataFormat неверный формат данных
func (uc *authUsecase) RegisterAndGetUserJWT(ctx context.Context, login string, password string) (string, error) {
	if login == "" || password == "" {
		logger.Log.Sugar().Infoln("register: wrong data format",
			"login", login,
			"password", password,
		)
		return "", model.ErrWrongDataFormat
	}

	hashPswd, err := uc.passwordHasher.Hash(password)
	if err != nil {
		logger.Log.Error("register:", zap.Error(err))
		return "", err
	}

	userID, err := uc.repo.AddUser(ctx, login, hashPswd)
	if err != nil {
		var errLoginAlreadyExists *repoAuth.ErrLoginAlreadyExists
		if errors.As(err, &errLoginAlreadyExists) {
			logger.Log.Sugar().Infoln("register: user already exists",
				"login", login,
				"password", password,
			)
			return "", model.ErrLoginIsOccupiedByAnotherUser
		}

		logger.Log.Error("register:", zap.Error(err))
		return "", err
	}

	return uc.buildJWTStringWithUserID(userID)
}

// LoginAndGetUserJWT ищёт пользователя по логину и паролю и если находит - возвращает JWT токен, принадлежащий пользователю
//
// Может вернуть следующие ошибки:
//
//	errUserNotFound - пользователь не найден
func (uc *authUsecase) LoginAndGetUserJWT(ctx context.Context, login string, password string) (string, error) {
	hashPswd, err := uc.passwordHasher.Hash(password)
	if err != nil {
		logger.Log.Error("login:", zap.Error(err))
		return "", err
	}

	user, err := uc.repo.GetUserByLoginAndPassword(ctx, login, hashPswd)
	if err != nil {
		var errUserNotFound *repoAuth.ErrUserNotFound
		if errors.As(err, &errUserNotFound) {
			return "", model.ErrUserNotFound
		}

		logger.Log.Error("login:", zap.Error(err))
		return "", err
	}

	return uc.buildJWTStringWithUserID(user.ID)
}

func (uc *authUsecase) ValidateJWTAndGetUserID(tokenString string) (model.UserID, error) {
	claims := &model.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(uc.confJWT.SecretJWTKey), nil
		})

	if err != nil {
		return -1, err
	}

	if !token.Valid {
		return -1, fmt.Errorf("token is not valid")
	}

	return claims.UserID, nil
}

// buildJWTString создаёт токен и возвращает его в виде строки.
func (uc *authUsecase) buildJWTStringWithUserID(userID model.UserID) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.Claims{
		UserID: userID,
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(uc.confJWT.SecretJWTKey))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}
