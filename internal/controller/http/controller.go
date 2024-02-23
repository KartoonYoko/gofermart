package http

import (
	"context"
	"gofermart/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type HttpController struct {
	r    *chi.Mux
	conf *config.Config
	uc   authUsercase
}

type authUsercase interface {
	RegisterAndGetUserJWT(ctx context.Context, login string, password string) (string, error)
	LoginAndGetUserJWT(ctx context.Context, login string, password string) (string, error)
}

func New(conf *config.Config, uc authUsercase) *HttpController {
	c := &HttpController{
		conf: conf,
		uc:   uc,
	}
	r := chi.NewRouter()
	c.r = r

	// TODO вынести создание роутов в отдельный метод
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	c.routeAPI()

	return c
}

func (c *HttpController) routeAPI() {
	userRoute := c.routeUser()
	c.r.Mount("/api/user", userRoute)
}

func (c *HttpController) routeUser() *chi.Mux {
	userRouter := chi.NewRouter()

	// общедоступные конечные точки
	userRouter.Group(func(r chi.Router) {
		userRouter.Post("/register", c.handlerUserRegisterPOST)
		userRouter.Post("/login", c.handlerUserLoginPOST)
	})

	// конечные точки для аутентифицированных пользователей
	userRouter.Group(func(r chi.Router) {
		// TODO добавить middleware aутентификации
		userRouter.Post("/orders", c.handlerUserOrdersPOST)
		userRouter.Get("/orders", c.handlerUserOrdersGET)

		userRouter.Get("/balance", c.handlerUserBalanceGET)
		userRouter.Post("/balance/withdraw", c.handlerUserBalanceWithdrawPOST)

		userRouter.Get("/withdrawals", c.handlerUserWithdrawalsGET)
	})

	return userRouter
}

func (c *HttpController) Serve() {
	server := &http.Server{Addr: c.conf.RunAddress, Handler: c.r}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run the server
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
