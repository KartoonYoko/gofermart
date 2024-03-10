package auth

import (
	"context"
	"gofermart/internal/common/suithelp"
	model "gofermart/internal/model/auth"
	"gofermart/internal/repository/pgsql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

type PostgresTestSuite struct {
	suite.Suite
	authRepository

	tc  *tcpostgres.PostgresContainer
	cfg *pgsql.Config
}

func (r *authRepository) clean(ctx context.Context) error {
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
	ts.authRepository = *New(ctx, db)

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

func TestPostgres(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}

func (ts *PostgresTestSuite) Test_authRepository_AddUser() {
	ctx := context.Background()
	type args struct {
		login    string
		password string
	}
	tests := []struct {
		name                      string
		args                      args
		wantErrLoginAlreadyExists bool
	}{
		{
			name: "User #1",
			args: args{
				login:    "Ivan",
				password: "Password",
			},
		},
		{
			name: "User with same password",
			args: args{
				login:    "Vlad",
				password: "Password",
			},
		},
		{
			name: "User with same login",
			args: args{
				login:    "Ivan",
				password: "123",
			},
			wantErrLoginAlreadyExists: true,
		},
	}
	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			userID, err := ts.AddUser(ctx, tt.args.login, tt.args.password)
			if tt.wantErrLoginAlreadyExists {
				var e *ErrLoginAlreadyExists
				require.ErrorAs(ts.T(), err, &e)
				return
			}
			require.NoError(ts.T(), err)

			// проверим существование пользователя
			type userGETModel struct {
				Password string `db:"password"`
				Login    string `db:"login"`
			}
			model := &userGETModel{}
			query := `SELECT login, password FROM users WHERE id=$1`
			err = ts.authRepository.conn.GetContext(ctx, model, query, userID)
			require.NoError(ts.T(), err)
			require.Equal(ts.T(), tt.args.login, model.Login)
			require.Equal(ts.T(), tt.args.password, model.Password)
		})
	}
}

func (ts *PostgresTestSuite) Test_authRepository_GetUserByLoginAndPassword() {
	ctx := context.Background()
	type args struct {
		login    string
		password string
	}
	tests := []struct {
		name                string
		args                args
		want                *model.GetUserByLoginAndPasswordModel
		wantErrUserNotFound bool
	}{
		{
			name: "Not found user",
			args: args{
				login:    "Ivan",
				password: "password",
			},
			want:                nil,
			wantErrUserNotFound: true,
		},
		{
			name: "Found user",
			args: args{
				login:    "Ivan 2",
				password: "password",
			},
			want: &model.GetUserByLoginAndPasswordModel{
				ID:       1,
				Login:    "Ivan 2",
				Password: "password",
			},
		},
	}
	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			if tt.want != nil {
				var id model.UserID
				query := `
					INSERT INTO users (login, password, loyality_balance_current, loyality_balance_withdrawn)
					VALUES ($1, $2, 0, 0)
					RETURNING id
				`
				err := ts.authRepository.conn.QueryRowContext(ctx, query, tt.want.Login, tt.want.Password).Scan(&id)
				require.NoError(ts.T(), err)
				tt.want.ID = id
			}

			got, err := ts.authRepository.GetUserByLoginAndPassword(ctx, tt.args.login, tt.args.password)
			if tt.wantErrUserNotFound {
				var e *ErrUserNotFound
				require.ErrorAs(ts.T(), err, &e)
				return
			}
			require.NoError(ts.T(), err)

			require.Equal(ts.T(), tt.want.ID, got.ID)
			require.Equal(ts.T(), tt.want.Login, got.Login)
			require.Equal(ts.T(), tt.want.Password, got.Password)
		})
	}
}
