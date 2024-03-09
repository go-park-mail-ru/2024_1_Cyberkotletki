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
		httputil.NewError(w, http.StatusBadRequest, httputil.BadJSON)
		return
	}
	key, err := auth.Register(*registerData)
	if err != nil {
		if exc.Is(err, exc.ServerErr) {
			httputil.NewError(w, http.StatusInternalServerError, err)
		} else {
			httputil.NewError(w, http.StatusBadRequest, err)
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
