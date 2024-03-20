package http

import (
	"encoding/json"
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/DTO"
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
// @Param 	registerData	body	DTO.Register	true	"Данные для регистрации"
// @Success     200
// @Failure		400	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router /auth/register [post]
func (h *AuthEndpoints) Register(c echo.Context) error {
	decoder := json.NewDecoder(c.Request().Body)
	registerData := new(DTO.Register)
	if err := decoder.Decode(registerData); err != nil {
		return echoutil.NewError(c, http.StatusBadRequest, echoutil.BadJSON)
	}
	key, err := h.useCase.Register(*registerData)
	if err != nil {
		switch {
		// будем считать, что ErrAlreadyExists это тоже BadRequest
		case errors.Is(err, entity.ErrBadRequest) || errors.Is(err, entity.ErrAlreadyExists):
			return echoutil.NewError(c, http.StatusBadRequest, err)
		default:
			return echoutil.NewError(c, http.StatusInternalServerError, err)
		}
	}
	cookie := http.Cookie{
		Name:     "session",
		Value:    key,
		Expires:  time.Now().Add(time.Duration(c.Get("params").(config.Config).Auth.SessionAliveTime) * time.Second),
		HttpOnly: true,
		Secure:   c.Get("params").(config.Config).HTTP.SecureCookies,
		Path:     "/",
	}
	c.SetCookie(&cookie)
	return nil
}

// Login
// @Tags Auth
// @Description Авторизация пользователя. При успешной авторизации отправляет куки с сессией. Если пользователь уже авторизован, то прежний cookies с сессией перезаписывается
// @Accept json
// @Header 	200		{string}	Set-Cookie		"возвращает cookies с полученной сессией"
// @Param 	loginData	body	DTO.Login	true	"Данные для входа"
// @Success     200
// @Failure		400	{object}	echo.HTTPError
// @Failure		404	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router /auth/login [post]
func (h *AuthEndpoints) Login(c echo.Context) error {
	decoder := json.NewDecoder(c.Request().Body)
	loginData := new(DTO.Login)
	if err := decoder.Decode(loginData); err != nil {
		return echoutil.NewError(c, http.StatusBadRequest, echoutil.BadJSON)
	}
	key, err := h.useCase.Login(*loginData)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFound):
			return echoutil.NewError(c, http.StatusNotFound, err)
		case errors.Is(err, entity.ErrBadRequest):
			return echoutil.NewError(c, http.StatusBadRequest, err)
		default:
			return echoutil.NewError(c, http.StatusInternalServerError, err)
		}
	}
	cookie := http.Cookie{
		Name:     "session",
		Value:    key,
		Expires:  time.Now().Add(time.Duration(c.Get("params").(config.Config).Auth.SessionAliveTime) * time.Second),
		HttpOnly: true,
		Secure:   c.Get("params").(config.Config).HTTP.SecureCookies,
		Path:     "/",
	}
	c.SetCookie(&cookie)
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
func (h *AuthEndpoints) IsAuth(c echo.Context) error {
	session, err := c.Cookie("session")
	if err != nil {
		// нет cookies == не авторизован
		return echoutil.NewError(c, http.StatusUnauthorized, entity.NewClientError("не авторизован"))
	}
	isAuth, err := h.useCase.IsAuth(session.Value)
	if err != nil {
		return echoutil.NewError(c, http.StatusInternalServerError, err)
	}
	if !isAuth {
		return echoutil.NewError(c, http.StatusUnauthorized, entity.NewClientError("не авторизован"))
	}
	c.Response().WriteHeader(http.StatusOK)
	return nil
}

// Logout
// @Tags Auth
// @Description Удаляет сессию
// @Param 	Cookie header string  true "session"     default(session=xxx)
// @Success     200
// @Router /auth/logout [post]
func (h *AuthEndpoints) Logout(c echo.Context) error {
	cookie, err := c.Cookie("session")
	if err != nil {
		// сессия в куках не найдена, значит пользователь уже вышел
		c.Response().WriteHeader(http.StatusOK)
		return nil
	}
	// если сессии не было в базе сессий, то это не имеет значения - пользователь в любом случае вышел, поэтому
	// ошибку игнорируем
	err = h.useCase.Logout(cookie.Value)
	cookies := http.Cookie{
		Name:     "session",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   c.Get("params").(config.Config).HTTP.SecureCookies,
		Path:     "/",
	}
	c.SetCookie(&cookies)
	return nil
}
