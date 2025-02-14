package unit_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hamillka/avitoTechWinter25/internal/repositories"
	"github.com/hamillka/avitoTechWinter25/internal/repositories/models"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionsRepository_GetOutTransactions(t *testing.T) {
	tests := []struct {
		name        string
		userID      int64
		mockSetup   func(mock sqlmock.Sqlmock)
		expected    []*models.Transaction
		expectedErr error
	}{
		{
			name:   "Success case with multiple transactions",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "sender_id", "receiver_id", "amount"}).
					AddRow(1, 1, 2, 100).
					AddRow(2, 1, 3, 200)
				mock.ExpectQuery("SELECT \\* FROM transactions WHERE sender_id = \\$1").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expected: []*models.Transaction{
				{ID: 1, SenderID: 1, ReceiverID: 2, Amount: 100},
				{ID: 2, SenderID: 1, ReceiverID: 3, Amount: 200},
			},
			expectedErr: nil,
		},
		{
			name:   "Empty result",
			userID: 2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "sender_id", "receiver_id", "amount"})
				mock.ExpectQuery("SELECT \\* FROM transactions WHERE sender_id = \\$1").
					WithArgs(2).
					WillReturnRows(rows)
			},
			expected:    []*models.Transaction{},
			expectedErr: nil,
		},
		{
			name:   "Scan error",
			userID: 4,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "sender_id", "receiver_id", "amount"}).
					AddRow("invalid", 4, 5, 100)
				mock.ExpectQuery("SELECT \\* FROM transactions WHERE sender_id = \\$1").
					WithArgs(4).
					WillReturnRows(rows)
			},
			expected:    nil,
			expectedErr: repositories.ErrDatabaseReadingError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			defer sqlxDB.Close()

			tt.mockSetup(mock)

			repo := repositories.NewTransactionsRepository(sqlxDB)
			result, err := repo.GetOutTransactions(tt.userID)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTransactionsRepository_GetInTransactions(t *testing.T) {
	tests := []struct {
		name        string
		userID      int64
		mockSetup   func(mock sqlmock.Sqlmock)
		expected    []*models.Transaction
		expectedErr error
	}{
		{
			name:   "Success case with multiple transactions",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "sender_id", "receiver_id", "amount"}).
					AddRow(1, 2, 1, 100).
					AddRow(2, 3, 1, 200)
				mock.ExpectQuery("SELECT \\* FROM transactions WHERE receiver_id = \\$1").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expected: []*models.Transaction{
				{ID: 1, SenderID: 2, ReceiverID: 1, Amount: 100},
				{ID: 2, SenderID: 3, ReceiverID: 1, Amount: 200},
			},
			expectedErr: nil,
		},
		{
			name:   "Empty result",
			userID: 2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "sender_id", "receiver_id", "amount"})
				mock.ExpectQuery("SELECT \\* FROM transactions WHERE receiver_id = \\$1").
					WithArgs(2).
					WillReturnRows(rows)
			},
			expected:    []*models.Transaction{},
			expectedErr: nil,
		},
		{
			name:   "Scan error",
			userID: 4,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "sender_id", "receiver_id", "amount"}).
					AddRow("invalid", 5, 4, 100)
				mock.ExpectQuery("SELECT \\* FROM transactions WHERE receiver_id = \\$1").
					WithArgs(4).
					WillReturnRows(rows)
			},
			expected:    nil,
			expectedErr: repositories.ErrDatabaseReadingError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			defer sqlxDB.Close()

			tt.mockSetup(mock)

			repo := repositories.NewTransactionsRepository(sqlxDB)
			result, err := repo.GetInTransactions(tt.userID)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
