package app

import (
	"context"
	"fmt"
	"gofermart/config"
	"gofermart/internal/controller/httpserver"
	"gofermart/internal/logger"
	apiOrderAccrual "gofermart/internal/repository/api/order_accrual"
	"gofermart/internal/repository/pgsql"
	repoAuth "gofermart/internal/repository/pgsql/auth"
	repoBalance "gofermart/internal/repository/pgsql/balance"
	repoOrder "gofermart/internal/repository/pgsql/order"
	repoOrderAccrual "gofermart/internal/repository/pgsql/order_accrual"
	repoWithdraw "gofermart/internal/repository/pgsql/withdraw"
	usecaseAuth "gofermart/internal/usecase/auth"
	usecaseBalance "gofermart/internal/usecase/balance"
	usecaseOrder "gofermart/internal/usecase/order"
	usecaseOrderAccrual "gofermart/internal/usecase/order_accrual"
	usecaseWithdraw "gofermart/internal/usecase/withdraw"
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
	repositoryAuth := repoAuth.New(ctx, db)
	repoOrder := repoOrder.New(ctx, db)
	repoOrderAccrual := repoOrderAccrual.New(ctx, db)
	apiOrderAccrual := apiOrderAccrual.New(*conf)
	repoBalance := repoBalance.New(ctx, db)
	repoWithdraw := repoWithdraw.New(ctx, db)

	// хешер для паролей
	hasher := hash.NewSHA1PasswordHasher(confAuth.Sault)

	// usecases
	usecaseAuth := usecaseAuth.New(confJWT, confAuth, repositoryAuth, hasher)
	usecaseOrder := usecaseOrder.New(repoOrder)
	usecaseOrderAccrual := usecaseOrderAccrual.New(repoOrderAccrual, apiOrderAccrual)
	usecaseBalance := usecaseBalance.New(repoBalance)
	usecaseWithdraw := usecaseWithdraw.New(repoWithdraw)

	//
	controller := httpserver.New(conf,
		usecaseAuth,
		usecaseOrder,
		usecaseOrderAccrual,
		usecaseBalance,
		usecaseWithdraw,
	)
	controller.StartWorkers(ctx)
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
