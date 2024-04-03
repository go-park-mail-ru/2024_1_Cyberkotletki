package http

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type AuthEndpoints struct {
	authUC usecase.Auth
}

func NewAuthEndpoints(authUC usecase.Auth) AuthEndpoints {
	return AuthEndpoints{authUC: authUC}
}

func (h *AuthEndpoints) Configure(e *echo.Group) {
	e.GET("/isAuth", h.IsAuth)
	e.POST("/logout", h.Logout)
	e.POST("/logoutAll", h.LogoutAll)
}

// IsAuth
// @Tags Auth
// @Description Проверяет, авторизован ли пользователь
// @Param 	Cookie header string  true "session"     default(session=xxx)
// @Success     200
// @Failure		401	{object}	echo.HTTPError	"Не авторизован"
// @Failure		500	{object}	echo.HTTPError	"Внутренняя ошибка сервера"
// @Router /auth/isAuth [get]
func (h *AuthEndpoints) IsAuth(ctx echo.Context) error {
	if _, err := utils.GetUserIDFromSession(ctx, h.authUC); err != nil {
		return utils.NewError(ctx, http.StatusUnauthorized, err)
	}
	ctx.Response().WriteHeader(http.StatusOK)
	return nil
}

// Logout
// @Tags Auth
// @Description Удаляет сессию
// @Param 	Cookie header string  true "session"     default(session=xxx)
// @Success     200
// @Router /auth/logout [post]
func (h *AuthEndpoints) Logout(ctx echo.Context) error {
	cookie, err := ctx.Cookie("session")
	if err != nil {
		// сессия в куках не найдена, значит пользователь уже вышел
		ctx.Response().WriteHeader(http.StatusOK)
		return nil
	}
	// если сессии не было в базе сессий, то это не имеет значения - пользователь в любом случае вышел, поэтому
	// ошибку игнорируем
	// no-lint
	_ = h.authUC.Logout(cookie.Value)
	utils.SessionSet(ctx, "", time.Unix(0, 0))
	return nil
}

// LogoutAll
// @Tags Auth
// @Description Удаляет все сессии пользователя
// @Param 	Cookie header string  true "session"     default(session=xxx)
// @Success     200
// @Failure		500	{object}	echo.HTTPError	"Внутренняя ошибка сервера"
// @Router /auth/logoutAll [post]
func (h *AuthEndpoints) LogoutAll(ctx echo.Context) error {
	cookie, err := ctx.Cookie("session")
	if err != nil {
		// сессия в куках не найдена, значит считаем, что пользователь уже вышел
		ctx.Response().WriteHeader(http.StatusOK)
		return nil
	}
	userID, err := h.authUC.GetUserIDBySession(cookie.Value)
	if err != nil {
		return utils.NewError(ctx, http.StatusInternalServerError, err)
	}
	if err := h.authUC.LogoutAll(userID); err != nil {
		return utils.NewError(ctx, http.StatusInternalServerError, err)
	}
	utils.SessionSet(ctx, "", time.Unix(0, 0))
	return nil
}
