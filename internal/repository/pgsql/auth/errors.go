package auth

import "fmt"

type ErrLoginAlreadyExists struct {
	Login string
	Err   error
}

func NewErrLoginAlreadyExists(login string, err error) *ErrLoginAlreadyExists {
	return &ErrLoginAlreadyExists{
		Login: login,
		Err:   err,
	}
}

func (e *ErrLoginAlreadyExists) Error() string {
	return fmt.Sprintf("login %s already exists: %s", e.Login, e.Err)
}

func (e *ErrLoginAlreadyExists) Unwrap() error {
	return e.Err
}

type ErrUserNotFound struct {
	Login    string
	Password string
	Err      error
}

func NewErrUserNotFound(login string, password string, err error) *ErrLoginAlreadyExists {
	return &ErrLoginAlreadyExists{
		Login: login,
		Err:   err,
	}
}

func (e *ErrUserNotFound) Error() string {
	return fmt.Sprintf("login %s already exists: %s", e.Login, e.Err)
}

func (e *ErrUserNotFound) Unwrap() error {
	return e.Err
}
