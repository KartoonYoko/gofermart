package httpserver

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func (ts *HTTPControllerTestSuite) TestHTTPController_handlerUserRegisterPOST() {
	ctx := context.TODO()
	ctrl := gomock.NewController(ts.T())
	defer ctrl.Finish()
	controller := createTestController(ctrl, ctx)
	srv := httptest.NewServer(controller.r)
	// останавливаем сервер после завершения теста
	defer srv.Close()
	controller.conf.RunAddress = srv.URL

	type RegisterModel struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	// какой результат хотим получить
	type want struct {
		code int
	}

	httpClient := resty.New()

	tests := []struct {
		name string
		body RegisterModel
		want want
	}{
		{
			name: "Success registry",
			body: RegisterModel{
				Login:    "testuser",
				Password: "testpassword",
			},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "Login is occupied by another user",
			body: RegisterModel{
				Login:    "testuser",
				Password: "testpassword",
			},
			want: want{
				code: http.StatusConflict,
			},
		},
		{
			name: "Wrong format data",
			body: RegisterModel{
				Login:    "",
				Password: "",
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			res, err := httpClient.R().SetBody(tt.body).Post(srv.URL + "/api/user/register")
			require.NoError(t, err)
			// проверяем код ответа
			assert.Equal(t, tt.want.code, res.StatusCode())

			// если пользователь зарегистрирован, проверяем его наличие
			// if tt.want.code == http.StatusOK {
			// 	_, err := controller.usecaseAuth.LoginAndGetUserJWT(ctx, tt.body.Login, tt.body.Password)
			// 	require.NoError(t, err)
			// }
		})
	}
}
