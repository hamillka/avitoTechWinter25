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

func TestMerchRepository_GetMerchByID(t *testing.T) {
	tests := []struct {
		name        string
		merchID     int64
		mockSetup   func(mock sqlmock.Sqlmock)
		expected    models.Merch
		expectedErr error
	}{
		{
			name:    "Success case",
			merchID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "type", "cost"}).
					AddRow(1, "t-shirt", 80)
				mock.ExpectQuery("SELECT id, type, cost FROM merch WHERE id = \\$1").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expected: models.Merch{
				ID:   1,
				Type: "t-shirt",
				Cost: 80,
			},
			expectedErr: nil,
		},
		{
			name:    "Record not found",
			merchID: 999,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, type, cost  FROM merch WHERE id = \\$1").
					WithArgs(999).
					WillReturnError(sql.ErrNoRows)
			},
			expected:    models.Merch{},
			expectedErr: repositories.ErrRecordNotFound,
		},
		{
			name:    "Database error",
			merchID: 2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, type, cost FROM merch WHERE id = \\$1").
					WithArgs(2).
					WillReturnError(sql.ErrConnDone)
			},
			expected:    models.Merch{},
			expectedErr: sql.ErrConnDone,
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

			repo := repositories.NewMerchRepository(sqlxDB)
			result, err := repo.GetMerchByID(tt.merchID)

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

func TestMerchRepository_GetMerchByType(t *testing.T) {
	tests := []struct {
		name        string
		merchType   string
		mockSetup   func(mock sqlmock.Sqlmock)
		expected    models.Merch
		expectedErr error
	}{
		{
			name:      "Success case",
			merchType: "t-shirt",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "type", "cost"}).
					AddRow(1, "t-shirt", 80)
				mock.ExpectQuery("SELECT id, type, cost FROM merch WHERE type = \\$1").
					WithArgs("t-shirt").
					WillReturnRows(rows)
			},
			expected: models.Merch{
				ID:   1,
				Type: "t-shirt",
				Cost: 80,
			},
			expectedErr: nil,
		},
		{
			name:      "Record not found",
			merchType: "non-existent",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, type, cost FROM merch WHERE type = \\$1").
					WithArgs("non-existent").
					WillReturnError(sql.ErrNoRows)
			},
			expected:    models.Merch{},
			expectedErr: repositories.ErrRecordNotFound,
		},
		{
			name:      "Database error",
			merchType: "hoodie",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, type, cost FROM merch WHERE type = \\$1").
					WithArgs("hoodie").
					WillReturnError(sql.ErrConnDone)
			},
			expected:    models.Merch{},
			expectedErr: sql.ErrConnDone,
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

			repo := repositories.NewMerchRepository(sqlxDB)
			result, err := repo.GetMerchByType(tt.merchType)

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
