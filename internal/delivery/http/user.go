package http

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"strconv"
	"time"
)

type UserEndpoints struct {
	userUC   usecase.User
	authUC   usecase.Auth
	staticUC usecase.Static
}

func NewUserEndpoints(userUC usecase.User, authUC usecase.Auth, staticUC usecase.Static) UserEndpoints {
	return UserEndpoints{userUC: userUC, authUC: authUC, staticUC: staticUC}
}

func (h *UserEndpoints) Configure(e *echo.Group) {
	e.POST("/register", h.Register)
	e.POST("/login", h.Login)
	e.PUT("/password", h.UpdatePassword)
	e.PUT("/avatar", h.UploadAvatar)
	e.PUT("/profile", h.UpdateInfo)
	e.GET("/profile", h.GetProfile)
	e.GET("/me", h.GetMyID)
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
// @Router /user/register [post]
func (h *UserEndpoints) Register(ctx echo.Context) error {
	registerData := new(dto.Register)
	if err := ctx.Bind(registerData); err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, utils.ErrBadJSON)
	}
	userId, err := h.userUC.Register(registerData)
	if err != nil {
		switch {
		case entity.Contains(err, entity.ErrAlreadyExists):
			return utils.NewError(ctx, http.StatusConflict, err)
		case entity.Contains(err, entity.ErrBadRequest):
			return utils.NewError(ctx, http.StatusBadRequest, err)
		default:
			return utils.NewError(ctx, http.StatusInternalServerError, err)
		}
	}
	session, err := h.authUC.CreateSession(userId)
	if err != nil {
		return utils.NewError(ctx, http.StatusInternalServerError, err)
	}
	utils.SessionSet(
		ctx,
		session,
		time.Now().Add(time.Duration(ctx.Get("params").(config.Config).Auth.SessionAliveTime)*time.Second),
	)
	return nil
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
// @Router /user/login [post]
func (h *UserEndpoints) Login(ctx echo.Context) error {
	loginData := new(dto.Login)
	if err := ctx.Bind(loginData); err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, utils.ErrBadJSON)
	}
	userId, err := h.userUC.Login(loginData)
	if err != nil {
		switch {
		case entity.Contains(err, entity.ErrNotFound):
			return utils.NewError(ctx, http.StatusNotFound, err)
		case entity.Contains(err, entity.ErrForbidden):
			return utils.NewError(ctx, http.StatusForbidden, err)
		default:
			return utils.NewError(ctx, http.StatusInternalServerError, err)
		}
	}
	session, err := h.authUC.CreateSession(userId)
	if err != nil {
		return utils.NewError(ctx, http.StatusInternalServerError, err)
	}
	utils.SessionSet(
		ctx,
		session,
		time.Now().Add(time.Duration(ctx.Get("params").(config.Config).Auth.SessionAliveTime)*time.Second),
	)
	return nil
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
// @Router /user/password [put]
func (h *UserEndpoints) UpdatePassword(ctx echo.Context) error {
	userID, err := utils.GetUserIDFromSession(ctx, h.authUC)
	if err != nil {
		return err
	}
	updateData := new(dto.UpdatePassword)
	if err = ctx.Bind(updateData); err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, utils.ErrBadJSON)
	}
	if err = h.userUC.UpdatePassword(userID, updateData); err != nil {
		switch {
		case entity.Contains(err, entity.ErrBadRequest):
			return utils.NewError(ctx, http.StatusBadRequest, err)
		default:
			return utils.NewError(ctx, http.StatusInternalServerError, err)
		}
	}
	session, err := h.authUC.CreateSession(userID)
	if err != nil {
		return utils.NewError(ctx, http.StatusInternalServerError, err)
	}
	utils.SessionSet(
		ctx,
		session,
		time.Now().Add(time.Duration(ctx.Get("params").(config.Config).Auth.SessionAliveTime)*time.Second),
	)
	return nil
}

// UploadAvatar
// @Tags User
// @Description Позволяет загрузить аватарку пользователя. Необходимо быть авторизованным
// @Accept json
// @Param 	Cookie header string  true "session"     default(session=xxx)
// @Param 	avatar formData file  true "файл с аватаркой"
// @Success     200
// @Failure		400	{object}	echo.HTTPError	"Невалидное изображение"
// @Failure		401	{object}	echo.HTTPError	"Не авторизован"
// @Failure		500	{object}	echo.HTTPError	"Внутренняя ошибка сервера"
// @Router /user/avatar [put]
func (h *UserEndpoints) UploadAvatar(ctx echo.Context) error {
	userID, err := utils.GetUserIDFromSession(ctx, h.authUC)
	if err != nil {
		return err
	}
	fileHeader, err := ctx.FormFile("avatar")
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, entity.NewClientError("файл не прикреплен"))
	}
	file, err := fileHeader.Open()
	if err != nil {
		return utils.NewError(ctx, http.StatusInternalServerError, err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return utils.NewError(ctx, http.StatusInternalServerError, err)
	}
	if err = file.Close(); err != nil {
		return utils.NewError(ctx, http.StatusInternalServerError, err)
	}
	uploadID, err := h.staticUC.UploadAvatar(data)
	if err != nil {
		switch {
		case entity.Contains(err, entity.ErrBadRequest):
			return utils.NewError(ctx, http.StatusBadRequest, err)
		default:
			return utils.NewError(ctx, http.StatusInternalServerError, err)
		}
	}
	if err = h.userUC.UpdateAvatar(userID, uploadID); err != nil {
		return utils.NewError(ctx, http.StatusInternalServerError, err)
	}
	return nil
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
// @Router /user/profile [put]
func (h *UserEndpoints) UpdateInfo(ctx echo.Context) error {
	userID, err := utils.GetUserIDFromSession(ctx, h.authUC)
	if err != nil {
		return err
	}
	updateData := new(dto.UserUpdate)
	if err = ctx.Bind(updateData); err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, utils.ErrBadJSON)
	}
	if err = h.userUC.UpdateInfo(userID, updateData); err != nil {
		switch {
		case entity.Contains(err, entity.ErrBadRequest):
			return utils.NewError(ctx, http.StatusBadRequest, err)
		default:
			return utils.NewError(ctx, http.StatusInternalServerError, err)
		}
	}
	return nil
}

// GetProfile
// @Tags User
// @Description Возвращает профиль пользователя по id
// @Accept json
// @Param 	id	query	int	true	"ID пользователя"
// @Success     200
// @Failure		400	{object}	echo.HTTPError	"Неверный id"
// @Failure		404	{object}	echo.HTTPError	"Пользователь не найден"
// @Failure		500	{object}	echo.HTTPError	"Внутренняя ошибка сервера"
// @Router /user/profile [get]
func (h *UserEndpoints) GetProfile(ctx echo.Context) error {
	userID, err := strconv.ParseInt(ctx.QueryParam("id"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, entity.NewClientError("неверный id"))
	}
	user, err := h.userUC.GetUser(int(userID))
	if err != nil {
		switch {
		case entity.Contains(err, entity.ErrNotFound):
			return utils.NewError(ctx, http.StatusNotFound, err)
		default:
			return utils.NewError(ctx, http.StatusInternalServerError, err)
		}
	}
	return ctx.JSON(http.StatusOK, user)
}

// GetMyID
// @Tags User
// @Description Возвращает id авторизованного пользователя
// @Accept json
// @Param 	Cookie header string  true "session"     default(session=xxx)
// @Success     200
// @Failure		401	{object}	echo.HTTPError	"Не авторизован"
// @Failure		500	{object}	echo.HTTPError	"Внутренняя ошибка сервера"
// @Router /user/me [get]
func (h *UserEndpoints) GetMyID(ctx echo.Context) error {
	userID, err := utils.GetUserIDFromSession(ctx, h.authUC)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, map[string]int{"id": userID})
}
