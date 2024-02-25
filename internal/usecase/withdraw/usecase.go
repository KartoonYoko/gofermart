package withdraw

import (
	"context"
	"gofermart/internal/common/monetary"
	"gofermart/internal/model/auth"
	modelWithdraw "gofermart/internal/model/withdraw"
	"strconv"
)

type usecaseWithdraw struct {
	repo repositoryWithdraw
}

type repositoryWithdraw interface {
	WithdrawFromUserBalance(ctx context.Context, addModel modelWithdraw.AddUserWithdrawModel) error
	GetUserWithdrawals(ctx context.Context, userID auth.UserID) ([]modelWithdraw.GetUserWithdrawModel, error)
}

func New(repo repositoryWithdraw) *usecaseWithdraw {
	return &usecaseWithdraw{
		repo: repo,
	}
}

func (uc *usecaseWithdraw) WithdrawFromUserBalance(ctx context.Context, userID auth.UserID, orderID int64, sum float64) error {
	return uc.repo.WithdrawFromUserBalance(ctx, modelWithdraw.AddUserWithdrawModel{
		UserID:  userID,
		OrderID: orderID,
		Sum:     monetary.GetCurencyFromFloat64(sum),
	})
}

func (uc *usecaseWithdraw) GetUserWithdrawals(ctx context.Context, userID auth.UserID) ([]modelWithdraw.GetUserWithdrawAPIModel, error) {
	res, err := uc.repo.GetUserWithdrawals(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := []modelWithdraw.GetUserWithdrawAPIModel{}
	for _, r := range res {
		result = append(result, modelWithdraw.GetUserWithdrawAPIModel{
			OrderID:     strconv.FormatInt(r.OrderID, 10),
			ProcessedAt: r.ProcessedAt,
			Sum:         monetary.GetFloat64FromCurrency(r.Sum),
		})
	}
	return result, nil
}
