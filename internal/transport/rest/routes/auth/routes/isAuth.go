package routes

import (
	"fmt"
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/exceptions"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/services/auth"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/httputil"
	"net/http"
	"time"
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
	fmt.Println(r.Cookie("session"))
	if session, err := r.Cookie("session"); err != nil {
		httputil.NewError(w, 403, exc.Exception{
			When:  time.Now(),
			What:  "Не авторизован",
			Layer: exc.Service,
			Type:  exc.Forbidden,
		})
	} else {
		if _, err := auth.IsAuth(session.Value); err != nil {
			httputil.NewError(w, 403, *err)
		} else {
			if _, err := w.Write([]byte("authenticated")); err != nil {
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

// 8d2528b8-7daf-42c4-b0b3-d368ce561b8d
