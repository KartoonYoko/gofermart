package app

import (
	"context"
	"gofermart/config"
	"gofermart/internal/controller/http"
	"gofermart/internal/repository/pgsql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func Run() {
	ctx := context.TODO()
	conf := config.New()

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

	controller := http.New(conf)
	controller.Serve()
}
