package http

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/echoutil"
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

func NewUserEndpoints(userUC usecase.User, authUC usecase.Auth) UserEndpoints {
	return UserEndpoints{userUC: userUC, authUC: authUC}
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
// @Router /user/register [post]
func (h *UserEndpoints) Register(ctx echo.Context) error {
	registerData := new(dto.Register)
	if err := ctx.Bind(registerData); err != nil {
		return echoutil.NewError(ctx, http.StatusBadRequest, echoutil.ErrBadJSON)
	}
	userId, err := h.userUC.Register(registerData)
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
	session, err := h.authUC.CreateSession(userId)
	if err != nil {
		return echoutil.NewError(ctx, http.StatusInternalServerError, err)
	}
	utils.SessionSet(
		ctx,
		session,
		time.Now().Add(time.Duration(ctx.Get("params").(config.Config).Auth.SessionAliveTime)*time.Second),
	)
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
// @Failure		403	{object}	echo.HTTPError
// @Failure		404	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router /user/login [post]
func (h *UserEndpoints) Login(ctx echo.Context) error {
	loginData := new(dto.Login)
	if err := ctx.Bind(loginData); err != nil {
		return echoutil.NewError(ctx, http.StatusBadRequest, echoutil.ErrBadJSON)
	}
	userId, err := h.userUC.Login(loginData)
	if err != nil {
		switch {
		case entity.Contains(err, entity.ErrNotFound):
			return echoutil.NewError(ctx, http.StatusNotFound, err)
		case entity.Contains(err, entity.ErrForbidden):
			return echoutil.NewError(ctx, http.StatusForbidden, err)
		default:
			return echoutil.NewError(ctx, http.StatusInternalServerError, err)
		}
	}
	session, err := h.authUC.CreateSession(userId)
	if err != nil {
		return echoutil.NewError(ctx, http.StatusInternalServerError, err)
	}
	utils.SessionSet(
		ctx,
		session,
		time.Now().Add(time.Duration(ctx.Get("params").(config.Config).Auth.SessionAliveTime)*time.Second),
	)
	return nil
}

// UpdatePassword
// @Tags Auth
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
	cookie, err := ctx.Cookie("session")
	if err != nil {
		return echoutil.NewError(ctx, http.StatusUnauthorized, entity.NewClientError("не авторизован"))
	}
	userID, err := h.authUC.GetUserIDBySession(cookie.Value)
	if err != nil {
		return echoutil.NewError(ctx, http.StatusInternalServerError, err)
	}
	updateData := new(dto.UpdatePassword)
	if err = ctx.Bind(updateData); err != nil {
		return echoutil.NewError(ctx, http.StatusBadRequest, echoutil.ErrBadJSON)
	}
	if err = h.userUC.UpdatePassword(userID, updateData); err != nil {
		switch {
		case entity.Contains(err, entity.ErrBadRequest):
			return echoutil.NewError(ctx, http.StatusBadRequest, err)
		default:
			return echoutil.NewError(ctx, http.StatusInternalServerError, err)
		}
	}
	session, err := h.authUC.CreateSession(userID)
	if err != nil {
		return echoutil.NewError(ctx, http.StatusInternalServerError, err)
	}
	utils.SessionSet(
		ctx,
		session,
		time.Now().Add(time.Duration(ctx.Get("params").(config.Config).Auth.SessionAliveTime)*time.Second),
	)
	return nil
}

func (h *UserEndpoints) UploadAvatar(ctx echo.Context) error {
	cookie, err := ctx.Cookie("session")
	if err != nil {
		return echoutil.NewError(ctx, http.StatusUnauthorized, entity.NewClientError("не авторизован"))
	}
	userID, err := h.authUC.GetUserIDBySession(cookie.Value)
	if err != nil {
		return echoutil.NewError(ctx, http.StatusInternalServerError, err)
	}
	fileHeader, err := ctx.FormFile("avatar")
	if err != nil {
		return echoutil.NewError(ctx, http.StatusBadRequest, entity.NewClientError("файл не прикреплен"))
	}
	file, err := fileHeader.Open()
	if err != nil {
		return echoutil.NewError(ctx, http.StatusInternalServerError, err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return echoutil.NewError(ctx, http.StatusInternalServerError, err)
	}
	if err = file.Close(); err != nil {
		return echoutil.NewError(ctx, http.StatusInternalServerError, err)
	}
	uploadID, err := h.staticUC.UploadAvatar(data)
	if err != nil {
		switch {
		case entity.Contains(err, entity.ErrBadRequest):
			return echoutil.NewError(ctx, http.StatusBadRequest, err)
		default:
			return echoutil.NewError(ctx, http.StatusInternalServerError, err)
		}
	}
	if err = h.userUC.UpdateAvatar(userID, uploadID); err != nil {
		return echoutil.NewError(ctx, http.StatusInternalServerError, err)
	}
	return nil
}

func (h *UserEndpoints) UpdateInfo(ctx echo.Context) error {
	cookie, err := ctx.Cookie("session")
	if err != nil {
		return echoutil.NewError(ctx, http.StatusUnauthorized, entity.NewClientError("не авторизован"))
	}
	userID, err := h.authUC.GetUserIDBySession(cookie.Value)
	if err != nil {
		return echoutil.NewError(ctx, http.StatusInternalServerError, err)
	}
	updateData := new(dto.UserUpdate)
	if err = ctx.Bind(updateData); err != nil {
		return echoutil.NewError(ctx, http.StatusBadRequest, echoutil.ErrBadJSON)
	}
	if err = h.userUC.UpdateInfo(userID, updateData); err != nil {
		switch {
		case entity.Contains(err, entity.ErrBadRequest):
			return echoutil.NewError(ctx, http.StatusBadRequest, err)
		default:
			return echoutil.NewError(ctx, http.StatusInternalServerError, err)
		}
	}
	return nil
}

func (h *UserEndpoints) GetProfile(ctx echo.Context) error {
	userID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return echoutil.NewError(ctx, http.StatusBadRequest, entity.NewClientError("неверный id"))
	}
	user, err := h.userUC.GetUser(int(userID))
	if err != nil {
		return echoutil.NewError(ctx, http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, user)
}
