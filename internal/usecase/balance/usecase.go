package balance

import (
	"context"
	"gofermart/internal/model/auth"
	modelBalance "gofermart/internal/model/balance"
)

type usecaseBalance struct {
	repo RepositoryBalance
}

type RepositoryBalance interface {
	GetUserBalance(ctx context.Context, userID auth.UserID) (*modelBalance.GetUserBalanceModel, error)
}

func New(repo RepositoryBalance) *usecaseBalance {
	return &usecaseBalance{
		repo: repo,
	}
}

func (uc *usecaseBalance) GetUserBalance(ctx context.Context, userID auth.UserID) (*modelBalance.GetUserBalanceAPIModel, error) {
	res, err := uc.repo.GetUserBalance(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &modelBalance.GetUserBalanceAPIModel{
		Withdrawn: res.Withdrawn,
		Current:   res.Current,
	}, nil
}
