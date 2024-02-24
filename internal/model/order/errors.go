package order

import "fmt"

type ErrWrongFormatOfOrderID struct {
	OrderID int64
	Err     error
}

func NewErrWrongFormatOfOrderID(err error, orderID int64) *ErrWrongFormatOfOrderID {
	return &ErrWrongFormatOfOrderID{
		OrderID: orderID,
		Err:     err,
	}
}

func (e *ErrWrongFormatOfOrderID) Error() string {
	return fmt.Sprintf("wrong order id (%d) format: %s", e.OrderID, e.Err)
}

func (e *ErrWrongFormatOfOrderID) Unwrap() error {
	return e.Err
}

type ErrOrderAlreadyExists struct {
	OrderID int64
	Err     error
}

func NewErrOrderAlreadyExists(err error, orderID int64) *ErrOrderAlreadyExists {
	return &ErrOrderAlreadyExists{
		OrderID: orderID,
		Err:     err,
	}
}

func (e *ErrOrderAlreadyExists) Error() string {
	return fmt.Sprintf("order %d already exists: %s", e.OrderID, e.Err)
}

func (e *ErrOrderAlreadyExists) Unwrap() error {
	return e.Err
}

type ErrOrderBelongsToAnotherUser struct {
	OrderID int64
	Err     error
}

func NewErrOrderBelongsToAnotherUser(err error, orderID int64) *ErrOrderBelongsToAnotherUser {
	return &ErrOrderBelongsToAnotherUser{
		OrderID: orderID,
		Err:     err,
	}
}

func (e *ErrOrderBelongsToAnotherUser) Error() string {
	return fmt.Sprintf("order %d belongs to another user: %s", e.OrderID, e.Err)
}

func (e *ErrOrderBelongsToAnotherUser) Unwrap() error {
	return e.Err
}