package routes

import (
	"encoding/json"
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/exceptions"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/services/auth"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/httputil"
	"net/http"
	"time"
)

// Login
// @Tags Auth
// @Description Авторизация пользователя. При успешной авторизации отправляет куки с сессией. Если пользователь уже авторизован, то прежний cookies с сессией перезаписывается
// @Accept json
// @Header 	200		{string}	Set-Cookie		"возвращает cookies с полученной сессией"
// @Param 	loginData	body	auth.LoginData	true	"Данные для входа"
// @Success     200
// @Failure		400	{object}	httputil.HTTPError
// @Failure		404	{object}	httputil.HTTPError
// @Failure		500	{object}	httputil.HTTPError
// @Router /auth/login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	loginData := new(auth.LoginData)
	if err := decoder.Decode(loginData); err != nil {
		httputil.NewError(w, 400, exc.Exception{
			When:  time.Now(),
			What:  "Невалидный JSON",
			Layer: exc.Transport,
			Type:  exc.Unprocessable,
		})
	} else if key, err := auth.Login(*loginData); err != nil {
		// forbidden тоже будет восприниматься как not found
		if err.Type == exc.NotFound || err.Type == exc.Forbidden {
			httputil.NewError(w, 404, *err)
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
			Secure:   true,
		}
		http.SetCookie(w, &cookie)
	}
}
