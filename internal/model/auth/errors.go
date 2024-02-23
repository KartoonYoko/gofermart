package auth

import "errors"

// Логин уже занят
var ErrLoginIsOccupiedByAnotherUser = errors.New("auth usecase: login is occupied by another user")

// Пользователь не найден
var ErrUserNotFound = errors.New("auth usecase: user not found")

// Неверный формат данных
var ErrWrongDataFormat = errors.New("auth usecase: wrong data format")
