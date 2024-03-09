package routes

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/services/auth"
	"net/http"
	"time"
)

// Logout
// @Tags Auth
// @Description Удаляет сессию
// @Param 	Cookie header string  true "session"     default(session=xxx)
// @Success     200
// @Router /auth/logout [post]
func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		// сессия в куках не найдена, значит пользователь уже вышел
		w.WriteHeader(http.StatusOK)
		return
	}
	// если сессии не было в базе сессий, то это не имеет значения - пользователь в любом случае вышел
	_ = auth.Logout(cookie.Value)
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   r.Context().Value("params").(config.InitParams).CookiesSecure,
		Path:     "/",
	})
}
