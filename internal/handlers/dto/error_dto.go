package dto

import "errors"

type ErrorDto struct {
	Error string `json:"error"`
}

var (
	ErrNotEnoughCoins = errors.New("not enough coins error")
)
