package http

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPlaygroundEndpoints_Ping(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		ExpectedErr error
		ExpectedRes string
	}{
		{
			Name:        "Пинг",
			ExpectedErr: nil,
			ExpectedRes: "pong",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			e := echo.New()
			playgroundEndpoints := NewPlaygroundEndpoints()
			req := httptest.NewRequest(http.MethodGet, "/playground/ping", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := playgroundEndpoints.Ping(c)
			require.Equal(t, tc.ExpectedErr, err)
			require.Equal(t, tc.ExpectedRes, rec.Body.String())
		})
	}
}
