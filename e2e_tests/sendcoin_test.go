package e2e

import (
	"net/http"
	"testing"

	"github.com/hamillka/avitoTechWinter25/internal/handlers/dto"
	"github.com/stretchr/testify/require"
)

func TestSendCoins(t *testing.T) {
	client := http.Client{}

	var senderAuthResp dto.AuthResponseDto
	var infoResp dto.InfoResponseDto

	tests := []struct {
		Name           string
		Method         string
		route          string
		body           string
		headers        [][2]string
		expectedStatus int
		parseResp      bool
		respBody       interface{}
	}{
		{
			Name:           "sender register",
			Method:         http.MethodPost,
			route:          "/api/auth",
			body:           `{"username": "sender", "password": "senderPswd"}`,
			headers:        [][2]string{{"Content-Type", "application/json"}},
			expectedStatus: http.StatusOK,
			parseResp:      true,
			respBody:       &senderAuthResp,
		},
		{
			Name:           "receiver register",
			Method:         http.MethodPost,
			route:          "/api/auth",
			body:           `{"username": "receiver", "password": "receiverPswd"}`,
			headers:        [][2]string{{"Content-Type", "application/json"}},
			expectedStatus: http.StatusOK,
			parseResp:      false,
			respBody:       nil,
		},
		{
			Name:           "send coins",
			Method:         http.MethodPost,
			route:          "/api/sendCoin",
			body:           "{\"toUser\": \"receiver\", \"amount\": 250}",
			headers:        [][2]string{{"Content-Type", "application/json"}},
			expectedStatus: http.StatusOK,
			parseResp:      false,
			respBody:       nil,
		},
		{
			Name:           "get user info",
			Method:         http.MethodGet,
			route:          "/api/info",
			headers:        [][2]string{{"Content-Type", "application/json"}},
			expectedStatus: http.StatusOK,
			parseResp:      true,
			respBody:       &infoResp,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			var token string
			if tt.Name == "send coins" || tt.Name == "get user info" {
				token = senderAuthResp.Token
			}

			req, err := createRequest(tt.Method, tt.route, tt.body, append(tt.headers, [2]string{"Authorization", "Bearer " + token}))
			require.NoError(t, err)

			sendRequest(t, &client, req, tt.expectedStatus, tt.parseResp, tt.respBody)
		})
	}
}
