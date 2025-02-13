package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hamillka/avitoTechWinter25/internal/handlers/dto"
	"github.com/hamillka/avitoTechWinter25/internal/repositories"
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
		ch.logger.Errorf("coin handler: send coin service method: %s", err)
		var errorDto *dto.ErrorDto
		if errors.Is(err, repositories.ErrRecordNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			errorDto = &dto.ErrorDto{
				Error: "Неверный запрос",
			}
		} else if errors.Is(err, dto.ErrNotEnoughCoins) {
			w.WriteHeader(http.StatusBadRequest)
			errorDto = &dto.ErrorDto{
				Error: "Недостаточно монет для перевода",
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
