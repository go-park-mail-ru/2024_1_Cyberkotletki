package http

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPlaygroundEndpoints_Ping(t *testing.T) {
	h := NewPlaygroundEndpoints()

	e := echo.New()

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
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodGet, "/playground/ping", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := h.Ping(c)

			if err != nil && tc.ExpectedErr == nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err == nil && tc.ExpectedErr != nil {
				t.Fatalf("expected error: %v, got nil", tc.ExpectedErr)
			}
			if err != nil && tc.ExpectedErr != nil && err.Error() != tc.ExpectedErr.Error() {
				t.Fatalf("expected error: %v, got error: %v", tc.ExpectedErr, err)
			}
			if rec.Code != http.StatusOK {
				t.Fatalf("expected status 200, got %v", rec.Code)
			}
			if rec.Body.String() != tc.ExpectedRes {
				t.Fatalf("expected body %q, got %q", tc.ExpectedRes, rec.Body.String())
			}
		})
	}
}
