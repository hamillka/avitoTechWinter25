package repositories

import (
	"github.com/hamillka/avitoTechWinter25/internal/repositories/models"
	"github.com/jmoiron/sqlx"
)

type TransactionsRepository struct {
	db *sqlx.DB
}

const (
	getOutTransactionsByUID = "SELECT * FROM transactions WHERE sender_id = $1"
	getInTransactionsByUID  = "SELECT * FROM transactions WHERE receiver_id = $1"
)

func NewTransactionsRepository(db *sqlx.DB) *TransactionsRepository {
	return &TransactionsRepository{db: db}
}

func (tr *TransactionsRepository) GetOutTransactions(userID int64) ([]*models.Transaction, error) {
	transactions := make([]*models.Transaction, 0)
	rows, _ := tr.db.Query(getOutTransactionsByUID, userID)

	if err := rows.Err(); err != nil {
		return transactions, ErrDatabaseReadingError
	}

	for rows.Next() {
		transaction := new(models.Transaction)
		if err := rows.Scan(
			&transaction.ID,
			&transaction.SenderID,
			&transaction.ReceiverID,
			&transaction.Amount,
		); err != nil {
			return nil, ErrDatabaseReadingError
		}
		transactions = append(transactions, transaction)
	}
	defer rows.Close()

	return transactions, nil
}

func (tr *TransactionsRepository) GetInTransactions(userID int64) ([]*models.Transaction, error) {
	transactions := make([]*models.Transaction, 0)
	rows, _ := tr.db.Query(getInTransactionsByUID, userID)

	if err := rows.Err(); err != nil {
		return transactions, ErrDatabaseReadingError
	}

	for rows.Next() {
		transaction := new(models.Transaction)
		if err := rows.Scan(
			&transaction.ID,
			&transaction.SenderID,
			&transaction.ReceiverID,
			&transaction.Amount,
		); err != nil {
			return nil, ErrDatabaseReadingError
		}
		transactions = append(transactions, transaction)
	}
	defer rows.Close()

	return transactions, nil
}
