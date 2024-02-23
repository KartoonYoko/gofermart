package app

import (
	"context"
	"fmt"
	"gofermart/config"
	"gofermart/internal/controller/http"
	"gofermart/internal/logger"
	"gofermart/internal/repository/pgsql"
	repoAuth "gofermart/internal/repository/pgsql/auth"
	usecaseAuth "gofermart/internal/usecase/auth"
	"gofermart/pkg/hash"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func Run() {
	ctx := context.TODO()
	// логгер
	if err := logger.Initialize("Info"); err != nil {
		log.Fatal(fmt.Errorf("logger init error: %w", err))
	}
	defer logger.Log.Sync()

	// конфигурация
	conf := config.New()
	confJWT, err := config.NewJWTConfig()
	if err != nil {
		log.Fatal(err)
	}
	confAuth := config.NewAuthConfig("some sault")

	// хранилище
	db, err := sqlx.Connect("pgx", conf.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	init, err := pgsql.NewInitializer(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
	err = init.Init(ctx)
	if err != nil {
		log.Fatal(err)
	}
	repositoryAuth, err := repoAuth.New(ctx, db)
	if err != nil {
		log.Fatal(err)
	}

	// хешер для паролей
	hasher := hash.NewSHA1PasswordHasher(confAuth.Sault)

	// usecases
	usecaseAuth, err := usecaseAuth.New(confJWT, confAuth, repositoryAuth, hasher)
	if err != nil {
		log.Fatal(err)
	}

	//
	controller := http.New(conf, usecaseAuth)
	controller.Serve()
}
