package dto

type SendCoinRequestDto struct {
	ToUser string `json:"toUser"`
	Amount int64  `json:"amount"`
}
