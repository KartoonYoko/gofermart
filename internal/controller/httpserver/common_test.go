package httpserver

import (
	"context"
	"errors"
	"gofermart/config"
	"testing"
	"time"

	repoAuth "gofermart/internal/repository/pgsql/auth"

	mocksAuth "gofermart/internal/mocks/auth"
	mocksBalance "gofermart/internal/mocks/balance"
	mocksOrder "gofermart/internal/mocks/order"
	mocksOrderAccrual "gofermart/internal/mocks/order_accrual"
	mocksWithdraw "gofermart/internal/mocks/withdraw"
	"gofermart/internal/model/auth"
	usecaseAuthPackage "gofermart/internal/usecase/auth"

	usecaseBalancePackage "gofermart/internal/usecase/balance"
	usecaseOrderPackage "gofermart/internal/usecase/order"
	usecaseOrderAccrualPackage "gofermart/internal/usecase/order_accrual"
	usecaseWithdrawPackage "gofermart/internal/usecase/withdraw"
	"gofermart/pkg/hash"

	"go.uber.org/mock/gomock"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type HTTPControllerTestSuite struct {
	suite.Suite
	HTTPController

	tc *tcpostgres.PostgresContainer
}

func (ts *HTTPControllerTestSuite) SetupSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pgc, err := tcpostgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:16.2"),
		tcpostgres.WithDatabase("gophermart"),
		tcpostgres.WithUsername("postgres"),
		tcpostgres.WithPassword("123"),
		tcpostgres.WithInitScripts(),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)

	require.NoError(ts.T(), err)

	// cfg.Host, err = pgc.Host(ctx)
	// require.NoError(ts.T(), err)

	// port, err := pgc.MappedPort(ctx, "5432")
	// require.NoError(ts.T(), err)

	// cfg.Port = uint16(port.Int())

	ts.tc = pgc
	// ts.cfg = cfg

	// db, err := New(cfg)
	// require.NoError(ts.T(), err)

	// ts.testStorager = db

	// ts.T().Logf("stared postgres at %s:%d", cfg.Host, cfg.Port)
}

func (ts *HTTPControllerTestSuite) TearDownSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	require.NoError(ts.T(), ts.tc.Terminate(ctx))
}

// func (ts *PostgresTestSuite) SetupTest() {
// 	ts.Require().NoError(ts.clean(context.Background()))
// }

// func (ts *PostgresTestSuite) TearDownTest() {
// 	ts.Require().NoError(ts.clean(context.Background()))
// }

func TestPostgres(t *testing.T) {
	suite.Run(t, new(HTTPControllerTestSuite))
}

func createTestController(ctrl *gomock.Controller, ctx context.Context) *HTTPController {
	// ctrl := gomock.NewController(t)
	// defer ctrl.Finish()
	repositoryAuth := mocksAuth.NewMockAuthRepository(ctrl)
	repoBalance := mocksBalance.NewMockRepositoryBalance(ctrl)
	repoOrder := mocksOrder.NewMockOrderRepository(ctrl)
	repoOrderAccrual := mocksOrderAccrual.NewMockOrderAccrualRepository(ctrl)
	apiOrderAccrual := mocksOrderAccrual.NewMockOrderAccrualAPI(ctrl)
	repoWithdraw := mocksWithdraw.NewMockRepositoryWithdraw(ctrl)

	repositoryAuth.EXPECT().AddUser(gomock.Any(), "testuser", "8bb6118f8fd6935ad0876a3be34a717d32708ffd").Return(auth.UserID(1), nil)
	repositoryAuth.EXPECT().AddUser(gomock.Any(), "testuser", "8bb6118f8fd6935ad0876a3be34a717d32708ffd").Return(auth.UserID(-1), repoAuth.NewErrLoginAlreadyExists("testuser", errors.New("asdad")))
	repositoryAuth.EXPECT().AddUser(gomock.Any(), "", "").AnyTimes().Return(auth.UserID(-1), errors.New("wrong format"))

	confJWT := &config.JWTConfig{}
	confAuth := &config.AuthConfig{}
	conf := &config.Config{}

	// хешер для паролей
	hasher := hash.NewSHA1PasswordHasher(confAuth.Sault)

	// usecases
	usecaseAuth := usecaseAuthPackage.New(confJWT, confAuth, repositoryAuth, hasher)
	usecaseOrder := usecaseOrderPackage.New(repoOrder)
	usecaseOrderAccrual := usecaseOrderAccrualPackage.New(repoOrderAccrual, apiOrderAccrual)
	usecaseBalance := usecaseBalancePackage.New(repoBalance)
	usecaseWithdraw := usecaseWithdrawPackage.New(repoWithdraw)

	controller := New(conf,
		usecaseAuth,
		usecaseOrder,
		usecaseOrderAccrual,
		usecaseBalance,
		usecaseWithdraw,
	)

	return controller
}
