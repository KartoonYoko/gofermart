package httpserver

import (
	"context"
	"errors"
	"gofermart/internal/logger"
	model "gofermart/internal/model/auth"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

// middlewareAuth проверяет наличие куки "Authorization", валидирует находящийся там JWT токен
func (c *HTTPController) middlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userID model.UserID
		cookie, err := r.Cookie("Authorization")

		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			logger.Log.Error("middleware auth error: ", zap.Error(err))
			http.Error(w, "unexpected auth error", http.StatusInternalServerError)
			return
		}

		var cookieValue string
		arr := strings.Split(cookie.Value, " ")
		if len(arr) == 2 {
			cookieValue = arr[1]
		}
		userID, err = c.uc.ValidateJWTAndGetUserID(cookieValue)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), keyUserID, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
