package withdraw

import (
	"context"
	"database/sql"
	"errors"
	"gofermart/internal/common/suithelp"
	"gofermart/internal/model/auth"
	modelWithdraw "gofermart/internal/model/withdraw"
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
	repositoryWithdraw

	tc  *tcpostgres.PostgresContainer
	cfg *pgsql.Config
}

func (r *repositoryWithdraw) clean(ctx context.Context) error {
	query := `DELETE FROM withdrawals`
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
	ts.repositoryWithdraw = *New(ctx, db)

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

func (ts *PostgresTestSuite) Test_repositoryWithdraw_WithdrawFromUserBalance() {
	ctx := context.Background()

	type args struct {
		createUserModel createUser
		addModel        modelWithdraw.AddUserWithdrawModel
	}
	tests := []struct {
		name                           string
		args                           args
		wantErrUserHasNotEnoughBalance bool
	}{
		{
			name: "Success withdraw",
			args: args{
				createUserModel: createUser{
					Login:     "User 1",
					Password:  "1",
					Balance:   100,
					Withdrawn: 0,
				},
				addModel: modelWithdraw.AddUserWithdrawModel{
					UserID:  -1,
					OrderID: 1,
					Sum:     100,
				},
			},
			wantErrUserHasNotEnoughBalance: false,
		},
		{
			name: "Withdraw with error ErrUserHasNotEnoughBalance",
			args: args{
				createUserModel: createUser{
					Login:     "User 2",
					Password:  "1",
					Balance:   100,
					Withdrawn: 0,
				},
				addModel: modelWithdraw.AddUserWithdrawModel{
					UserID:  -1,
					OrderID: 2,
					Sum:     101,
				},
			},
			wantErrUserHasNotEnoughBalance: true,
		},
	}
	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			userID, err := ts.createUserIfNotExist(ctx, tt.args.createUserModel)
			require.NoError(ts.T(), err)
			tt.args.addModel.UserID = userID

			err = ts.repositoryWithdraw.WithdrawFromUserBalance(ctx, tt.args.addModel)
			if tt.wantErrUserHasNotEnoughBalance {
				require.ErrorIs(ts.T(), err, modelWithdraw.ErrUserHasNotEnoughBalance)
				return
			}
			require.NoError(ts.T(), err)

			query := `
				SELECT loyality_balance_current FROM users WHERE id=$1
			`
			type GetBalanceModel struct {
				Balance int `db:"loyality_balance_current"`
			}
			getModel := GetBalanceModel{}
			err= ts.conn.GetContext(ctx, &getModel, query, tt.args.addModel.UserID)
			require.NoError(ts.T(), err)

			expectedBalance := tt.args.createUserModel.Balance - tt.args.addModel.Sum
			require.Equal(ts.T(), expectedBalance, getModel.Balance)
		})
	}
}

type createUser struct {
	Login     string `db:"login"`
	Password  string `db:"password"`
	Balance   int    `db:"loyality_balance_current"`
	Withdrawn int    `db:"loyality_balance_withdrawn"`
}

func (ts *PostgresTestSuite) createUserIfNotExist(ctx context.Context, addModel createUser) (auth.UserID, error) {
	query := `SELECT id FROM users WHERE login=$1`
	type GetUserIDModel struct {
		ID auth.UserID `db:"id"`
	}
	getID := GetUserIDModel{}
	err := ts.conn.GetContext(ctx, &getID, query, addModel.Login)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return getID.ID, err
		}
	}

	query = `
		INSERT INTO users (login, password, loyality_balance_current, loyality_balance_withdrawn)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var userID auth.UserID
	err = ts.conn.
		QueryRowContext(ctx, query, addModel.Login, addModel.Password, addModel.Balance, addModel.Withdrawn).
		Scan(&userID)
	if err != nil {
		return -1, err
	}

	return userID, nil
}

func Test_repositoryWithdraw_GetUserWithdrawals(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID auth.UserID
	}
	tests := []struct {
		name    string
		r       *repositoryWithdraw
		args    args
		want    []modelWithdraw.GetUserWithdrawModel
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.GetUserWithdrawals(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("repositoryWithdraw.GetUserWithdrawals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("repositoryWithdraw.GetUserWithdrawals() = %v, want %v", got, tt.want)
			}
		})
	}
}
