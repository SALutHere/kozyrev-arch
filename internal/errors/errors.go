package errors

import "errors"

var (
	ErrNotFound       = errors.New("деталь не найдена")
	ErrNotEnoughParts = errors.New("недостаточно деталей")
)
