package repositories

import (
	"database/sql"
	"errors"

	"github.com/hamillka/avitoTechWinter25/internal/repositories/models"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

const (
	createUser                = "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id"
	getUserByID               = "SELECT * FROM users WHERE id = $1"
	getUserByUsernamePassword = "SELECT * FROM users WHERE username = $1 AND password = $2"
	getUser                   = "SELECT * FROM users WHERE username = $1"
	subtractCoins             = "UPDATE users SET coins = coins - $1 WHERE id = $2"
	addCoins                  = "UPDATE users SET coins = coins + $1 WHERE id = $2"
	addTransaction            = "INSERT INTO transactions (sender_id, receiver_id, amount) VALUES ($1, $2, $3) RETURNING id"
	inventoryItemExists       = "SELECT 1 FROM inventory WHERE user_id = $1 AND merch_id = $2"
	incrementItemAmount       = "UPDATE inventory SET amount = amount + 1 WHERE user_id = $1 AND merch_id = $2"
	addItemToInventory        = "INSERT INTO inventory (user_id, merch_id, amount) VALUES ($1, $2, 1) RETURNING id"
)

const (
	AvitoShopID  = 0
	DefaultCoins = 1000
)

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) GetUserByUsernamePassword(username, password string) (models.User, error) {
	var user models.User
	err := ur.db.QueryRow(getUserByUsernamePassword, username, password).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Coins,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, ErrRecordNotFound
		}
		return models.User{}, err
	}

	return user, nil
}

func (ur *UserRepository) CreateUser(username, password string) (models.User, error) {
	var id int64
	err := ur.db.QueryRow(createUser, username, password).Scan(&id)
	if err != nil {
		return models.User{}, ErrRecordAlreadyExists
	}

	return models.User{
		ID:       id,
		Username: username,
		Password: password,
		Coins:    DefaultCoins,
	}, nil
}

func (ur *UserRepository) GetUserByUsername(username string) (models.User, error) {
	var user models.User
	err := ur.db.QueryRow(getUser, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Coins,
	)
	if err != nil {
		return models.User{}, ErrRecordNotFound
	}

	return user, nil
}

func (ur *UserRepository) GetUserByID(id int64) (models.User, error) {
	var user models.User
	err := ur.db.QueryRow(getUserByID, id).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Coins,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, ErrRecordNotFound
		}
		return models.User{}, err
	}

	return user, nil
}

func (ur *UserRepository) TransferCoins(senderID, receiverID, amount int64) error {
	tx, err := ur.db.Beginx()
	if err != nil {
		return ErrStartTxError
	}

	_, err = ur.db.Exec(subtractCoins, amount, senderID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return ErrRollbackTxError
		}
		return ErrDatabaseUpdatingError
	}

	_, err = tx.Exec(addCoins, amount, receiverID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return ErrRollbackTxError
		}
		return ErrDatabaseUpdatingError
	}

	var txID int64
	err = ur.db.QueryRow(addTransaction, senderID, receiverID, amount).Scan(&txID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return ErrRollbackTxError
		}
		return ErrDatabaseUpdatingError
	}

	if err := tx.Commit(); err != nil {
		return ErrCommitError
	}

	return nil
}

func (ur *UserRepository) BuyItemFromAvitoShop(buyerID, itemID, itemCost int64) error {
	tx, err := ur.db.Beginx()
	if err != nil {
		return ErrStartTxError
	}

	_, err = ur.db.Exec(subtractCoins, itemCost, buyerID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return ErrRollbackTxError
		}
		return ErrDatabaseUpdatingError
	}

	_, err = tx.Exec(addCoins, itemCost, AvitoShopID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return ErrRollbackTxError
		}
		return ErrDatabaseUpdatingError
	}

	var transactionID int64
	err = ur.db.QueryRow(addTransaction, buyerID, AvitoShopID, itemCost).Scan(&transactionID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return ErrRollbackTxError
		}
		return ErrDatabaseUpdatingError
	}

	var exists bool
	_ = tx.QueryRow(inventoryItemExists, buyerID, itemID).Scan(&exists)
	if exists {
		_, err = tx.Exec(incrementItemAmount, buyerID, itemID)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return ErrRollbackTxError
			}
			return ErrDatabaseUpdatingError
		}
	} else {
		var newItemID int64
		err = ur.db.QueryRow(addItemToInventory, buyerID, itemID).Scan(&newItemID)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return ErrRollbackTxError
			}
			return ErrDatabaseUpdatingError
		}
	}

	if err := tx.Commit(); err != nil {
		return ErrCommitError
	}

	return nil
}
