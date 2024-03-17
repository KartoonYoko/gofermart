package balance

import (
	"context"
	"errors"
	"gofermart/internal/common/suithelp"
	"gofermart/internal/model/auth"
	modelBalance "gofermart/internal/model/balance"
	"gofermart/internal/repository/pgsql"
	"testing"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type PostgresTestSuite struct {
	suite.Suite
	repositoryBalance

	tc  *tcpostgres.PostgresContainer
	cfg *pgsql.Config
}

func (r *repositoryBalance) clean(ctx context.Context) error {
	query := `DELETE FROM users`
	_, err := r.conn.ExecContext(ctx, query)
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
	ts.repositoryBalance = *New(ctx, db)

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

func TestBalanceRepository(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}

func (ts *PostgresTestSuite) Test_repositoryBalance_GetUserBalance() {
	ctx := context.Background()
	type args struct {
		userLogin    string
		userPassword string
		userBalance  int
		userWithdraw int
	}
	tests := []struct {
		name    string
		args    args
		want    *modelBalance.GetUserBalanceModel
		wantErr bool
	}{
		{
			name: "#1",
			args: args{
				userLogin:    "User 1",
				userPassword: "password",
				userBalance:  100,
				userWithdraw: 0,
			},
		},
		{
			name: "#2",
			args: args{
				userLogin:    "User 2",
				userPassword: "password",
				userBalance:  0,
				userWithdraw: 0,
			},
		},
	}
	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			setUser := setUserModel{
				Login:    tt.args.userLogin,
				Password: tt.args.userPassword,
				Withdraw: tt.args.userWithdraw,
				Balance:  tt.args.userBalance,
			}
			id, err := ts.setUserWithBalanceIfNotExists(setUser)
			require.NoError(ts.T(), err)
			got, err := ts.repositoryBalance.GetUserBalance(ctx, id)
			require.NoError(ts.T(), err)
			require.Equal(ts.T(), got.Current, tt.args.userBalance)
		})
	}
}

type setUserModel struct {
	Login    string `db:"login"`
	Password string `db:"password"`
	Balance  int    `db:"loyality_balance_current"`
	Withdraw int    `db:"loyality_balance_withdrawn"`
}

func (ts *PostgresTestSuite) setUserWithBalanceIfNotExists(model setUserModel) (auth.UserID, error) {
	_, err := ts.repositoryBalance.conn.NamedExec(`INSERT INTO users (login, password, loyality_balance_current, loyality_balance_withdrawn)
        VALUES (:login, :password, :loyality_balance_current, :loyality_balance_withdrawn)`, model)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation != pgErr.Code {
			return -1, err
		}
	}

	var id auth.UserID
	err = ts.repositoryBalance.conn.Get(&id, `SELECT id FROM users WHERE login=$1 AND password=$2`, model.Login, model.Password)
	if err != nil {
		return -1, err
	}
	return id, nil
}
