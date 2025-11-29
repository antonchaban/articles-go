package database

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestNewPostgresConnectionSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT current_database()`)).
		WillReturnRows(sqlmock.NewRows([]string{"current_database"}).AddRow("testdb"))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "articles" LIMIT 1`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "created_at"}))

	mock.ExpectExec(regexp.QuoteMeta(`CREATE TABLE "articles"`)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	dialector := postgres.New(postgres.Config{
		Conn: db,
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})

	assert.NoError(t, err)
	assert.NotNil(t, gormDB)
}

func TestNewPostgresConnectionWithInvalidDSN(t *testing.T) {
	invalidDSN := "invalid connection string"

	db, err := NewPostgresConnection(invalidDSN)

	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestNewPostgresConnectionWithAutoMigrateError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT current_database()`)).
		WillReturnRows(sqlmock.NewRows([]string{"current_database"}).AddRow("testdb"))

	expectedError := errors.New("migration failed")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "articles" LIMIT 1`)).
		WillReturnError(expectedError)

	dialector := postgres.New(postgres.Config{
		Conn: db,
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	err = gormDB.AutoMigrate(&struct {
		ID    uint
		Title string
	}{})

	assert.Error(t, err)
}
