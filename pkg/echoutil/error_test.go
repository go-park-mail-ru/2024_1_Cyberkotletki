package echoutil

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Error(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c := e.NewContext(req, nil)

	testCases := []struct {
		Name     string
		Status   int
		Err      error
		Expected *echo.HTTPError
	}{
		{
			Name:   "Client",
			Status: http.StatusBadRequest,
			Err:    entity.NewClientError("client error"),
			Expected: &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "client error",
			},
		},
		{
			Name:   "Internal server",
			Status: http.StatusInternalServerError,
			Err:    errors.New("internal error"),
			Expected: &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := NewError(c, tc.Status, tc.Err)
			if err.Code != tc.Expected.Code {
				t.Errorf("expected status code %v, got %v", tc.Expected.Code, err.Code)
			}
			if err.Message != tc.Expected.Message {
				t.Errorf("expected message %q, got %q", tc.Expected.Message, err.Message)
			}
		})
	}
}
