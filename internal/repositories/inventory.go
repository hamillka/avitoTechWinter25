package repositories

import (
	"github.com/hamillka/avitoTechWinter25/internal/repositories/models"
	"github.com/jmoiron/sqlx"
)

type InventoryRepository struct {
	db *sqlx.DB
}

const (
	getUserInventory = "SELECT * FROM inventory WHERE user_id = $1"
)

func NewInventoryRepository(db *sqlx.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (ir *InventoryRepository) GetInventoryByUserID(userID int64) ([]*models.Inventory, error) {
	var inventory []*models.Inventory
	rows, err := ir.db.Query(getUserInventory, userID)
	if err != nil {
		return nil, err
	}

	if err = rows.Err(); err != nil {
		return nil, ErrDatabaseReadingError
	}

	for rows.Next() {
		item := new(models.Inventory)
		if err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.MerchID,
			&item.Amount,
		); err != nil {
			return nil, ErrDatabaseReadingError
		}
		inventory = append(inventory, item)
	}
	defer rows.Close()

	return inventory, nil
}
