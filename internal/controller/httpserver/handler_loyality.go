package httpserver

import (
	"encoding/json"
	"errors"
	"gofermart/internal/model/auth"
	modelOrder "gofermart/internal/model/order"
	"gofermart/pkg/luhn"
	"io"
	"net/http"
)

// handlerUserOrdersPOST загрузка пользователем номера заказа для расчёта
//
// Возможные коды ответа:
//
//	200 — номер заказа уже был загружен этим пользователем;
//	202 — новый номер заказа принят в обработку;
//	400 — неверный формат запроса;
//	401 — пользователь не аутентифицирован;
//	409 — номер заказа уже был загружен другим пользователем;
//	422 — неверный формат номера заказа;
//	500 — внутренняя ошибка сервера.
func (c *HTTPController) handlerUserOrdersPOST(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Can not parse body", http.StatusBadRequest)
		return
	}
	orderID, err := luhn.CheckStr(string(body))
	if err != nil {
		http.Error(w, "Wrong order id format", http.StatusUnprocessableEntity)
		return
	}
	ctxUserID := ctx.Value(keyUserID)
	userID, ok := ctxUserID.(auth.UserID)
	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	err = c.usecaseOrder.CreateNewOrder(ctx, userID, orderID)
	if err != nil {
		var errOrderAlreadyExists *modelOrder.ErrOrderAlreadyExists
		if errors.As(err, &errOrderAlreadyExists) {
			http.Error(w, "", http.StatusOK)
			return
		}
		var errOrderBelongsToAnotherUser *modelOrder.ErrOrderBelongsToAnotherUser
		if errors.As(err, &errOrderBelongsToAnotherUser) {
			http.Error(w, "", http.StatusConflict)
			return
		}

		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// handlerUserOrdersGET - получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
//
// Ответы:
//
//	200 — успешная обработка запроса.
//	204 — нет данных для ответа.
//	401 — пользователь не авторизован.
//	500 — внутренняя ошибка сервера.
func (c *HTTPController) handlerUserOrdersGET(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctxUserID := ctx.Value(keyUserID)
	userID, ok := ctxUserID.(auth.UserID)
	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	res, err := c.usecaseOrder.GetUserOrders(ctx, userID)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if len(res) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	jsonStr, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonStr))
}

// получение текущего баланса счёта баллов лояльности пользователя;
func (c *HTTPController) handlerUserBalanceGET(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	// var request model.CreateShortenURLRequest
	// if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
	// 	http.Error(w, "Can not parse body", http.StatusBadRequest)
	// 	return
	// }
	w.WriteHeader(http.StatusInternalServerError)
}

// запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
func (c *HTTPController) handlerUserBalanceWithdrawPOST(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	// var request model.CreateShortenURLRequest
	// if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
	// 	http.Error(w, "Can not parse body", http.StatusBadRequest)
	// 	return
	// }

	// TODO проверить достаточно ли баллов на счету;
	// в одной ТРАНЗАКЦИИ списать баллы и добавить записб о списании в таблицу списания баллов;
	// нужно добавлять заказ к себе в БД;
	w.WriteHeader(http.StatusInternalServerError)
}

// получение информации о выводе средств с накопительного счёта пользователем
func (c *HTTPController) handlerUserWithdrawalsGET(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	// var request model.CreateShortenURLRequest
	// if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
	// 	http.Error(w, "Can not parse body", http.StatusBadRequest)
	// 	return
	// }
	w.WriteHeader(http.StatusInternalServerError)
}
