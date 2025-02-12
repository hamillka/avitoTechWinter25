package handlers

import (
	"github.com/gorilla/mux"
	"github.com/hamillka/avitoTechWinter25/internal/handlers/middlewares"
	"go.uber.org/zap"
)

func Router(s AvitoShopService, logger *zap.SugaredLogger) *mux.Router {
	router := mux.NewRouter()
	auth := router.PathPrefix("").Subrouter()
	api := router.PathPrefix("/api").Subrouter()

	ah := NewAuthHandler(s, logger)
	ch := NewCoinHandler(s, logger)
	ih := NewInfoHandler(s, logger)
	ph := NewPurchaseHandler(s, logger)

	auth.HandleFunc("/api/auth", ah.Auth).Methods("POST")

	api.HandleFunc("/info", ih.GetInfo).Methods("GET")
	api.HandleFunc("/sendCoin", ch.SendCoin).Methods("POST")
	api.HandleFunc("/buy/{item}", ph.BuyItem).Methods("GET")

	api.Use(middlewares.AuthMiddleware)

	return router
}
