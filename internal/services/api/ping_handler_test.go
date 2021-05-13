package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	pkg_testutil "tiniapp-backend-oauth-sample/pkg/testutil"
)

func TestIntegrationService_Ping(t *testing.T) {
	r := svc.GetRouter()
	require.NotNil(t, r)

	ctx := context.Background()
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/api/ping",
		nil,
	)
	require.NoError(t, err)
	require.NotEmpty(t, req)

	pkg_testutil.TestGinHTTPResponse(
		r,
		req,
		func(res *httptest.ResponseRecorder) {
			require.Equal(t, 200, res.Code)

			var body struct {
				Build     string `json:"build"`
				Timestamp int64  `json:"timestamp"`
			}

			err = json.Unmarshal(res.Body.Bytes(), &body)
			require.NoError(t, err)
			require.NotEmpty(t, body)
		},
	)
}
