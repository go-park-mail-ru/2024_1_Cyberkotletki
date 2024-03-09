package routes

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/services/auth"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/httputil"
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/exceptions"
	"net/http"
)

// IsAuth
// @Tags Auth
// @Description Проверяет, авторизован ли пользователь
// @Param 	Cookie header string  true "session"     default(session=xxx)
// @Success     200
// @Failure		403	{object}	httputil.HTTPError	"Не авторизован"
// @Failure		500	{object}	httputil.HTTPError	"Внутренняя ошибка сервера"
// @Router /auth/isAuth [get]
func IsAuth(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session")
	if err != nil {
		httputil.NewError(w, http.StatusForbidden, exc.New(exc.Service, exc.Forbidden, "не авторизован"))
		return
	}
	if _, err := auth.IsAuth(session.Value); err != nil {
		httputil.NewError(w, http.StatusForbidden, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
