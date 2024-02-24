package order

import (
	"context"
	"errors"
	"gofermart/internal/model/auth"
	model "gofermart/internal/model/order"
	"gofermart/internal/repository/pgsql/order"
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

// CreateNewOrder - создаёт новый заказ
// 
// Возможные ошибки:
// 	ErrOrderAlreadyExists - у данного пользователя уже существует заказ с таким ID
// 	ErrOrderBelongsToAnotherUser - данный заказ зарегистрировал другой пользователь
func (uc *orderUsecase) CreateNewOrder(ctx context.Context, userID auth.UserID, orderID int64) error {
	err := uc.repo.AddOrder(ctx, &model.AddOrderModel{
		UserID:  userID,
		OrderID: orderID,
		Status:  model.OrderStatusNew,
		Accrual: 0,
	})

	if err != nil {
		if errors.Is(err, order.ErrOrderAlreadyExists) {
			return model.NewErrOrderAlreadyExists(err, orderID)
		}
		if errors.Is(err, order.ErrOrderBelongsToAnotherUser) {
			return model.NewErrOrderBelongsToAnotherUser(err, orderID)
		}
		return err
	}

	return nil
}
