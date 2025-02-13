package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/hamillka/avitoTechWinter25/internal/handlers/dto"
	"github.com/hamillka/avitoTechWinter25/internal/repositories"
	"go.uber.org/zap"
)

type PurchaseHandler struct {
	service AvitoShopService
	logger  *zap.SugaredLogger
}

func NewPurchaseHandler(s AvitoShopService, logger *zap.SugaredLogger) *PurchaseHandler {
	return &PurchaseHandler{
		service: s,
		logger:  logger,
	}
}

func (ph *PurchaseHandler) BuyItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims := ctx.Value("props").(jwt.MapClaims)
	buyerName := claims["username"].(string)

	if buyerName == dto.AvitoShopName {
		ph.logger.Errorf("purchase handler: avito can't buy its own merch")
		w.WriteHeader(http.StatusBadRequest)
		errorDto := &dto.ErrorDto{
			Error: "Неверный запрос",
		}
		err := json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	itemName, ok := mux.Vars(r)["item"]
	if !ok {
		ph.logger.Errorf("purchase handler: %s", dto.ErrPathVarExtracting)
		w.WriteHeader(http.StatusBadRequest)
		errorDto := &dto.ErrorDto{
			Error: "Неверный запрос",
		}
		err := json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	err := ph.service.BuyItem(buyerName, itemName)
	if err != nil {
		ph.logger.Errorf("purchase handler: buy item service method: %s", err)
		var errorDto *dto.ErrorDto
		if errors.Is(err, repositories.ErrRecordNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			errorDto = &dto.ErrorDto{
				Error: "Неверный запрос",
			}
		} else if errors.Is(err, dto.ErrNotEnoughCoins) {
			w.WriteHeader(http.StatusBadRequest)
			errorDto = &dto.ErrorDto{
				Error: "Недостаточно монет для покупки",
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			errorDto = &dto.ErrorDto{
				Error: "Внутренняя ошибка сервера",
			}
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
