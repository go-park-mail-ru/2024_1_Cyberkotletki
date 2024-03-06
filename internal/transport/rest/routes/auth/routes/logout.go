package routes

import (
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/exceptions"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/services/auth"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/httputil"
	"net/http"
	"time"
)

// Logout
// @Tags Auth
// @Description Удаляет сессию
// @Param 	Cookie header string  true "session"     default(session=xxx)
// @Success     200
// @Failure		405	{object}	httputil.HTTPError	"Не авторизован"
// @Router /auth/logout [post]
func Logout(w http.ResponseWriter, r *http.Request) {
	if session, err := r.Cookie("session"); err != nil {
		httputil.NewError(w, 403, exc.Exception{
			When:  time.Now(),
			What:  "Не авторизован",
			Layer: exc.Service,
			Type:  exc.Forbidden,
		})
	} else {
		if err := auth.Logout(session.Value); err != nil {
			cookie := http.Cookie{
				Name:     "session",
				Value:    "",
				Expires:  time.Unix(0, 0),
				HttpOnly: true,
			}
			http.SetCookie(w, &cookie)
			httputil.NewError(w, 403, exc.Exception{
				When:  time.Now(),
				What:  "Не авторизован",
				Layer: exc.Service,
				Type:  exc.Forbidden,
			})
		} else {
			if _, err := w.Write([]byte("logout")); err != nil {
				httputil.NewError(w, 500, exc.Exception{
					When:  time.Now(),
					What:  "Внутренняя ошибка сервера",
					Layer: exc.Transport,
					Type:  exc.Untyped,
				})
			}
		}
	}
}
