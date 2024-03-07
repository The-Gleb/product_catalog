package v1

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequest(t *testing.T, sToken string, ts *httptest.Server, method, path string, body []byte) (*http.Response, string) {
	req, err := http.NewRequestWithContext(context.Background(), method, ts.URL+path, bytes.NewReader(body))
	require.NoError(t, err)

	if sToken != "" {
		c := http.Cookie{
			Name:  "sessionToken",
			Value: sToken,
		}
		req.AddCookie(&c)
	}

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}
