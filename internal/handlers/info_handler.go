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

type InfoHandler struct {
	service AvitoShopService
	logger  *zap.SugaredLogger
}

func NewInfoHandler(s AvitoShopService, logger *zap.SugaredLogger) *InfoHandler {
	return &InfoHandler{
		service: s,
		logger:  logger,
	}
}

// GetInfo godoc
//
//		@Summary		Получение информации о пользователе
//		@Description	Получить информацию об остатке монет, инвентаре и истории переводов конкретного пользователя
//		@ID				get-user-info
//		@Tags			info
//		@Accept			json
//		@Produce		json
//
//		@Success		200	    {object} 	dto.InfoResponseDto			"Успешный ответ"
//		@Failure		400	    {object}	dto.ErrorDto				"Неверный запрос"
//		@Failure		401	    {object}	dto.ErrorDto				"Неавторизован"
//		@Failure		500	    {object}	dto.ErrorDto				"Внутренняя ошибка сервера"
//	    @Security		ApiKeyAuth
//		@Router			/api/info [get]
func (ih *InfoHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims := ctx.Value("props").(jwt.MapClaims)
	username := claims["username"].(string)

	var infoResponseDto dto.InfoResponseDto
	info, err := ih.service.GetInfo(username)
	if err != nil {
		errorDto := &dto.ErrorDto{}
		ih.logger.Errorf("info handler: get info service method: %s", err)
		if errors.Is(err, repositories.ErrRecordNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			errorDto = &dto.ErrorDto{
				Error: "Неверный запрос",
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

	infoResponseDto = *dto.ConvertBLInfoToDto(info)

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(infoResponseDto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
