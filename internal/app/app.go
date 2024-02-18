package app

import (
	"gofermart/config"
	"gofermart/internal/controller/http"
)

func Run() {
	// ctx := context.TODO()
	conf := config.New()
	controller := http.New(conf)
	controller.Serve()
}
