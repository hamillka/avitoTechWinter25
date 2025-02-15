package unit_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hamillka/avitoTechWinter25/internal/repositories"
	"github.com/hamillka/avitoTechWinter25/internal/repositories/models"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInventoryRepository_GetInventoryByUserID(t *testing.T) {
	tests := []struct {
		name          string
		userID        int64
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedItems []*models.Inventory
		expectedErr   error
	}{
		{
			name:   "Success case",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "merch_id", "amount"}).
					AddRow(1, 1, 1, 5).
					AddRow(2, 1, 2, 3)
				mock.ExpectQuery("SELECT id, user_id, merch_id, amount FROM inventory WHERE user_id = \\$1").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedItems: []*models.Inventory{
				{ID: 1, UserID: 1, MerchID: 1, Amount: 5},
				{ID: 2, UserID: 1, MerchID: 2, Amount: 3},
			},
			expectedErr: nil,
		},
		{
			name:   "Empty result",
			userID: 2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "merch_id", "amount"})
				mock.ExpectQuery("SELECT id, user_id, merch_id, amount FROM inventory WHERE user_id = \\$1").
					WithArgs(2).
					WillReturnRows(rows)
			},
			expectedItems: []*models.Inventory{},
			expectedErr:   nil,
		},
		{
			name:   "Database error",
			userID: 3,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, user_id, merch_id, amount FROM inventory WHERE user_id = \\$1").
					WithArgs(3).
					WillReturnError(sql.ErrConnDone)
			},
			expectedItems: []*models.Inventory{},
			expectedErr:   sql.ErrConnDone,
		},
		{
			name:   "Scan error",
			userID: 4,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "merch_id", "amount"}).
					AddRow("invalid", 4, 5, 5)
				mock.ExpectQuery("SELECT id, user_id, merch_id, amount FROM inventory WHERE user_id = \\$1").
					WithArgs(4).
					WillReturnRows(rows)
			},
			expectedItems: nil,
			expectedErr:   repositories.ErrDatabaseReadingError,
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

			repo := repositories.NewInventoryRepository(sqlxDB)

			items, err := repo.GetInventoryByUserID(tt.userID)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedItems, items)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
