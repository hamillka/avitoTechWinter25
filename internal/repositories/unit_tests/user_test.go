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

func TestUserRepository_GetUserByUsernamePassword(t *testing.T) {
	tests := []struct {
		name        string
		username    string
		password    string
		mockSetup   func(mock sqlmock.Sqlmock)
		expected    models.User
		expectedErr error
	}{
		{
			name:     "Success case",
			username: "testuser",
			password: "testpass",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "username", "password", "coins"}).
					AddRow(1, "testuser", "testpass", 1000)
				mock.ExpectQuery("SELECT id, username, password, coins FROM users WHERE username = \\$1 AND password = \\$2").
					WithArgs("testuser", "testpass").
					WillReturnRows(rows)
			},
			expected: models.User{
				ID:       1,
				Username: "testuser",
				Password: "testpass",
				Coins:    1000,
			},
			expectedErr: nil,
		},
		{
			name:     "User not found",
			username: "nonexistent",
			password: "wrongpass",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, username, password, coins FROM users WHERE username = \\$1 AND password = \\$2").
					WithArgs("nonexistent", "wrongpass").
					WillReturnError(sql.ErrNoRows)
			},
			expected:    models.User{},
			expectedErr: repositories.ErrRecordNotFound,
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

			repo := repositories.NewUserRepository(sqlxDB)
			result, err := repo.GetUserByUsernamePassword(tt.username, tt.password)

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

func TestUserRepository_CreateUser(t *testing.T) {
	tests := []struct {
		name        string
		username    string
		password    string
		mockSetup   func(mock sqlmock.Sqlmock)
		expected    models.User
		expectedErr error
	}{
		{
			name:     "Success case",
			username: "newuser",
			password: "newpass",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO users \\(username, password\\) VALUES \\(\\$1, \\$2\\) RETURNING id").
					WithArgs("newuser", "newpass").
					WillReturnRows(rows)
			},
			expected: models.User{
				ID:       1,
				Username: "newuser",
				Password: "newpass",
				Coins:    repositories.DefaultCoins,
			},
			expectedErr: nil,
		},
		{
			name:     "User already exists",
			username: "existing",
			password: "pass",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("INSERT INTO users \\(username, password\\) VALUES \\(\\$1, \\$2\\) RETURNING id").
					WithArgs("existing", "pass").
					WillReturnError(sql.ErrConnDone)
			},
			expected:    models.User{},
			expectedErr: repositories.ErrRecordAlreadyExists,
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

			repo := repositories.NewUserRepository(sqlxDB)
			result, err := repo.CreateUser(tt.username, tt.password)

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

func TestUserRepository_TransferCoins(t *testing.T) {
	tests := []struct {
		name        string
		senderID    int64
		receiverID  int64
		amount      int64
		mockSetup   func(mock sqlmock.Sqlmock)
		expectedErr error
	}{
		{
			name:       "Success case",
			senderID:   1,
			receiverID: 2,
			amount:     100,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE users SET coins = coins - \\$1 WHERE id = \\$2").
					WithArgs(100, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec("UPDATE users SET coins = coins \\+ \\$1 WHERE id = \\$2").
					WithArgs(100, 2).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectQuery("INSERT INTO transactions \\(sender_id, receiver_id, amount\\) VALUES \\(\\$1, \\$2, \\$3\\) RETURNING id").
					WithArgs(1, 2, 100).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			expectedErr: nil,
		},
		{
			name:       "Failed to subtract coins",
			senderID:   1,
			receiverID: 2,
			amount:     100,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE users SET coins = coins - \\$1 WHERE id = \\$2").
					WithArgs(100, 1).
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			expectedErr: repositories.ErrDatabaseUpdatingError,
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

			repo := repositories.NewUserRepository(sqlxDB)
			err = repo.TransferCoins(tt.senderID, tt.receiverID, tt.amount)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_BuyItemFromAvitoShop(t *testing.T) {
	tests := []struct {
		name        string
		buyerID     int64
		itemID      int64
		itemCost    int64
		mockSetup   func(mock sqlmock.Sqlmock)
		expectedErr error
	}{
		{
			name:     "Success case - new item",
			buyerID:  1,
			itemID:   100,
			itemCost: 500,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE users SET coins = coins - \\$1 WHERE id = \\$2").
					WithArgs(500, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec("UPDATE users SET coins = coins \\+ \\$1 WHERE id = \\$2").
					WithArgs(500, repositories.AvitoShopID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectQuery("INSERT INTO transactions \\(sender_id, receiver_id, amount\\) VALUES \\(\\$1, \\$2, \\$3\\) RETURNING id").
					WithArgs(1, repositories.AvitoShopID, 500).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectQuery("SELECT 1 FROM inventory WHERE user_id = \\$1 AND merch_id = \\$2").
					WithArgs(1, 100).
					WillReturnError(sql.ErrNoRows)
				mock.ExpectQuery("INSERT INTO inventory \\(user_id, merch_id, amount\\) VALUES \\(\\$1, \\$2, 1\\) RETURNING id").
					WithArgs(1, 100).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			expectedErr: nil,
		},
		{
			name:     "Success case - existing item",
			buyerID:  1,
			itemID:   100,
			itemCost: 500,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE users SET coins = coins - \\$1 WHERE id = \\$2").
					WithArgs(500, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec("UPDATE users SET coins = coins \\+ \\$1 WHERE id = \\$2").
					WithArgs(500, repositories.AvitoShopID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectQuery("INSERT INTO transactions \\(sender_id, receiver_id, amount\\) VALUES \\(\\$1, \\$2, \\$3\\) RETURNING id").
					WithArgs(1, repositories.AvitoShopID, 500).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectQuery("SELECT 1 FROM inventory").
					WithArgs(1, 100).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
				mock.ExpectExec("UPDATE inventory SET amount = amount \\+ 1").
					WithArgs(1, 100).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			expectedErr: nil,
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

			repo := repositories.NewUserRepository(sqlxDB)
			err = repo.BuyItemFromAvitoShop(tt.buyerID, tt.itemID, tt.itemCost)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				require.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
