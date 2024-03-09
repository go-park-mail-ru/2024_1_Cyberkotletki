package routes

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/services/auth"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/httputil"
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/exceptions"
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
		httputil.NewError(w, http.StatusBadRequest, httputil.BadJSON)
		return
	}
	key, err := auth.Login(*loginData)
	if err != nil {
		// forbidden (неверный пароль) тоже будет восприниматься как not found (пользователь не найден)
		if exc.Is(err, exc.NotFoundErr) || exc.Is(err, exc.ForbiddenErr) {
			httputil.NewError(w, http.StatusBadRequest, err)
		} else {
			httputil.NewError(w, http.StatusInternalServerError, err)
		}
		return
	}
	cookie := http.Cookie{
		Name:     "session",
		Value:    key,
		Expires:  time.Now().Add(r.Context().Value("params").(config.InitParams).SessionAliveTime),
		HttpOnly: true,
		Secure:   r.Context().Value("params").(config.InitParams).CookiesSecure,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
}
