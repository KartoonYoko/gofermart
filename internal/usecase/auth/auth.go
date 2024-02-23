package auth

import (
	"context"
	"gofermart/config"
	"gofermart/internal/logger"
	model "gofermart/internal/model/auth"
	repoAuth "gofermart/internal/repository/pgsql/auth"
	"strconv"

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
	AddUser(ctx context.Context, login string, password string) (int64, error)
	GetUserByLoginAndPassword(ctx context.Context, login string, password string) (*model.GetUserByLoginAndPasswordModel, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
}

func New(confJWT *config.JWTConfig, confAuth *config.AuthConfig, repo AuthRepository, passwordHasher PasswordHasher) (*authUsecase, error) {
	uc := &authUsecase{
		confJWT:        confJWT,
		confAuth:       confAuth,
		repo:           repo,
		passwordHasher: passwordHasher,
	}

	return uc, nil
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

func (uc *authUsecase) ValidateJWT() error {
	return errors.New("not implemented")
}

func (uc *authUsecase) Login(ctx context.Context, login string, password string) (string, error) {
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

// Claims — структура утверждений, которая включает стандартные утверждения
// и одно пользовательское — UserID
type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

// buildJWTString создаёт токен и возвращает его в виде строки.
func (uc *authUsecase) buildJWTStringWithUserID(userID int64) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: strconv.FormatInt(userID, 10),
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(uc.confJWT.SecretJWTKey))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}
