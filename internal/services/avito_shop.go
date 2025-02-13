package services

import (
	"errors"

	"github.com/hamillka/avitoTechWinter25/internal/handlers/dto"
	"github.com/hamillka/avitoTechWinter25/internal/repositories"
	repoModels "github.com/hamillka/avitoTechWinter25/internal/repositories/models"
	serviceModels "github.com/hamillka/avitoTechWinter25/internal/services/models"
)

type UserRepository interface {
	GetUserByUsernamePassword(username, password string) (repoModels.User, error)
	CreateUser(username, password string) (repoModels.User, error)
	GetUserByUsername(username string) (repoModels.User, error)
	GetUserByID(id int64) (repoModels.User, error)
	TransferCoins(senderID, receiverID, amount int64) error
	BuyItemFromAvitoShop(buyerID, itemID, itemCost int64) error
}

type InventoryRepository interface {
	GetInventoryByUserID(userID int64) ([]*repoModels.Inventory, error)
}

type MerchRepository interface {
	GetMerchByID(id int64) (repoModels.Merch, error)
	GetMerchByType(merchType string) (repoModels.Merch, error)
}

type TransactionRepository interface {
	GetOutTransactions(userID int64) ([]*repoModels.Transaction, error)
	GetInTransactions(userID int64) ([]*repoModels.Transaction, error)
}

type AvitoShopService struct {
	userRepo        UserRepository
	inventoryRepo   InventoryRepository
	merchRepo       MerchRepository
	transactionRepo TransactionRepository
}

func NewAvitoShopService(
	userRepository UserRepository,
	inventoryRepository InventoryRepository,
	merchRepository MerchRepository,
	transactionRepository TransactionRepository,
) *AvitoShopService {
	return &AvitoShopService{
		userRepo:        userRepository,
		inventoryRepo:   inventoryRepository,
		merchRepo:       merchRepository,
		transactionRepo: transactionRepository,
	}
}

func (as *AvitoShopService) GetInfo(username string) (*serviceModels.Info, error) {
	user, err := as.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	inventory, err := as.inventoryRepo.GetInventoryByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	inventoryMap := make(map[string]int64)
	for _, item := range inventory {
		merch, err := as.merchRepo.GetMerchByID(item.MerchID)
		if err != nil {
			continue
		}
		inventoryMap[merch.Type] += item.Amount
	}

	inventoryItems := make([]serviceModels.InventoryItem, 0)
	for itemType, quantity := range inventoryMap {
		inventoryItems = append(inventoryItems, serviceModels.InventoryItem{
			Type:     itemType,
			Quantity: quantity,
		})
	}

	outTransactions, err := as.transactionRepo.GetOutTransactions(user.ID)
	if err != nil {
		return nil, err
	}

	sentTransactions := make([]serviceModels.OutgoingTransactionInfo, 0)
	for _, tx := range outTransactions {
		receiver, err := as.userRepo.GetUserByID(tx.ReceiverID)
		if err != nil {
			continue
		}
		sentTransactions = append(sentTransactions, serviceModels.OutgoingTransactionInfo{
			ToUser: receiver.Username,
			Amount: tx.Amount,
		})
	}

	inTransactions, err := as.transactionRepo.GetInTransactions(user.ID)
	if err != nil {
		return nil, err
	}

	receivedTransactions := make([]serviceModels.IncomingTransactionInfo, 0)
	for _, tx := range inTransactions {
		sender, err := as.userRepo.GetUserByID(tx.SenderID)
		if err != nil {
			continue
		}
		receivedTransactions = append(receivedTransactions, serviceModels.IncomingTransactionInfo{
			FromUser: sender.Username,
			Amount:   tx.Amount,
		})
	}

	info := &serviceModels.Info{
		Coins:     user.Coins,
		Inventory: inventoryItems,
		CoinHistory: serviceModels.CoinHistory{
			Received: receivedTransactions,
			Sent:     sentTransactions,
		},
	}

	return info, nil
}

func (as *AvitoShopService) SendCoin(sender, receiver string, amount int64) error {
	senderUser, err := as.userRepo.GetUserByUsername(sender)
	if err != nil {
		return err
	}

	receiverUser, err := as.userRepo.GetUserByUsername(receiver)
	if err != nil {
		return err
	}

	if senderUser.Coins < amount {
		return dto.ErrNotEnoughCoins
	}

	err = as.userRepo.TransferCoins(senderUser.ID, receiverUser.ID, amount)
	if err != nil {
		return err
	}

	return nil
}

func (as *AvitoShopService) BuyItem(buyerName, itemName string) error {
	merch, err := as.merchRepo.GetMerchByType(itemName)
	if err != nil {
		return err
	}

	buyerUser, err := as.userRepo.GetUserByUsername(buyerName)
	if err != nil {
		return err
	}

	if merch.Cost > buyerUser.Coins {
		return dto.ErrNotEnoughCoins
	}

	err = as.userRepo.BuyItemFromAvitoShop(buyerUser.ID, merch.ID, merch.Cost)
	if err != nil {
		return err
	}

	return nil
}

func (as *AvitoShopService) Login(username, password string) (serviceModels.User, error) {
	user, err := as.userRepo.GetUserByUsernamePassword(username, password)
	if errors.Is(err, repositories.ErrRecordNotFound) {
		user, err = as.userRepo.CreateUser(username, password)
		if err != nil {
			return serviceModels.User{}, err
		}
	}

	return *serviceModels.ConvertUserToBL(user), nil
}
