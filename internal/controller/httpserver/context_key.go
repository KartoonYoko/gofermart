package httpserver

// тип ключа контекста для middleware аутентификации
type controllerContextKey int

const (
	keyUserID controllerContextKey = iota // ключ для ID пользователя
)
