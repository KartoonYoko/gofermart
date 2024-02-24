package orderaccrual

import (
	"context"
	model "gofermart/internal/model/order_accrual"
)

type orderAccrualUsecase struct {
	repo OrderAccrualRepository
}

type OrderAccrualRepository interface {
	GetUnhandledOrders(ctx context.Context) ([]model.GetOrderModel, error)
}

type OrderAccrualAPI interface {
	GetOrderAccrual(orderID int64) (*model.GetOrderAccrualAPIModel, error)
}

func New(repo OrderAccrualRepository) *orderAccrualUsecase {
	return &orderAccrualUsecase{
		repo: repo,
	}
}
