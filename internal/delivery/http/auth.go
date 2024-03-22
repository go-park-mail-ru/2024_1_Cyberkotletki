package http

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/echoutil"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type AuthEndpoints struct {
	useCase usecase.Auth
}

func NewAuthEndpoints(useCase usecase.Auth) AuthEndpoints {
	return AuthEndpoints{useCase: useCase}
}

// Register
// @Tags Auth
// @Description Регистрация пользователя. Сразу же возвращает сессию в cookies
// @Accept json
// @Header 	200		{string}	Set-Cookie		"возвращает cookies с полученной сессией"
// @Param 	registerData	body	dto.Register	true	"Данные для регистрации"
// @Success     200
// @Failure		400	{object}	echo.HTTPError
// @Failure		409	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router /auth/register [post]
func (h *AuthEndpoints) Register(ctx echo.Context) error {
	registerData := new(dto.Register)
	if err := ctx.Bind(registerData); err != nil {
		return echoutil.NewError(ctx, http.StatusBadRequest, echoutil.ErrBadJSON)
	}
	key, err := h.useCase.Register(registerData)
	if err != nil {
		switch {
		case entity.Contains(err, entity.ErrAlreadyExists):
			return echoutil.NewError(ctx, http.StatusConflict, err)
		case entity.Contains(err, entity.ErrBadRequest):
			return echoutil.NewError(ctx, http.StatusBadRequest, err)
		default:
			return echoutil.NewError(ctx, http.StatusInternalServerError, err)
		}
	}
	cookie := http.Cookie{
		Name:     "session",
		Value:    key,
		Expires:  time.Now().Add(time.Duration(ctx.Get("params").(config.Config).Auth.SessionAliveTime) * time.Second),
		HttpOnly: true,
		Secure:   ctx.Get("params").(config.Config).HTTP.SecureCookies,
		Path:     "/",
	}
	ctx.SetCookie(&cookie)
	return nil
}

// Login
// @Tags Auth
// @Description Авторизация пользователя. При успешной авторизации отправляет куки с сессией. Если пользователь уже
// авторизован, то прежний cookies с сессией перезаписывается
// @Accept json
// @Header 	200		{string}	Set-Cookie		"возвращает cookies с полученной сессией"
// @Param 	loginData	body	dto.Login	true	"Данные для входа"
// @Success     200
// @Failure		400	{object}	echo.HTTPError
// @Failure		404	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router /auth/login [post]
func (h *AuthEndpoints) Login(ctx echo.Context) error {
	loginData := new(dto.Login)
	if err := ctx.Bind(loginData); err != nil {
		return echoutil.NewError(ctx, http.StatusBadRequest, echoutil.ErrBadJSON)
	}
	key, err := h.useCase.Login(loginData)
	if err != nil {
		switch {
		case entity.Contains(err, entity.ErrNotFound):
			return echoutil.NewError(ctx, http.StatusNotFound, err)
		case entity.Contains(err, entity.ErrBadRequest):
			return echoutil.NewError(ctx, http.StatusBadRequest, err)
		default:
			return echoutil.NewError(ctx, http.StatusInternalServerError, err)
		}
	}
	CookieSet(
		ctx,
		"session",
		key,
		time.Now().Add(time.Duration(ctx.Get("params").(config.Config).Auth.SessionAliveTime)*time.Second),
	)
	return nil
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
	session, err := ctx.Cookie("session")
	if err != nil {
		// нет cookies == не авторизован
		return echoutil.NewError(ctx, http.StatusUnauthorized, entity.NewClientError("не авторизован"))
	}
	isAuth, err := h.useCase.IsAuth(session.Value)
	if err != nil {
		return echoutil.NewError(ctx, http.StatusInternalServerError, err)
	}
	if !isAuth {
		return echoutil.NewError(ctx, http.StatusUnauthorized, entity.NewClientError("не авторизован"))
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
	_ = h.useCase.Logout(cookie.Value)
	CookieSet(ctx, "session", "", time.Unix(0, 0))
	return nil
}

func CookieSet(ctx echo.Context, key, value string, expires time.Time) {
	cookie := http.Cookie{
		Name:     key,
		Value:    value,
		Expires:  expires,
		HttpOnly: true,
		Secure:   ctx.Get("params").(config.Config).HTTP.SecureCookies,
		Path:     "/",
	}
	ctx.SetCookie(&cookie)
}
