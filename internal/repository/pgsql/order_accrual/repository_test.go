package orderaccrual

import (
	"context"
	"gofermart/internal/common/suithelp"
	"gofermart/internal/model/auth"
	modelOrder "gofermart/internal/model/order"
	"gofermart/internal/repository/pgsql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

type PostgresTestSuite struct {
	suite.Suite
	orderAccrualRepository

	tc  *tcpostgres.PostgresContainer
	cfg *pgsql.Config
}

func (r *orderAccrualRepository) clean(ctx context.Context) error {
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
	ts.orderAccrualRepository = *New(ctx, db)

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

func TestOrderAccrualRepository(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}

func (ts *PostgresTestSuite) Test_orderAccrualRepository_GetUnhandledOrders() {
	ctx := context.Background()

	query := `
		INSERT INTO users (login, password, loyality_balance_current, loyality_balance_withdrawn)
		VALUES ($1, $2, 0, 0)
		RETURNING id
	`
	var userID auth.UserID
	err := ts.conn.QueryRowContext(ctx, query, "User 1", "password").Scan(&userID)
	require.NoError(ts.T(), err)

	type AddOrderModel struct {
		OrderID int                    `db:"order_id"`
		Status  modelOrder.OrderStatus `db:"status"`
		Accrual int                    `db:"accrual"`
		UserID  auth.UserID            `db:"user_id"`
	}
	ordersStructs := []AddOrderModel{
		{OrderID: 1, Status: modelOrder.OrderStatusNew, Accrual: 0, UserID: userID},
		{OrderID: 2, Status: modelOrder.OrderStatusProcessing, Accrual: 0, UserID: userID},
		{OrderID: 3, Status: modelOrder.OrderStatusProcessing, Accrual: 0, UserID: userID},
	}
	query = `
	INSERT INTO orders (id, status, accrual, user_id) 
	VALUES (:order_id, :status, :accrual, :user_id)
	`

	_, err = ts.orderAccrualRepository.conn.NamedExec(query, ordersStructs)
	require.NoError(ts.T(), err)
	got, err := ts.orderAccrualRepository.GetUnhandledOrders(ctx)
	require.NoError(ts.T(), err)
	for _, v := range got {
		var foundOrder *AddOrderModel
		for _, order := range ordersStructs {
			if order.OrderID == int(v.ID) {
				foundOrder = &order
				break
			}
		}

		require.NotNil(ts.T(), foundOrder)
		assert.Equal(ts.T(), v.Status, string(foundOrder.Status))
		assert.Equal(ts.T(), v.UserID, foundOrder.UserID)
	}

	require.Equal(ts.T(), len(ordersStructs), len(got))
}

func (ts *PostgresTestSuite) Test_orderAccrualRepository_AccrualToOrderAndUser() {
	// Создать пользователя
	// Создать заказ
	// Вызвать метод начисления
	// Проверить начисление в заказе и у пользователя
	ctx := context.Background()
	accrual := 1200
	var orderID int64 = 123123
	orderStatus := modelOrder.OrderStatusProcessed

	query := `
	INSERT INTO users (login, password, loyality_balance_current, loyality_balance_withdrawn)
	VALUES ($1, $2, 0, 0)
	RETURNING id
	`
	var userID int
	err := ts.conn.QueryRowContext(ctx, query, "User 1", "password").Scan(&userID)
	require.NoError(ts.T(), err)

	query = `
	INSERT INTO orders (id, status, accrual, user_id) 
	VALUES (:order_id, :status, :accrual, :user_id)
	`
	_, err = ts.conn.NamedExecContext(ctx, query, map[string]any{
		"order_id": orderID,
		"status":   orderStatus,
		"accrual":  accrual,
		"user_id":  userID,
	})
	require.NoError(ts.T(), err)

	err = ts.orderAccrualRepository.AccrualToOrderAndUser(ctx, orderID, auth.UserID(userID), accrual, orderStatus)
	require.NoError(ts.T(), err)

	query = `
	SELECT loyality_balance_current FROM users WHERE id=$1
	`
	var currentUserBalance int
	err = ts.conn.QueryRowContext(ctx, query, userID).Scan(&currentUserBalance)
	require.NoError(ts.T(), err)
	assert.Equal(ts.T(), accrual, currentUserBalance)

	query = `
		SELECT accrual FROM orders WHERE id=$1
	`
	var orderAccrual int
	err = ts.conn.QueryRowContext(ctx, query, orderID).Scan(&orderAccrual)
	require.NoError(ts.T(), err)
	assert.Equal(ts.T(), accrual, orderAccrual)
}
