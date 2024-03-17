package orderaccrual

import (
	"context"
	model "gofermart/internal/model/order_accrual"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

type AccrualTestSuite struct {
	suite.Suite
	orderAccrualAPI

	tc  testcontainers.Container
	ctx context.Context
}

func (ts *AccrualTestSuite) SetupSuite() {
	// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel()
	// ts.ctx = context.Background()

	// cd, err := suithelp.NewAccrualContainer(ctx)
	// require.NoError(ts.T(), err)

	// ts.orderAccrualAPI = *New(config.Config{
	// 	AccrualSystemAddress: fmt.Sprintf("http://%s:%d", cd.Host, cd.Port),
	// })
	// ts.tc = cd.Tc

	// ts.T().Logf("stared accrual at %s:%d", cd.Host, cd.Port)
}

func (ts *AccrualTestSuite) TearDownSuite() {
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	// require.NoError(ts.T(), ts.tc.Terminate(ctx))
}

func (ts *AccrualTestSuite) SetupTest() {
	// ts.Require().NoError(ts.clean(context.Background()))
}

func (ts *AccrualTestSuite) TearDownTest() {
	// ts.Require().NoError(ts.clean(context.Background()))
}

func TestOrderAccrualRepository(t *testing.T) {
	suite.Run(t, new(AccrualTestSuite))
}

func (ts *AccrualTestSuite) Test_orderAccrualAPI_GetOrderAccrual() {
	return
	type args struct {
		orderID int64
	}
	tests := []struct {
		name    string
		args    args
		want    *model.GetOrderAccrualFromRemoteModel
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			got, err := ts.orderAccrualAPI.GetOrderAccrual(ts.ctx, tt.args.orderID)
			if (err != nil) != tt.wantErr {
				t.Errorf("orderAccrualAPI.GetOrderAccrual() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("orderAccrualAPI.GetOrderAccrual() = %v, want %v", got, tt.want)
			}
		})
	}
}

type createRewardModel struct {
	Match       string `json:"match"`
	Reward      string `json:"reward"`
	Reward_type string `json:"reward_type"`
}

func (ts *AccrualTestSuite) createRewardIfNotExists() error {
	return nil
}
