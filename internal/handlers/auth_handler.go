package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hamillka/avitoTechWinter25/internal/handlers/dto"
	"github.com/hamillka/avitoTechWinter25/internal/handlers/middlewares"
	"go.uber.org/zap"
)

type AuthHandler struct {
	service AvitoShopService
	logger  *zap.SugaredLogger
}

func NewAuthHandler(s AvitoShopService, logger *zap.SugaredLogger) *AuthHandler {
	return &AuthHandler{
		service: s,
		logger:  logger,
	}
}

func (ah *AuthHandler) Auth(w http.ResponseWriter, r *http.Request) {
	var authReqDto dto.AuthRequestDto

	w.Header().Add("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&authReqDto)
	if err != nil {
		ah.logger.Errorf("user handler: json decode %s", err)
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

	user, err := ah.service.Login(authReqDto.Username, authReqDto.Password)
	if err != nil {
		ah.logger.Errorf("auth handler: login service method: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		errorDto := &dto.ErrorDto{
			Error: "Неавторизован",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	t, err := createToken(user.Username)
	if err != nil {
		ah.logger.Errorf("auth handler: create token method: %s", err)
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

	authRespDto := dto.AuthResponseDto{
		Token: t,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(authRespDto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func createToken(username string) (string, error) {
	payload := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 8760).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString(middlewares.Secret)
	if err != nil {
		return "", err
	}

	return t, nil
}
