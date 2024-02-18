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
