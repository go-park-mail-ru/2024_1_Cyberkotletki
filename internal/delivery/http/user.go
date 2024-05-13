package http

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type UserEndpoints struct {
	userUC         usecase.User
	authUC         usecase.Auth
	staticUC       usecase.Static
	sessionManager *utils.SessionManager
}

func NewUserEndpoints(
	userUC usecase.User,
	authUC usecase.Auth,
	staticUC usecase.Static,
	sessionManager *utils.SessionManager,
) UserEndpoints {
	return UserEndpoints{userUC: userUC, authUC: authUC, staticUC: staticUC, sessionManager: sessionManager}
}

func (h *UserEndpoints) Configure(server *echo.Group) {
	server.POST("/register", h.Register)
	server.POST("/login", h.Login)
	server.PUT("/password", h.UpdatePassword)
	server.PUT("/avatar", h.UploadAvatar)
	server.PUT("/profile", h.UpdateInfo)
	server.GET("/profile", h.GetProfile)
	server.GET("/me", h.GetMyID)
}

// Register
// @Tags User
// @Description Регистрация пользователя. Сразу же возвращает сессию в cookies
// @Accept json
// @Header 	200		{string}	Set-Cookie		"возвращает cookies с полученной сессией"
// @Param 	registerData	body	dto.Register	true	"Данные для регистрации"
// @Success     200
// @Failure		400	{object}	echo.HTTPError
// @Failure		409	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router /api/user/register [post]
// @Security _csrf
func (h *UserEndpoints) Register(ctx echo.Context) error {
	registerData := new(dto.Register)
	if err := utils.ReadJSON(ctx, registerData); err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, "Невалидный JSON", nil)
	}
	userID, err := h.userUC.Register(registerData)
	var errUserIncorrectData usecase.UserIncorrectDataError
	switch {
	case errors.Is(err, usecase.ErrUserAlreadyExists):
		return utils.NewError(ctx, http.StatusConflict, "Пользователь с такой почтой уже существует", err)
	case errors.As(err, &errUserIncorrectData):
		return utils.NewError(ctx, http.StatusBadRequest, errUserIncorrectData.Err.Error(), err)
	case err != nil:
		return utils.NewError(ctx, http.StatusInternalServerError, "Внутренняя ошибка сервера", err)
	default:
		err = h.sessionManager.CreateSession(ctx, h.authUC, userID)
		if err != nil {
			return utils.NewError(ctx, http.StatusInternalServerError, "Внутренняя ошибка сервера", err)
		}
		return ctx.NoContent(http.StatusOK)
	}
}

// Login
// @Tags User
// @Description Авторизация пользователя. При успешной авторизации отправляет куки с сессией. Если пользователь уже
// авторизован, то прежний cookies с сессией перезаписывается
// @Accept json
// @Header 	200		{string}	Set-Cookie		"возвращает cookies с полученной сессией"
// @Param 	loginData	body	dto.Login	true	"Данные для входа"
// @Success     200
// @Failure		400	{object}	echo.HTTPError
// @Failure		403	{object}	echo.HTTPError
// @Failure		404	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router /api/user/login [post]
// @Security _csrf
func (h *UserEndpoints) Login(ctx echo.Context) error {
	loginData := new(dto.Login)
	if err := utils.ReadJSON(ctx, loginData); err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, "Невалидный JSON", nil)
	}
	userID, err := h.userUC.Login(loginData)
	var errUserIncorrectData usecase.UserIncorrectDataError
	switch {
	case errors.Is(err, usecase.ErrUserNotFound):
		return utils.NewError(ctx, http.StatusNotFound, "Пользователь не найден", err)
	case errors.As(err, &errUserIncorrectData):
		return utils.NewError(ctx, http.StatusForbidden, errUserIncorrectData.Err.Error(), err)
	case err != nil:
		return utils.NewError(ctx, http.StatusInternalServerError, "Внутренняя ошибка сервера", err)
	default:
		if err = h.sessionManager.CreateSession(ctx, h.authUC, userID); err != nil {
			return utils.NewError(ctx, http.StatusInternalServerError, "Внутренняя ошибка сервера", err)
		}
		return ctx.NoContent(http.StatusOK)
	}
}

// UpdatePassword
// @Tags User
// @Description Обновляет пароль пользователя. Необходимо быть авторизованным, при этом все сессии пользователя
// обнуляются
// @Accept json
// @Header 	200		{string}	Set-Cookie		"возвращает cookies с полученной сессией, если пароль успешно обновлён"
// @Param 	Cookie header string  true "session"     default(session=xxx)
// @Param 	loginData	body	dto.UpdatePassword	true	"Данные для обновления пароля"
// @Success     200
// @Failure		400	{object}	echo.HTTPError	"Неверный пароль или невалидный payload"
// @Failure		401	{object}	echo.HTTPError	"Не авторизован"
// @Failure		500	{object}	echo.HTTPError	"Внутренняя ошибка сервера"
// @Router /api/user/password [put]
// @Security _csrf
func (h *UserEndpoints) UpdatePassword(ctx echo.Context) error {
	userID, err := utils.GetUserIDFromSession(ctx, h.authUC)
	if err != nil {
		return utils.NewError(ctx, http.StatusUnauthorized, "Не авторизован", err)
	}
	updateData := new(dto.UpdatePassword)
	if err = utils.ReadJSON(ctx, updateData); err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, "Невалидный JSON", nil)
	}
	err = h.userUC.UpdatePassword(userID, updateData)
	switch {
	case errors.Is(err, usecase.ErrUserNotFound):
		return utils.NewError(ctx, http.StatusNotFound, "Пользователь не найден", err)
	case errors.As(err, &usecase.UserIncorrectDataError{}):
		return utils.NewError(ctx, http.StatusBadRequest, err.Error(), err)
	case err != nil:
		return utils.NewError(ctx, http.StatusInternalServerError, "Внутренняя ошибка сервера", err)
	}
	_, err = h.authUC.CreateSession(userID)
	if err != nil {
		return utils.NewError(ctx, http.StatusInternalServerError, "Внутренняя ошибка сервера", err)
	}
	return ctx.NoContent(http.StatusOK)
}

// UploadAvatar
// @Tags User
// @Description Позволяет загрузить аватарку пользователя. Необходимо быть авторизованным
// @Param 	Cookie header string  true "session"     default(session=xxx)
// @Param 	avatar formData file  true "файл с аватаркой"
// @Success     200
// @Failure		400	{object}	echo.HTTPError	"Невалидное изображение"
// @Failure		401	{object}	echo.HTTPError	"Не авторизован"
// @Failure		500	{object}	echo.HTTPError	"Внутренняя ошибка сервера"
// @Router /api/user/avatar [put]
// @Security _csrf
func (h *UserEndpoints) UploadAvatar(ctx echo.Context) error {
	userID, err := utils.GetUserIDFromSession(ctx, h.authUC)
	if err != nil {
		return utils.NewError(ctx, http.StatusUnauthorized, "Не авторизован", err)
	}
	fileHeader, err := ctx.FormFile("avatar")
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, "Файл не прикреплён", nil)
	}
	file, err := fileHeader.Open()
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, "Невалидный файл", nil)
	}
	err = h.userUC.UpdateAvatar(userID, file)
	var errUserIncorrectData usecase.UserIncorrectDataError
	switch {
	case errors.Is(err, usecase.ErrUserNotFound):
		return utils.NewError(ctx, http.StatusNotFound, "Пользователь не найден", err)
	case errors.As(err, &errUserIncorrectData):
		return utils.NewError(ctx, http.StatusBadRequest, errUserIncorrectData.Err.Error(), err)
	case err != nil:
		return utils.NewError(ctx, http.StatusInternalServerError, "Внутренняя ошибка сервера", err)
	}
	// ну не закрылся и не закрылся, чего бубнить-то
	// no-lint
	_ = file.Close()
	return ctx.NoContent(http.StatusOK)
}

// UpdateInfo
// @Tags User
// @Description Позволяет обновить следующие данные пользователя: почта, имя (никнейм). Необходимо быть авторизованным.
// @Accept json
// @Param 	Cookie header string  true "session"     default(session=xxx)
// @Param 	updateData	body	dto.UserUpdate	true	"Данные для обновления профиля"
// @Success     200
// @Failure		400	{object}	echo.HTTPError	"Невалидные данные для обновления профиля"
// @Failure		401	{object}	echo.HTTPError	"Не авторизован"
// @Failure		500	{object}	echo.HTTPError	"Внутренняя ошибка сервера"
// @Router /api/user/profile [put]
// @Security _csrf
func (h *UserEndpoints) UpdateInfo(ctx echo.Context) error {
	userID, err := utils.GetUserIDFromSession(ctx, h.authUC)
	if err != nil {
		return utils.NewError(ctx, http.StatusUnauthorized, "Не авторизован", err)
	}
	updateData := new(dto.UserUpdate)
	if err = utils.ReadJSON(ctx, updateData); err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, "Невалидный JSON", nil)
	}
	err = h.userUC.UpdateInfo(userID, updateData)
	var errUserIncorrectData usecase.UserIncorrectDataError
	switch {
	case errors.Is(err, usecase.ErrUserNotFound):
		return utils.NewError(ctx, http.StatusNotFound, "Пользователь не найден", err)
	case errors.As(err, &errUserIncorrectData):
		return utils.NewError(ctx, http.StatusBadRequest, err.Error(), err)
	case err != nil:
		return utils.NewError(ctx, http.StatusInternalServerError, "Внутренняя ошибка сервера", err)
	default:
		return ctx.NoContent(http.StatusOK)
	}
}

// GetProfile
// @Tags User
// @Description Возвращает профиль пользователя по id
// @Accept json
// @Produce json
// @Param 	id	query	int	true	"ID пользователя"
// @Success     200 {object}	dto.UserProfile
// @Failure		400	{object}	echo.HTTPError	"Неверный id"
// @Failure		404	{object}	echo.HTTPError	"Пользователь не найден"
// @Failure		500	{object}	echo.HTTPError	"Внутренняя ошибка сервера"
// @Router /api/user/profile [get]
func (h *UserEndpoints) GetProfile(ctx echo.Context) error {
	userID, err := strconv.ParseInt(ctx.QueryParam("id"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, "Неверный id", nil)
	}
	user, err := h.userUC.GetUser(int(userID))
	switch {
	case errors.Is(err, usecase.ErrUserNotFound):
		return utils.NewError(ctx, http.StatusNotFound, "Пользователь не найден", err)
	case err != nil:
		return utils.NewError(ctx, http.StatusInternalServerError, "Внутренняя ошибка сервера", err)
	default:
		return utils.WriteJSON(ctx, user)
	}
}

// GetMyID
// @Tags User
// @Description Возвращает id авторизованного пользователя
// @Accept json
// @Param 	Cookie header string  true "session"     default(session=xxx)
// @Success     200
// @Failure		401	{object}	echo.HTTPError	"Не авторизован"
// @Failure		500	{object}	echo.HTTPError	"Внутренняя ошибка сервера"
// @Router /api/user/me [get]
func (h *UserEndpoints) GetMyID(ctx echo.Context) error {
	userID, err := utils.GetUserIDFromSession(ctx, h.authUC)
	if err != nil {
		return utils.NewError(ctx, http.StatusUnauthorized, "Не авторизован", err)
	}
	response := dto.MyID{ID: userID}
	return utils.WriteJSON(ctx, response)
}
