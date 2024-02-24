package orderaccrual

import "errors"

// заказ не анйден в системе расчёта баллов
var ErrOrderNotFound = errors.New("order not found")

// слишком много запросов к системе расчёта баллов
var ErrTooManyRequests = errors.New("too many requests")

// статус код не соответсвующий ожиданиям от АПИ
var ErrUndefinedStatusCode = errors.New("undefined status code")