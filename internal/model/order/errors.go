package order

import "fmt"

type ErrWrongFormatOfOrderID struct {
	OrderID int
	Err     error
}

func NewErrWrongFormatOfOrderID(err error, orderID int) *ErrWrongFormatOfOrderID {
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
