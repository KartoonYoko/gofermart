package httpserver

import "net/http"

// загрузка пользователем номера заказа для расчёта;
func (c *HttpController) handlerUserOrdersPOST(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	// ctxUserID := ctx.Value(keyUserID)
	// var request model.CreateShortenURLRequest
	// if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
	// 	http.Error(w, "Can not parse body", http.StatusBadRequest)
	// 	return
	// }

	// TODO сделать отдельную горутину для проверки статусов заказов;
	// при успешной обработке в одной ТРАНЗАКЦИИ изменить статус заказа, сохранить начисленные баллы
	// и прибавить эти баллы к существующем у пользователя на данный момент
	w.WriteHeader(http.StatusInternalServerError)
}

// получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
func (c *HttpController) handlerUserOrdersGET(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	// var request model.CreateShortenURLRequest
	// if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
	// 	http.Error(w, "Can not parse body", http.StatusBadRequest)
	// 	return
	// }
	w.WriteHeader(http.StatusInternalServerError)
}

// получение текущего баланса счёта баллов лояльности пользователя;
func (c *HttpController) handlerUserBalanceGET(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	// var request model.CreateShortenURLRequest
	// if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
	// 	http.Error(w, "Can not parse body", http.StatusBadRequest)
	// 	return
	// }
	w.WriteHeader(http.StatusInternalServerError)
}

// запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
func (c *HttpController) handlerUserBalanceWithdrawPOST(w http.ResponseWriter, r *http.Request) {
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
func (c *HttpController) handlerUserWithdrawalsGET(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	// var request model.CreateShortenURLRequest
	// if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
	// 	http.Error(w, "Can not parse body", http.StatusBadRequest)
	// 	return
	// }
	w.WriteHeader(http.StatusInternalServerError)
}
