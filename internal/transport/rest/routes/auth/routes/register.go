package routes

import (
	"encoding/json"
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/exceptions"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/services/auth"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/httputil"
	"net/http"
	"time"
)

// Register
// @Tags Auth
// @Description Регистрация пользователя. Сразу же возвращает сессию в cookies
// @Accept json
// @Header 	200		{string}	Set-Cookie		"возвращает cookies с полученной сессией"
// @Param 	registerData	body	auth.RegisterData	true	"Данные для регистрации"
// @Success     200
// @Failure		400	{object}	httputil.HTTPError
// @Failure		404	{object}	httputil.HTTPError
// @Failure		500	{object}	httputil.HTTPError
// @Router /auth/register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	registerData := new(auth.RegisterData)
	if err := decoder.Decode(registerData); err != nil {
		httputil.NewError(w, 400, exc.Exception{
			When:  time.Now(),
			What:  "Невалидный JSON",
			Layer: exc.Transport,
			Type:  exc.Unprocessable,
		})
	} else if key, err := auth.Register(*registerData); err != nil {
		if err.Layer != exc.Server {
			httputil.NewError(w, 400, *err)
		} else {
			httputil.NewError(w, 500, *err)
		}
	} else {
		// todo expiration в конфиг
		cookie := http.Cookie{
			Name:     "session",
			Value:    key,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
	}
}
