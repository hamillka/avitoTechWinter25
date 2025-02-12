package handlers

import "github.com/hamillka/avitoTechWinter25/internal/services/models"

type AvitoShopService interface {
	GetInfo(username string) (*models.Info, error)
	SendCoin(sender, receiver string, amount int64) error
	BuyItem(buyerName, itemName string) error
	Login(username, password string) (models.User, error)
}
