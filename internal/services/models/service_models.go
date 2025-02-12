package models

import (
	repoModels "github.com/hamillka/avitoTechWinter25/internal/repositories/models"
)

type User struct {
	ID       int64
	Username string
	Password string
	Coins    int64
}

type Transaction struct {
	ID           int64
	SenderName   int64
	ReceiverName int64
	Amount       int64
}

type Merch struct {
	ID   int64
	Type string
	Cost int64
}

type Inventory struct {
	ID        int64
	Username  int64
	MerchName int64
	Amount    int64
}

//двадвыда

// InventoryItem представляет предмет в инвентаре
type InventoryItem struct {
	Type     string
	Quantity int64
}

// IncomingTransactionInfo Информация о входящей транзакции
type IncomingTransactionInfo struct {
	FromUser string
	Amount   int64
}

// OutgoingTransactionInfo Информация об исходящей транзакции
type OutgoingTransactionInfo struct {
	ToUser string
	Amount int64
}

// CoinHistory История транзакций
type CoinHistory struct {
	Received []IncomingTransactionInfo // Входящие транзакции
	Sent     []OutgoingTransactionInfo // Исходящие транзакции
}

// Info представляет полную информацию о пользователе
type Info struct {
	Coins       int64
	Inventory   []InventoryItem
	CoinHistory CoinHistory
}

func ConvertUserToBL(user repoModels.User) *User {
	return &User{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Coins:    user.Coins,
	}
}
