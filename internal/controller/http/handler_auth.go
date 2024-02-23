package http

import (
	"encoding/json"
	"errors"
	"fmt"
	model "gofermart/internal/model/auth"
	"net/http"
)

// регистрация пользователя
func (c *HttpController) handlerUserRegisterPOST(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	// var request model.CreateShortenURLRequest
	// if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
	// 	http.Error(w, "Can not parse body", http.StatusBadRequest)
	// 	return
	// }
	w.WriteHeader(http.StatusInternalServerError)
}

// handlerUserLoginPOST - аутентификация пользователя
//
// Ответы:
//
//	200 — пользователь успешно аутентифицирован;
//	400 — неверный формат запроса;
//	401 — неверная пара логин/пароль;
//	500 — внутренняя ошибка сервера.
func (c *HttpController) handlerUserLoginPOST(w http.ResponseWriter, r *http.Request) {
	type LoginModel struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	ctx := r.Context()
	var request LoginModel
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Can not parse body", http.StatusBadRequest)
		return
	}

	jwt, err := c.uc.Login(ctx, request.Login, request.Password)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			http.Error(w, "Can not parse body", http.StatusUnauthorized)
			return
		}

		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	saveJWTAsAuthorizationCookie(w, jwt)
	w.WriteHeader(http.StatusOK)
}

// saveJWTAsAuthorizationCookie сохранит JWT строку в куку "Authorization"
func saveJWTAsAuthorizationCookie(w http.ResponseWriter, jwt string) {
	bearerStr := fmt.Sprintf("Bearer %s", jwt)

	cookie := http.Cookie{
		Name:     "Authorization",
		Value:    bearerStr,
		Path:     "/",
		MaxAge:   0,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
}
