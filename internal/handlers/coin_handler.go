package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hamillka/avitoTechWinter25/internal/handlers/dto"
	"go.uber.org/zap"
)

type CoinHandler struct {
	service AvitoShopService
	logger  *zap.SugaredLogger
}

func NewCoinHandler(s AvitoShopService, logger *zap.SugaredLogger) *CoinHandler {
	return &CoinHandler{
		service: s,
		logger:  logger,
	}
}

func (ch *CoinHandler) SendCoin(w http.ResponseWriter, r *http.Request) {
	var sendCoinReqDto dto.SendCoinRequestDto

	w.Header().Add("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&sendCoinReqDto)
	if err != nil {
		ch.logger.Errorf("coin handler: json decode %s", err)
		w.WriteHeader(http.StatusBadRequest)
		errorDto := &dto.ErrorDto{
			Error: "Неверный запрос",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	ctx := r.Context()
	claims := ctx.Value("props").(jwt.MapClaims)
	fromUser := claims["username"].(string)

	if fromUser == sendCoinReqDto.ToUser {
		ch.logger.Errorf("coin handler: same usernames")
		w.WriteHeader(http.StatusBadRequest)
		errorDto := &dto.ErrorDto{
			Error: "Неверный запрос",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	err = ch.service.SendCoin(fromUser, sendCoinReqDto.ToUser, sendCoinReqDto.Amount)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
}
