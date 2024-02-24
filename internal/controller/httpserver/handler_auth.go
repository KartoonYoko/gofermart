package httpserver

import (
	"encoding/json"
	"errors"
	"fmt"
	model "gofermart/internal/model/auth"
	"net/http"
)

// handlerUserRegisterPOST - регистрация пользователя
//
// Ответы:
//
//	200 — пользователь успешно зарегистрирован и аутентифицирован;
//	400 — неверный формат запроса;
//	409 — логин уже занят;
//	500 — внутренняя ошибка сервера.
func (c *HTTPController) handlerUserRegisterPOST(w http.ResponseWriter, r *http.Request) {
	type RegisterModel struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	ctx := r.Context()
	var request RegisterModel
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Can not parse body", http.StatusBadRequest)
		return
	}

	jwt, err := c.usecaseAuth.RegisterAndGetUserJWT(ctx, request.Login, request.Password)
	if err != nil {
		if errors.Is(err, model.ErrWrongDataFormat) {
			http.Error(w, "Wrong data format", http.StatusBadRequest)
			return
		}
		if errors.Is(err, model.ErrLoginIsOccupiedByAnotherUser) {
			http.Error(w, "Login is occupied", http.StatusConflict)
			return
		}

		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	saveJWTAsAuthorizationCookie(w, jwt)
	w.WriteHeader(http.StatusOK)
}

// handlerUserLoginPOST - аутентификация пользователя
//
// Ответы:
//
//	200 — пользователь успешно аутентифицирован;
//	400 — неверный формат запроса;
//	401 — неверная пара логин/пароль;
//	500 — внутренняя ошибка сервера.
func (c *HTTPController) handlerUserLoginPOST(w http.ResponseWriter, r *http.Request) {
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

	jwt, err := c.usecaseAuth.LoginAndGetUserJWT(ctx, request.Login, request.Password)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			http.Error(w, "Wrong login or password", http.StatusUnauthorized)
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
