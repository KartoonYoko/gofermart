package withdraw

import "errors"

// ErrUserHasNotEnoughBalance - сигнализирует, что на счету пользователя недостаточно средств
var ErrUserHasNotEnoughBalance = errors.New("user has not enough balance")
