package orderaccrual

import (
	"context"
	"crypto/rand"
	"gofermart/internal/logger"
	model "gofermart/internal/model/order_accrual"
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
}

type OrderAccrualAPI interface {
	GetOrderAccrual(ctx context.Context, orderID int64) (*model.GetOrderAccrualAPIModel, error)
}

func New(repo OrderAccrualRepository, api OrderAccrualAPI) *orderAccrualUsecase {
	return &orderAccrualUsecase{
		repo: repo,
		api:  api,
	}
}

type requestAPIResult struct {
	orderID       int64
	orderAPIModel *model.GetOrderAccrualAPIModel
	err           error
}

func (uc *orderAccrualUsecase) Go(ctx context.Context) {
	unhandledOrders, err := uc.repo.GetUnhandledOrders(ctx)
	if err != nil {
		logger.Log.Error("get unhandled orders accrual usecase", zap.Error(err))
		return
	}

	result := uc.get(ctx, unhandledOrders)
	successResult := make([]requestAPIResult, 0, len(result))
	unSuccessRsult := make([]requestAPIResult, 0, len(result))
	for _, r := range result {
		if r.err == nil {
			successResult = append(successResult, r)
		} else {
			unSuccessRsult = append(unSuccessRsult, r)
		}
	}

	// лоигруем ошибки
	for _, r := range unSuccessRsult {
		logger.Log.Error("get order accrual usecase", zap.Int64("orderID", r.orderID), zap.Error(r.err))
	}

	// - получить все необработанные заказы
	// - сформировать данные для запросов к АПИ
	// - запустить горутины для сбора данных из АПИ
	// - сохранить статус обработки и начисленные баллы как в таблицу заказво так и в таблицу пользователей
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
				orderID: order.ID,
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

type Effector func(ctx context.Context, orderID int64) (*model.GetOrderAccrualAPIModel, error)

func Retry(effector Effector, retries int) Effector {
	return func(ctx context.Context, orderID int64) (*model.GetOrderAccrualAPIModel, error) {
		for r := 0; ; r++ {
			response, err := effector(ctx, orderID)
			if err == nil || r >= retries {
				return response, err
			}

			// стандартное время ожидания
			delay := time.Second * 5
			// генерируем случайное время ожидания
			b := make([]byte, 1)
			_, err = rand.Read(b)
			if err == nil {
				delay = time.Duration(b[0]) * time.Second
			}
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}
}
