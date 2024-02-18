package http

import "net/http"

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

// аутентификация пользователя
func (c *HttpController) handlerUserLoginPOST(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	// var request model.CreateShortenURLRequest
	// if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
	// 	http.Error(w, "Can not parse body", http.StatusBadRequest)
	// 	return
	// }
	w.WriteHeader(http.StatusInternalServerError)
}
