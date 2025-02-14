package dto

import "errors"

// ErrorDto model info
// @Description Информация об ошибке (DTO)
type ErrorDto struct {
	Error string `json:"error"` // Ошибка
}

var (
	ErrNotEnoughCoins    = errors.New("not enough coins error")
	ErrPathVarExtracting = errors.New("item type extracting error")
)
