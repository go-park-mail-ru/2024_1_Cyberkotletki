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

type FavouriteEndpoints struct {
	favouriteUC usecase.Favourite
	authUC      usecase.Auth
}

func NewFavouriteEndpoints(favouriteUC usecase.Favourite, authUC usecase.Auth) FavouriteEndpoints {
	return FavouriteEndpoints{favouriteUC: favouriteUC, authUC: authUC}
}

func (h *FavouriteEndpoints) Configure(server *echo.Group) {
	server.PUT("", h.CreateFavourite)
	server.DELETE("/:id", h.DeleteFavourite)
	server.GET("/:id", h.GetFavouritesByUser)
	server.GET("/my", h.GetMyFavourites)
	server.GET("/status/:id", h.GetStatus)
}

// CreateFavourite
// @Tags Favourite
// @Description Добавление в избранное. Если уже в избранном, то ошибка не возвращается (идемпотентный метод).
// Возможные категории для добавления в избранное: favourite, watching, watched, planned, rewatching, abandoned
// @Accept json
// @Param 	payload	body	dto.CreateFavouriteRequest	true	"Данные для добавления в избранное"
// @Success     200
// @Failure		400	{object}	echo.HTTPError
// @Failure		404	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router /api/favourite [put]
// @Security _csrf
func (h *FavouriteEndpoints) CreateFavourite(ctx echo.Context) error {
	userID, err := utils.GetUserIDFromSession(ctx, h.authUC)
	if err != nil {
		return utils.NewError(ctx, http.StatusUnauthorized, "Не авторизован", err)
	}
	favouriteData := new(dto.CreateFavouriteRequest)
	if err := ctx.Bind(favouriteData); err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, "Невалидный JSON", nil)
	}
	err = h.favouriteUC.CreateFavourite(userID, favouriteData.ContentID, favouriteData.Category)
	switch {
	case errors.Is(err, usecase.ErrFavouriteContentNotFound):
		return utils.NewError(ctx, http.StatusNotFound, "Контент не найден", err)
	case err != nil:
		return utils.NewError(ctx, http.StatusInternalServerError, "Внутренняя ошибка сервера", err)
	default:
		return ctx.NoContent(http.StatusOK)
	}
}

// DeleteFavourite
// @Tags Favourite
// @Description Удаление из избранного.
// @Param 	id	path	int	true	"Идентификатор контента"
// @Success     200
// @Failure		400	{object}	echo.HTTPError
// @Failure		401	{object}	echo.HTTPError
// @Failure		404	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router /api/favourite/{id} [delete]
// @Security _csrf
func (h *FavouriteEndpoints) DeleteFavourite(ctx echo.Context) error {
	userID, err := utils.GetUserIDFromSession(ctx, h.authUC)
	if err != nil {
		return utils.NewError(ctx, http.StatusUnauthorized, "Не авторизован", err)
	}
	contentID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, "Невалидный ID", err)
	}
	err = h.favouriteUC.DeleteFavourite(userID, int(contentID))
	switch {
	case errors.Is(err, usecase.ErrFavouriteNotFound):
		return utils.NewError(ctx, http.StatusNotFound, "Избранное не найдено", err)
	case err != nil:
		return utils.NewError(ctx, http.StatusInternalServerError, "Внутренняя ошибка сервера", err)
	default:
		return ctx.NoContent(http.StatusOK)
	}
}

// GetFavouritesByUser
// @Tags Favourite
// @Description Получение избранного пользователя
// @Param 	id	path	int	true	"Идентификатор пользователя"
// @Success     200 {object}    dto.FavouritesResponse
// @Failure		400	{object}	echo.HTTPError
// @Failure		404	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router /api/favourite/{id} [get]
func (h *FavouriteEndpoints) GetFavouritesByUser(ctx echo.Context) error {
	userID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, "Невалидный ID", err)
	}
	favourites, err := h.favouriteUC.GetFavourites(int(userID))
	switch {
	case errors.Is(err, usecase.ErrFavouriteUserNotFound):
		return utils.NewError(ctx, http.StatusNotFound, "Пользователь не найден", err)
	case err != nil:
		return utils.NewError(ctx, http.StatusInternalServerError, "Внутренняя ошибка сервера", err)
	default:
		return utils.WriteJSON(ctx, favourites)
	}
}

// GetMyFavourites
// @Tags Favourite
// @Description Получение избранного пользователя
// @Success     200 {object}    dto.FavouritesResponse
// @Failure		401	{object}	echo.HTTPError
// @Failure		404	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router /api/favourite/my [get]
func (h *FavouriteEndpoints) GetMyFavourites(ctx echo.Context) error {
	userID, err := utils.GetUserIDFromSession(ctx, h.authUC)
	if err != nil {
		return utils.NewError(ctx, http.StatusUnauthorized, "Не авторизован", err)
	}
	favourites, err := h.favouriteUC.GetFavourites(userID)
	switch {
	case errors.Is(err, usecase.ErrFavouriteUserNotFound):
		return utils.NewError(ctx, http.StatusNotFound, "Пользователь не найден", err)
	case err != nil:
		return utils.NewError(ctx, http.StatusInternalServerError, "Внутренняя ошибка сервера", err)
	default:
		return utils.WriteJSON(ctx, favourites)
	}
}

// GetStatus
// @Tags Favourite
// @Description Получение статуса контента в избранном
// @Param 	id	path	int	true	"Идентификатор контента"
// @Success     200 {object}    dto.FavouriteStatusResponse
// @Failure		400	{object}	echo.HTTPError
// @Failure		401	{object}	echo.HTTPError
// @Failure		404	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router /api/favourite/status/{id} [get]
func (h *FavouriteEndpoints) GetStatus(ctx echo.Context) error {
	userID, err := utils.GetUserIDFromSession(ctx, h.authUC)
	if err != nil {
		return utils.NewError(ctx, http.StatusUnauthorized, "Не авторизован", err)
	}
	contentID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, "Невалидный ID", err)
	}
	status, err := h.favouriteUC.GetStatus(userID, int(contentID))
	switch {
	case errors.Is(err, usecase.ErrFavouriteNotFound):
		return utils.NewError(ctx, http.StatusNotFound, "Не добавлено в избранное", err)
	case err != nil:
		return utils.NewError(ctx, http.StatusInternalServerError, "Внутренняя ошибка сервера", err)
	default:
		return utils.WriteJSON(ctx, status)
	}
}
