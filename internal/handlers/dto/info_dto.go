package dto

import "github.com/hamillka/avitoTechWinter25/internal/services/models"

type InventoryItemDto struct {
	Type     string `json:"type"`
	Quantity int64  `json:"quantity"`
}

type IncomingTransactionInfoDto struct {
	FromUser string `json:"fromUser"` // Будет отображаться в JSON как "fromUser"
	Amount   int64  `json:"amount"`
}

type OutgoingTransactionInfoDto struct {
	ToUser string `json:"toUser"` // Будет отображаться в JSON как "toUser"
	Amount int64  `json:"amount"`
}

type CoinHistoryDto struct {
	Received []IncomingTransactionInfoDto `json:"received"` // Входящие транзакции
	Sent     []OutgoingTransactionInfoDto `json:"sent"`     // Исходящие транзакции
}

type InfoResponseDto struct {
	Coins       int64              `json:"coins"`
	Inventory   []InventoryItemDto `json:"inventory"`
	CoinHistory CoinHistoryDto     `json:"coinHistory"`
}

func ConvertBLInventoryToDto(item models.InventoryItem) *InventoryItemDto {
	return &InventoryItemDto{
		Type:     item.Type,
		Quantity: item.Quantity,
	}
}

func ConvertBLIncomingTxToDto(inTx models.IncomingTransactionInfo) *IncomingTransactionInfoDto {
	return &IncomingTransactionInfoDto{
		FromUser: inTx.FromUser,
		Amount:   inTx.Amount,
	}
}

func ConvertBLOutgoingTxToDto(outTx models.OutgoingTransactionInfo) *OutgoingTransactionInfoDto {
	return &OutgoingTransactionInfoDto{
		ToUser: outTx.ToUser,
		Amount: outTx.Amount,
	}
}

func ConvertBLInfoToDto(info *models.Info) *InfoResponseDto {
	inventoryDto := make([]InventoryItemDto, len(info.Inventory))
	for idx, item := range info.Inventory {
		inventoryDto[idx] = *ConvertBLInventoryToDto(item)
	}

	outTx := make([]OutgoingTransactionInfoDto, len(info.CoinHistory.Sent))
	inTx := make([]IncomingTransactionInfoDto, len(info.CoinHistory.Received))

	for idx, tx := range info.CoinHistory.Sent {
		outTx[idx] = *ConvertBLOutgoingTxToDto(tx)
	}
	for idx, tx := range info.CoinHistory.Received {
		inTx[idx] = *ConvertBLIncomingTxToDto(tx)
	}

	return &InfoResponseDto{
		Coins:     info.Coins,
		Inventory: inventoryDto,
		CoinHistory: CoinHistoryDto{
			Received: inTx,
			Sent:     outTx,
		},
	}
}
