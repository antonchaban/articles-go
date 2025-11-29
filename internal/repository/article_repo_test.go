package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/antonchaban/articles-go/internal/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	cleanup := func() {
		sqlDB, _ := gormDB.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}

	return gormDB, mock, cleanup
}

func TestCreateArticleSuccessfully(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	logger := zap.NewNop()
	repo := NewPostgresRepo(db, logger)

	article := &entities.Article{
		Title:     "Test Article",
		CreatedAt: time.Now().UTC(),
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "articles"`)).
		WithArgs(article.Title, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), article)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateArticleWithDatabaseError(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	logger := zap.NewNop()
	repo := NewPostgresRepo(db, logger)

	article := &entities.Article{
		Title:     "Test Article",
		CreatedAt: time.Now().UTC(),
	}

	expectedError := errors.New("database connection failed")

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "articles"`)).
		WithArgs(article.Title, sqlmock.AnyArg()).
		WillReturnError(expectedError)
	mock.ExpectRollback()

	err := repo.Create(context.Background(), article)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByIDReturnsArticleSuccessfully(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	logger := zap.NewNop()
	repo := NewPostgresRepo(db, logger)

	expectedArticle := &entities.Article{
		ID:        1,
		Title:     "Test Article",
		CreatedAt: time.Now().UTC(),
	}

	rows := sqlmock.NewRows([]string{"id", "title", "created_at"}).
		AddRow(expectedArticle.ID, expectedArticle.Title, expectedArticle.CreatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "articles" WHERE "articles"."id" = $1 ORDER BY "articles"."id" LIMIT $2`)).
		WithArgs(1, 1).
		WillReturnRows(rows)

	article, err := repo.GetByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, article)
	assert.Equal(t, expectedArticle.ID, article.ID)
	assert.Equal(t, expectedArticle.Title, article.Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByIDReturnsNotFoundError(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	logger := zap.NewNop()
	repo := NewPostgresRepo(db, logger)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "articles" WHERE "articles"."id" = $1 ORDER BY "articles"."id" LIMIT $2`)).
		WithArgs(999, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	article, err := repo.GetByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, article)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByIDWithDatabaseError(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	logger := zap.NewNop()
	repo := NewPostgresRepo(db, logger)

	expectedError := errors.New("database query failed")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "articles" WHERE "articles"."id" = $1 ORDER BY "articles"."id" LIMIT $2`)).
		WithArgs(1, 1).
		WillReturnError(expectedError)

	article, err := repo.GetByID(context.Background(), 1)

	assert.Error(t, err)
	assert.Nil(t, article)
	assert.NoError(t, mock.ExpectationsWereMet())
}
