package withdraw

import (
	"context"
	"gofermart/internal/model/auth"
	modelWithdraw "gofermart/internal/model/withdraw"
)

type usecaseWithdraw struct {
	repo repositoryWithdraw
}

type repositoryWithdraw interface {
	WithdrawFromUserBalance(ctx context.Context, addModel modelWithdraw.AddUserWithdrawModel) error
}

func New(repo repositoryWithdraw) *usecaseWithdraw {
	return &usecaseWithdraw{
		repo: repo,
	}
}

func (uc *usecaseWithdraw) WithdrawFromUserBalance(ctx context.Context, userID auth.UserID, orderID int64, sum int) error {
	return uc.repo.WithdrawFromUserBalance(ctx, modelWithdraw.AddUserWithdrawModel{
		UserID:  userID,
		OrderID: orderID,
		Sum:     sum,
	})
}
