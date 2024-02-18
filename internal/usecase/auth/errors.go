package auth

import "errors"

var ErrLoginIsOccupiedByAnotherUser = errors.New("auth usecase: login is occupied by another user")
var ErrWrongDataFormat = errors.New("auth usecase: wrong data format")
