package http

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type AuthEndpoints struct {
	authUC         usecase.Auth
	sessionManager *utils.SessionManager
}

func NewAuthEndpoints(authUC usecase.Auth, sessionManager *utils.SessionManager) AuthEndpoints {
	return AuthEndpoints{authUC: authUC, sessionManager: sessionManager}
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
// @Router /api/auth/isAuth [get]
func (h *AuthEndpoints) IsAuth(ctx echo.Context) error {
	_, err := utils.GetUserIDFromSession(ctx, h.authUC)
	if errors.Is(err, utils.ErrUnauthorized) {
		return utils.NewError(ctx, http.StatusUnauthorized, "Не авторизован", nil)
	}
	ctx.Response().WriteHeader(http.StatusOK)
	return nil
}

// Logout
// @Tags Auth
// @Description Удаляет сессию
// @Param 	Cookie header string  true "session"     default(session=xxx)
// @Success     200
// @Router /api/auth/logout [post]
// @Security _csrf
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
	h.sessionManager.SessionSet(ctx, "session", time.Unix(0, 0))
	return nil
}

// LogoutAll
// @Tags Auth
// @Description Удаляет все сессии пользователя
// @Param 	Cookie header string  true "session"     default(session=xxx)
// @Success     200
// @Failure		500	{object}	echo.HTTPError	"Внутренняя ошибка сервера"
// @Router /api/auth/logoutAll [post]
// @Security _csrf
func (h *AuthEndpoints) LogoutAll(ctx echo.Context) error {
	cookie, err := ctx.Cookie("session")
	if err != nil {
		// сессия в куках не найдена, значит считаем, что пользователь уже вышел
		ctx.Response().WriteHeader(http.StatusOK)
		return nil
	}
	userID, err := h.authUC.GetUserIDBySession(cookie.Value)
	if errors.Is(err, usecase.ErrSessionNotFound) {
		// сессия в базе не найдена, значит пользователь уже вышел
		ctx.Response().WriteHeader(http.StatusOK)
		return nil
	}
	if err != nil {
		return utils.NewError(ctx, http.StatusInternalServerError, "Внутренняя ошибка сервера", err)
	}
	if err = h.authUC.LogoutAll(userID); err != nil {
		return utils.NewError(ctx, http.StatusInternalServerError, "Внутренняя ошибка сервера", err)
	}
	h.sessionManager.SessionSet(ctx, "session", time.Unix(0, 0))
	return nil
}
