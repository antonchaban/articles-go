package main

import (
	"fmt"
	"log"

	"github.com/antonchaban/articles-go/internal/api"
	v1 "github.com/antonchaban/articles-go/internal/api/v1"
	"github.com/antonchaban/articles-go/internal/config"
	logger "github.com/antonchaban/articles-go/internal/log"
	"github.com/antonchaban/articles-go/internal/repository"
	"github.com/antonchaban/articles-go/internal/services"
	"github.com/antonchaban/articles-go/pkg/database"

	"go.uber.org/zap" // Import Zap
)

func main() {
	// Load Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Init Zap Logger
	l, cleanup, err := logger.NewLogger(cfg.AppEnv)
	if err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}
	defer cleanup()

	l.Info("application starting",
		zap.String("env", cfg.AppEnv),
		zap.String("port", cfg.HTTPPort))

	// DB Connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := database.NewPostgresConnection(dsn)
	if err != nil {
		l.Fatal("failed to connect to db", zap.Error(err))
	}

	// init repo, service, handler, and server
	repo := repository.NewPostgresRepo(db, l)
	svc := services.NewArticleService(repo, l)
	handler := v1.NewArticleHandler(svc, l)
	r := api.NewServer(cfg, handler)

	if err := r.Run(":" + cfg.HTTPPort); err != nil {
		l.Fatal("server failed to start", zap.Error(err))
	}
}
