package app

import (
	"context"
	"fmt"
	"gofermart/config"
	"gofermart/internal/controller/httpserver"
	"gofermart/internal/logger"
	"gofermart/internal/repository/pgsql"
	repoAuth "gofermart/internal/repository/pgsql/auth"
	repoOrder "gofermart/internal/repository/pgsql/order"
	usecaseAuth "gofermart/internal/usecase/auth"
	usecaseOrder "gofermart/internal/usecase/order"
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
	db := initDB(ctx, *conf)
	repositoryAuth, err := repoAuth.New(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
	repoOrder, err := repoOrder.New(ctx, db)
	if err != nil {
		log.Fatal(err)
	}

	// хешер для паролей
	hasher := hash.NewSHA1PasswordHasher(confAuth.Sault)

	// usecases
	usecaseAuth := usecaseAuth.New(confJWT, confAuth, repositoryAuth, hasher)
	usecaseOrder := usecaseOrder.New(repoOrder)

	//
	controller := httpserver.New(conf, usecaseAuth, usecaseOrder)
	controller.Serve(ctx)
}

func initDB(ctx context.Context, conf config.Config) *sqlx.DB {
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
	return db
}