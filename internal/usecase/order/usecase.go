package order

import (
	"context"
	"gofermart/internal/model/auth"
	model "gofermart/internal/model/order"
)

type orderUsecase struct {
	repo OrderRepository
}

type OrderRepository interface {
	AddOrder(ctx context.Context, addModel *model.AddOrderModel) error
}

func New(repo OrderRepository) *orderUsecase {
	return &orderUsecase{
		repo: repo,
	}
}

func (uc *orderUsecase) CreateNewOrder(ctx context.Context, userID auth.UserID, orderID int) error {
	// TODO написать функцию валидации ID заказа
	err := uc.repo.AddOrder(ctx, &model.AddOrderModel{
		UserID:  userID,
		OrderID: orderID,
		Status:  model.OrderStatusNew,
		Accrual: 0,
	})

	if err != nil {
		return err
	}

	return nil
}

func (uc *orderUsecase) validateOrderID(orderID int) error {
	// TODO написать функцию валидации
	// https://ru.wikipedia.org/wiki/%D0%90%D0%BB%D0%B3%D0%BE%D1%80%D0%B8%D1%82%D0%BC_%D0%9B%D1%83%D0%BD%D0%B0
	// model.NewErrWrongFormatOfOrderID()
}
