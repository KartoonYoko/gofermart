package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gofermart/internal/common/suithelp"
	"gofermart/internal/model/auth"
	model "gofermart/internal/model/order"
	"gofermart/internal/repository/pgsql"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

type PostgresTestSuite struct {
	suite.Suite
	orderRepository

	tc  *tcpostgres.PostgresContainer
	cfg *pgsql.Config
}

func (r *orderRepository) clean(ctx context.Context) error {
	query := `DELETE FROM orders`
	_, err := r.conn.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	query = `DELETE FROM users`
	_, err = r.conn.ExecContext(ctx, query)

	return err
}

func (ts *PostgresTestSuite) SetupSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pgc, err := suithelp.NewPostgresContainer(ctx)
	require.NoError(ts.T(), err)
	containerData, err := suithelp.GetPostgresqlContainerData(ctx, pgc)
	require.NoError(ts.T(), err)
	cfg := &pgsql.Config{
		ConnectionString: containerData.ConnectionString,
	}
	ts.tc = pgc
	ts.cfg = cfg

	db, err := pgsql.NewSQLxConnection(ctx, cfg)
	require.NoError(ts.T(), err)
	ts.orderRepository = *New(ctx, db)

	ts.T().Logf("stared postgres at %s:%d", containerData.Host, containerData.Port)
}

func (ts *PostgresTestSuite) TearDownSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	require.NoError(ts.T(), ts.tc.Terminate(ctx))
}

func (ts *PostgresTestSuite) SetupTest() {
	ts.Require().NoError(ts.clean(context.Background()))
}

func (ts *PostgresTestSuite) TearDownTest() {
	ts.Require().NoError(ts.clean(context.Background()))
}

func TestOrderRepository(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}

func (ts *PostgresTestSuite) Test_orderRepository_AddOrder() {
	ctx := context.Background()
	type args struct {
		addModel *model.AddOrderModel
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Success 1",
			args: args{
				addModel: &model.AddOrderModel{
					UserID:  1,
					OrderID: 1,
					Status:  "",
					Accrual: 0,
				},
			},
			err: nil,
		},
		{
			name: "Order already exists error",
			args: args{
				addModel: &model.AddOrderModel{
					UserID:  1,
					OrderID: 1,
					Status:  "",
					Accrual: 0,
				},
			},
			err: ErrOrderAlreadyExists,
		},
		{
			name: "Order occupied by another user error",
			args: args{
				addModel: &model.AddOrderModel{
					UserID:  2,
					OrderID: 1,
					Status:  "",
					Accrual: 0,
				},
			},
			err: ErrOrderBelongsToAnotherUser,
		},
	}
	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			err := ts.createUserIfNotExist(ctx, tt.args.addModel.UserID)
			require.NoError(ts.T(), err)
			err = ts.orderRepository.AddOrder(ctx, tt.args.addModel)
			require.ErrorIs(ts.T(), err, tt.err)
		})
	}
}

func Test_orderRepository_GetUserOrders(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID auth.UserID
	}
	tests := []struct {
		name    string
		r       *orderRepository
		args    args
		want    []model.GetUserOrderModel
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.GetUserOrders(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("orderRepository.GetUserOrders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("orderRepository.GetUserOrders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (ts *PostgresTestSuite) createUserIfNotExist(ctx context.Context, userID auth.UserID) error {
	query := `
	SELECT id FROM users WHERE id=$1
	`
	var check auth.UserID
	err := ts.orderRepository.conn.QueryRowContext(ctx, query, userID).Scan(&check)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if check == userID {
		return nil
	}
	ts.orderRepository.conn.ExecContext(ctx, `INSERT INTO users (login, password, loyality_balance_current, loyality_balance_withdrawn)
	VALUES ($1, $2, 0, 0)`, fmt.Sprintf("User_%d", userID), "password")

	return nil
}
