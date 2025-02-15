package e2e

import (
	"net/http"
	"testing"

	"github.com/hamillka/avitoTechWinter25/internal/handlers/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetInfo(t *testing.T) {
	client := http.Client{}

	var authResp dto.AuthResponseDto
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
			Name:           "user register",
			Method:         http.MethodPost,
			route:          "/api/auth",
			body:           `{"username": "testUser", "password": "testPswd"}`,
			headers:        [][2]string{{"Content-Type", "application/json"}},
			expectedStatus: http.StatusOK,
			parseResp:      true,
			respBody:       &authResp,
		},
		{
			Name:           "buy item",
			Method:         http.MethodGet,
			route:          "/api/buy/wallet",
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
			if tt.Name == "buy item" || tt.Name == "get user info" {
				token = authResp.Token
			}

			req, err := createRequest(tt.Method, tt.route, tt.body, append(tt.headers, [2]string{"Authorization", "Bearer " + token}))
			require.NoError(t, err)

			sendRequest(t, &client, req, tt.expectedStatus, tt.parseResp, tt.respBody)

			if tt.Name == "get user info" {
				assert.Less(t, infoResp.Coins, int64(1000))
				assert.Equal(t, "wallet", infoResp.Inventory[0].Type)
			}
		})
	}
}
