package dto

// SendCoinRequestDto model info
// @Description Информация о получателе и количестве монет при переводе
type SendCoinRequestDto struct {
	ToUser string `json:"toUser"` // имя пользователя-получателя
	Amount int64  `json:"amount"` // количество переводимых монет
}
