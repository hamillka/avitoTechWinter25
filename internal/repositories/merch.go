package repositories

import (
	"database/sql"
	"errors"

	"github.com/hamillka/avitoTechWinter25/internal/repositories/models"
	"github.com/jmoiron/sqlx"
)

type MerchRepository struct {
	db *sqlx.DB
}

const (
	getMerchByID   = "SELECT id, type, cost FROM merch WHERE id = $1"
	getMerchByType = "SELECT id, type, cost FROM merch WHERE type = $1"
)

func NewMerchRepository(db *sqlx.DB) *MerchRepository {
	return &MerchRepository{db: db}
}

func (ir *MerchRepository) GetMerchByID(id int64) (models.Merch, error) {
	var merch models.Merch
	err := ir.db.QueryRow(getMerchByID, id).Scan(
		&merch.ID,
		&merch.Type,
		&merch.Cost,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Merch{}, ErrRecordNotFound
		}
		return models.Merch{}, err
	}

	return merch, nil
}

func (ir *MerchRepository) GetMerchByType(merchType string) (models.Merch, error) {
	var merch models.Merch
	err := ir.db.QueryRow(getMerchByType, merchType).Scan(
		&merch.ID,
		&merch.Type,
		&merch.Cost,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Merch{}, ErrRecordNotFound
		}
		return models.Merch{}, err
	}

	return merch, nil
}
