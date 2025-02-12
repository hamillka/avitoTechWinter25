package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hamillka/avitoTechWinter25/internal/handlers/dto"
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

func (ih *InfoHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims := ctx.Value("props").(jwt.MapClaims)
	username := claims["username"].(string)

	var infoResponseDto dto.InfoResponseDto
	info, err := ih.service.GetInfo(username)
	if err != nil {
		ih.logger.Errorf("info handler: get info service method: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		errorDto := &dto.ErrorDto{
			Error: "Внутренняя ошибка сервера",
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
