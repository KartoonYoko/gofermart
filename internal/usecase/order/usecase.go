package order

import (
	"context"
	"errors"
	"gofermart/internal/common/monetary"
	"gofermart/internal/logger"
	"gofermart/internal/model/auth"
	model "gofermart/internal/model/order"
	"gofermart/internal/repository/pgsql/order"
	"strconv"

	"go.uber.org/zap"
)

type orderUsecase struct {
	repo OrderRepository
}

type OrderRepository interface {
	AddOrder(ctx context.Context, addModel *model.AddOrderModel) error
	GetUserOrders(ctx context.Context, userID auth.UserID) ([]model.GetUserOrderModel, error)
}

func New(repo OrderRepository) *orderUsecase {
	return &orderUsecase{
		repo: repo,
	}
}

// CreateNewOrder - создаёт новый заказ
//
// Возможные ошибки:
//
//	ErrOrderAlreadyExists - у данного пользователя уже существует заказ с таким ID
//	ErrOrderBelongsToAnotherUser - данный заказ зарегистрировал другой пользователь
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

func (uc *orderUsecase) GetUserOrders(ctx context.Context, userID auth.UserID) ([]model.GetUserOrderAPIModel, error) {
	res, err := uc.repo.GetUserOrders(ctx, userID)
	if err != nil {
		logger.Log.Error("order usecase", zap.Error(err))
		return nil, err
	}

	response := []model.GetUserOrderAPIModel{}
	for _, r := range res {
		response = append(response, model.GetUserOrderAPIModel{
			OrderID:    strconv.FormatInt(r.OrderID, 10),
			Status:     r.Status,
			Accrual:    monetary.GetFloat64FromCurrency(r.Accrual),
			UploadedAt: r.CreatedAt,
		})
	}

	return response, nil
}
