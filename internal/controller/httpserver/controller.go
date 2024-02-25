package httpserver

import (
	"context"
	"gofermart/config"
	modelAuth "gofermart/internal/model/auth"
	modelBalance "gofermart/internal/model/balance"
	modelOrder "gofermart/internal/model/order"
	modelWithdraw "gofermart/internal/model/withdraw"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type HTTPController struct {
	r                   *chi.Mux
	conf                *config.Config
	usecaseAuth         usecaseAuth
	usecaseOrder        usecaseOrder
	usecaseOrderAccrual usecaseOrderAccrual
	usecaseBalance      usecaseBalance
	usecaseWithdraw     usecaseWithdraw
}

type usecaseWithdraw interface {
	WithdrawFromUserBalance(ctx context.Context, userID modelAuth.UserID, orderID int64, sum float64) error
	GetUserWithdrawals(ctx context.Context, userID modelAuth.UserID) ([]modelWithdraw.GetUserWithdrawAPIModel, error)
}

type usecaseBalance interface {
	GetUserBalance(ctx context.Context, userID modelAuth.UserID) (*modelBalance.GetUserBalanceAPIModel, error)
}

type usecaseOrder interface {
	CreateNewOrder(ctx context.Context, userID modelAuth.UserID, orderID int64) error
	GetUserOrders(ctx context.Context, userID modelAuth.UserID) ([]modelOrder.GetUserOrderAPIModel, error)
}

type usecaseOrderAccrual interface {
	StartWorkerToHandleOrderAccrual(ctx context.Context)
}

type usecaseAuth interface {
	RegisterAndGetUserJWT(ctx context.Context, login string, password string) (string, error)
	LoginAndGetUserJWT(ctx context.Context, login string, password string) (string, error)
	ValidateJWTAndGetUserID(tokenString string) (modelAuth.UserID, error)
}

func New(conf *config.Config,
	ucAuth usecaseAuth,
	ucOrder usecaseOrder,
	usecaseOrderAccrual usecaseOrderAccrual,
	usecaseBalance usecaseBalance,
	usecaseWithdraw usecaseWithdraw,
) *HTTPController {
	c := &HTTPController{
		conf:                conf,
		usecaseAuth:         ucAuth,
		usecaseOrder:        ucOrder,
		usecaseOrderAccrual: usecaseOrderAccrual,
		usecaseBalance:      usecaseBalance,
	}
	r := chi.NewRouter()
	c.r = r

	// TODO вынести создание роутов в отдельный метод
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	c.routeAPI()

	return c
}

func (c *HTTPController) routeAPI() {
	userRoute := c.routeUser()
	c.r.Mount("/api/user", userRoute)
}

func (c *HTTPController) routeUser() *chi.Mux {
	userRouter := chi.NewRouter()

	// общедоступные конечные точки
	userRouter.Group(func(r chi.Router) {
		r.Post("/register", c.handlerUserRegisterPOST)
		r.Post("/login", c.handlerUserLoginPOST)
	})

	// конечные точки для аутентифицированных пользователей
	userRouter.Group(func(r chi.Router) {
		r.Use(c.middlewareAuth)

		r.Post("/orders", c.handlerUserOrdersPOST)
		r.Get("/orders", c.handlerUserOrdersGET)

		r.Get("/balance", c.handlerUserBalanceGET)
		r.Post("/balance/withdraw", c.handlerUserBalanceWithdrawPOST)

		r.Get("/withdrawals", c.handlerUserWithdrawalsGET)
	})

	return userRouter
}

func (c *HTTPController) StartWorkers(ctx context.Context) {
	c.usecaseOrderAccrual.StartWorkerToHandleOrderAccrual(ctx)
}

func (c *HTTPController) Serve(ctx context.Context) {
	server := &http.Server{Addr: c.conf.RunAddress, Handler: c.r}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(ctx)

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

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
