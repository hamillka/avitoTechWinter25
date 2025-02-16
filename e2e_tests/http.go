package e2e

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRequest(method, route, body string, headers [][2]string) (*http.Request, error) {
	req, err := http.NewRequest(method, "http://localhost:8080"+route, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	for _, header := range headers {
		req.Header.Add(header[0], header[1])
	}

	return req, nil
}

func sendRequest(
	t *testing.T,
	client *http.Client,
	req *http.Request,
	expectedStatus int,
	parseResp bool,
	respBody interface{},
) {
	t.Helper()
	resp, err := client.Do(req)
	require.NoError(t, err)

	require.Equal(t, expectedStatus, resp.StatusCode)

	if parseResp {
		err = json.NewDecoder(resp.Body).Decode(respBody)
		require.NoError(t, err)
	}

	resp.Body.Close()
}
