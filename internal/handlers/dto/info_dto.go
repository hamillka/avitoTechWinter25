package dto

import "github.com/hamillka/avitoTechWinter25/internal/services/models"

// InventoryItemDto model info
// @Description Информация о предмете инвентаря
type InventoryItemDto struct {
	Type     string `json:"type"`     // Тип мерча
	Quantity int64  `json:"quantity"` // Количество мерча данного типа в инвентаре
}

// IncomingTransactionInfoDto model info
// @Description Информация о входящих транзакциях пользователя
type IncomingTransactionInfoDto struct {
	FromUser string `json:"fromUser"` // Отправитель
	Amount   int64  `json:"amount"`   // Количество полученных монет
}

// OutgoingTransactionInfoDto model info
// @Description Информация о исходящих транзакциях пользователя
type OutgoingTransactionInfoDto struct {
	ToUser string `json:"toUser"` // Получатель
	Amount int64  `json:"amount"` // Количество отправленных монет
}

// CoinHistoryDto model info
// @Description Объединенная информация о всех транзакциях пользователя
type CoinHistoryDto struct {
	Received []IncomingTransactionInfoDto `json:"received"` // Входящие транзакции
	Sent     []OutgoingTransactionInfoDto `json:"sent"`     // Исходящие транзакции
}

// InfoResponseDto model info
// @Description Полная информация о пользователе: количество монет, предметы инвентаря и история перемещения монет
type InfoResponseDto struct {
	Coins       int64              `json:"coins"`       // Количество монет
	Inventory   []InventoryItemDto `json:"inventory"`   // Инвентарь
	CoinHistory CoinHistoryDto     `json:"coinHistory"` // История перемещения монет
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
