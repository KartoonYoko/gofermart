package orderaccrual

import (
	"context"
	"gofermart/config"
	model "gofermart/internal/model/order_accrual"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
)

type orderAccrualAPI struct {
	client *resty.Client
}

func New(conf config.Config) *orderAccrualAPI {
	client := resty.New().SetBaseURL(conf.AccrualSystemAddress)

	return &orderAccrualAPI{
		client: client,
	}
}

func (api *orderAccrualAPI) GetOrderAccrual(ctx context.Context, orderID int64) (*model.GetOrderAccrualAPIModel, error) {
	var res *model.GetOrderAccrualAPIModel = &model.GetOrderAccrualAPIModel{}
	resp, err := api.client.R().
		SetContext(ctx).
		SetResult(res).
		Get("/api/orders/" + strconv.FormatInt(orderID, 10))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() == http.StatusOK {
		return res, nil
	}

	if resp.StatusCode() == http.StatusNoContent {
		return nil, model.ErrOrderNotFound
	}
	if resp.StatusCode() == http.StatusTooManyRequests {
		return nil, model.ErrTooManyRequests
	}

	return nil, model.ErrUndefinedStatusCode
}
