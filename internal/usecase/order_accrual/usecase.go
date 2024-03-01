package orderaccrual

import (
	"context"
	"gofermart/internal/common/monetary"
	"gofermart/internal/logger"
	"gofermart/internal/model/auth"
	modelOrder "gofermart/internal/model/order"
	model "gofermart/internal/model/order_accrual"
	"math/rand"
	"sync"
	"time"

	"go.uber.org/zap"
)

type orderAccrualUsecase struct {
	repo OrderAccrualRepository
	api  OrderAccrualAPI
}

type OrderAccrualRepository interface {
	GetUnhandledOrders(ctx context.Context) ([]model.GetOrderModel, error)
	AccrualToOrderAndUser(ctx context.Context, orderID int64, userID auth.UserID, sum int, orderStatus modelOrder.OrderStatus) error
}

type OrderAccrualAPI interface {
	GetOrderAccrual(ctx context.Context, orderID int64) (*model.GetOrderAccrualFromRemoteModel, error)
}

func New(repo OrderAccrualRepository, api OrderAccrualAPI) *orderAccrualUsecase {
	return &orderAccrualUsecase{
		repo: repo,
		api:  api,
	}
}

type requestAPIResult struct {
	orderModel    *model.GetOrderModel
	orderAPIModel *model.GetOrderAccrualFromRemoteModel
	err           error
}

func (uc *orderAccrualUsecase) StartWorkerToHandleOrderAccrual(ctx context.Context) {
	go func() {
		delay := time.Second * 5
		for {
			uc.updateUnhandledOrdersAccrual(ctx)

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return
			}
		}
	}()
}

// updateUnhandledOrdersAccrual - обновит начисления у необработанных заказов. Также начислит баллы пользователям, которые создали этот заказ
func (uc *orderAccrualUsecase) updateUnhandledOrdersAccrual(ctx context.Context) {
	unhandledOrders, err := uc.repo.GetUnhandledOrders(ctx)
	if err != nil {
		logger.Log.Error("get unhandled orders accrual usecase", zap.Error(err))
		return
	}

	if len(unhandledOrders) == 0 {
		return
	}

	result := uc.get(ctx, unhandledOrders)
	successResult := make([]requestAPIResult, 0, len(result))
	notSuccessResult := make([]requestAPIResult, 0, len(result))
	for _, r := range result {
		if r.err == nil {
			successResult = append(successResult, r)
		} else {
			notSuccessResult = append(notSuccessResult, r)
		}
	}

	// логируем ошибки
	for _, r := range notSuccessResult {
		logger.Log.Error("get order accrual usecase", zap.Int64("orderID", r.orderModel.ID), zap.Error(r.err))
	}

	for _, r := range successResult {
		sum := 0
		if r.orderAPIModel.Accrual != nil && *r.orderAPIModel.Accrual > 0 {
			sum = monetary.GetCurencyFromFloat64(*r.orderAPIModel.Accrual)
		}

		var orderStatus modelOrder.OrderStatus = "NEW"
		if r.orderAPIModel.Status == "REGISTERED" {
			orderStatus = "NEW"
		} else if r.orderAPIModel.Status == "PROCESSING" {
			orderStatus = "PROCESSING"
		} else if r.orderAPIModel.Status == "INVALID" {
			orderStatus = "INVALID"
		} else if r.orderAPIModel.Status == "PROCESSED" {
			orderStatus = "PROCESSED"
		}

		if !orderStatus.Valid() {
			logger.Log.Error("not valid order status", zap.String("status", string(orderStatus)))
			continue
		}
		err = uc.repo.AccrualToOrderAndUser(ctx, r.orderModel.ID, r.orderModel.UserID, sum, orderStatus)
		if err != nil {
			logger.Log.Error("error order accrual", zap.Error(err))
		}
	}
}

func (uc *orderAccrualUsecase) get(ctx context.Context, orders []model.GetOrderModel) []requestAPIResult {
	var m sync.Mutex
	var wg sync.WaitGroup

	results := make([]requestAPIResult, 0, len(orders))
	for i := range orders {
		wg.Add(1)
		order := orders[i]
		go func() {
			requestWithRetry := Retry(uc.api.GetOrderAccrual, 3)
			orderAccrualAPIModel, err := requestWithRetry(ctx, order.ID)
			m.Lock()
			results = append(results, requestAPIResult{
				orderModel:    &order,
				orderAPIModel: orderAccrualAPIModel,
				err:           err,
			})
			m.Unlock()
			wg.Done()
		}()
	}

	wg.Wait()

	return results
}

type Effector func(ctx context.Context, orderID int64) (*model.GetOrderAccrualFromRemoteModel, error)

func Retry(effector Effector, retries int) Effector {
	return func(ctx context.Context, orderID int64) (*model.GetOrderAccrualFromRemoteModel, error) {
		for r := 0; ; r++ {
			response, err := effector(ctx, orderID)
			if err == nil || r >= retries {
				return response, err
			}

			// генерируем случайное время ожидания
			delay := time.Duration(rand.Intn(5)) * time.Second
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}
}
