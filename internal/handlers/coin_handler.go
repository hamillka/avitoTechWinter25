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

// SendCoin godoc
//
//		@Summary		Отправка монет другому пользователю
//		@Description	Отправить пользователю монеты по его username
//		@ID				send-coin-to-user-by-username
//		@Tags			coin
//		@Accept			json
//		@Produce		json
//		@Param 			SendCoin	body	dto.SendCoinRequestDto		true	"Информация о пользователе и количество отправляемых монет"
//
//		@Success		200												"Успешный ответ"
//		@Failure		400	    {object}	dto.ErrorDto				"Неверный запрос or Недостаточно монет для перевода"
//		@Failure		401	    {object}	dto.ErrorDto				"Неавторизован"
//		@Failure		500	    {object}	dto.ErrorDto				"Внутренняя ошибка сервера"
//	    @Security		ApiKeyAuth
//		@Router			/api/sendCoin [post]
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
