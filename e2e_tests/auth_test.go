package e2e

import (
	"net/http"
	"testing"

	"github.com/hamillka/avitoTechWinter25/internal/handlers/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuth(t *testing.T) {
	client := http.Client{}

	var authResp dto.AuthResponseDto

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
			Name:           "register",
			Method:         http.MethodPost,
			route:          "/api/auth",
			body:           `{"username": "user", "password": "userPswd"}`,
			headers:        [][2]string{{"Content-Type", "application/json"}},
			expectedStatus: http.StatusOK,
			parseResp:      true,
			respBody:       &authResp,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			req, err := createRequest(tt.Method, tt.route, tt.body, tt.headers)
			require.NoError(t, err)

			sendRequest(t, &client, req, tt.expectedStatus, tt.parseResp, tt.respBody)

			assert.NotEmpty(t, authResp.Token)
		})
	}
}
