package order

import "errors"

// Заказ уже был создан пользователем
var ErrOrderAlreadyExists = errors.New("order repository: order was created by user")

// Заказ создан другому пользователем
var ErrOrderBelongsToAnotherUser = errors.New("order repository: order was created by another user")
