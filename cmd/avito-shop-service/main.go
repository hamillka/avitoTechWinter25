package main

import (
	"net/http"

	cfg "github.com/hamillka/avitoTechWinter25/internal/config"
	"github.com/hamillka/avitoTechWinter25/internal/db"
	"github.com/hamillka/avitoTechWinter25/internal/handlers"
	"github.com/hamillka/avitoTechWinter25/internal/logger"
	"github.com/hamillka/avitoTechWinter25/internal/repositories"
	"github.com/hamillka/avitoTechWinter25/internal/services"
)

func main() {
	config, err := cfg.New()
	logger := logger.CreateLogger(config.Log)

	defer func() {
		err := logger.Sync()
		if err != nil {
			logger.Errorf("Error while syncing logger: %v", err)
		}
	}()

	if err != nil {
		logger.Errorf("Something went wrong with config: %v", err)
	}

	db, err := db.CreateConnection(&config.DB)

	defer func() {
		err := db.Close()
		if err != nil {
			logger.Errorf("Error while closing connection to db: %v", err)
		}
	}()

	if err != nil {
		logger.Fatalf("Error while connecting to database: %v", err)
	}

	ur := repositories.NewUserRepository(db)
	ir := repositories.NewInventoryRepository(db)
	tr := repositories.NewTransactionsRepository(db)
	mr := repositories.NewMerchRepository(db)

	as := services.NewAvitoShopService(ur, ir, mr, tr)

	r := handlers.Router(as, logger)

	port := config.Port
	logger.Info("Server is started on port ", port)
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		logger.Fatalf("Error while starting server: %v", err)
	}
}
