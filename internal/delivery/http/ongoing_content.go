package http

import (
	"errors"
	"golang.org/x/net/websocket"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	"github.com/labstack/echo/v4"
)

type OngoingContentEndpoints struct {
	contentUC usecase.Content
	authUC    usecase.Auth
}

func NewOngoingContentEndpoints(contentUC usecase.Content, authUC usecase.Auth) *OngoingContentEndpoints {
	return &OngoingContentEndpoints{
		contentUC: contentUC,
		authUC:    authUC,
	}
}

func (h *OngoingContentEndpoints) Configure(server *echo.Group) {
	server.GET("/nearest", h.GetNearestOngoings)
	server.GET("/:year/:month", h.GetOngoingContentByMonthAndYear)
	server.GET("/years", h.GetAllReleaseYears)
	server.GET("/:id/is_released", h.IsReleased)
	server.PUT("/:id/is_released", h.SetReleasedState)
	server.POST("/:id/subscribe", h.SubscribeOnContent)
	server.DELETE("/:id/subscribe", h.UnsubscribeFromContent)
	server.GET("/subscriptions", h.GetSubscribedContentIDs)
}

// GetNearestOngoings godoc
// @Summary Получить ближайшие релизы
// @Tags ongoing_content
// @Produce json
// @Success 200 {array} dto.PreviewContent
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /api/ongoing/nearest [get]
func (h *OngoingContentEndpoints) GetNearestOngoings(ctx echo.Context) error {
	ongoingContent, err := h.contentUC.GetNearestOngoings()
	switch {
	case errors.Is(err, usecase.ErrContentNotFound):
		return utils.NewError(ctx, http.StatusNotFound, "контент календаря релизов не найден", err)
	case err != nil:
		return utils.NewError(ctx, http.StatusInternalServerError, "ошибка при получении ближайших релизов", err)
	default:
		return utils.WriteJSON(ctx, ongoingContent)
	}
}

// GetOngoingContentByMonthAndYear godoc
// @Summary Получить релизы по месяцу и году
// @Tags ongoing_content
// @Produce json
// @Param month path int true "Месяц"
// @Param year path int true "Год"
// @Success 200 {array} dto.PreviewContent
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /api/ongoing/{year}/{month} [get]
func (h *OngoingContentEndpoints) GetOngoingContentByMonthAndYear(ctx echo.Context) error {
	month, err := strconv.ParseInt(ctx.Param("month"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, "невалидный месяц", err)
	}
	year, err := strconv.ParseInt(ctx.Param("year"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, "невалидный год", err)
	}

	ongoingContent, err := h.contentUC.GetOngoingContentByMonthAndYear(int(month), int(year))
	switch {
	case errors.Is(err, usecase.ErrContentNotFound):
		return utils.NewError(ctx, http.StatusNotFound, "контент календаря релизов не найден", err)
	case err != nil:
		return utils.NewError(ctx, http.StatusInternalServerError, "ошибка при получении релизов по месяцу и году", err)
	default:
		return utils.WriteJSON(ctx, ongoingContent)
	}
}

// GetAllReleaseYears godoc
// @Summary Получить все года релизов
// @Tags ongoing_content
// @Produce json
// @Success 200 {array} int
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /api/ongoing/years [get]
func (h *OngoingContentEndpoints) GetAllReleaseYears(ctx echo.Context) error {
	years, err := h.contentUC.GetAllOngoingsYears()
	switch {
	case errors.Is(err, usecase.ErrContentNotFound):
		return utils.NewError(ctx, http.StatusNotFound, "года релизов не найдены", err)
	case err != nil:
		return utils.NewError(ctx, http.StatusInternalServerError, "ошибка при получении годов релизов", err)
	default:
		return utils.WriteJSON(ctx, years)
	}
}

// IsReleased godoc
// @Summary Проверить, вышел ли контент. Использует WebSocket
// @Tags ongoing_content
// @Param id path int true "ID контента"
// @Success 101 {object} string "WebSocket connection is established"
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /api/ongoing/{id}/is_released [get]
func (h *OngoingContentEndpoints) IsReleased(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, "невалидный ID", err)
	}

	ws := websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		releasedCh := make(chan bool)
		errCh := make(chan error)

		go h.contentUC.IsOngoingContentReleased(int(id), releasedCh, errCh)

		for {
			select {
			case isReleased := <-releasedCh:
				if err := websocket.Message.Send(ws, strconv.FormatBool(isReleased)); err != nil {
					return
				}
				return
			case err := <-errCh:
				utils.WebsocketError(ctx, err)
				return
			}
		}
	})

	ws.ServeHTTP(ctx.Response(), ctx.Request())
	return nil
}

// SetReleasedState godoc
// @Summary Установить состояние релиза
// @Tags ongoing_content
// @Accept json
// @Produce json
// @Param id path int true "ID контента"
// @Param secret_key query string true "Секретный ключ"
// @Param is_released query bool true "Вышел ли контент"
// @Success 200 {object} string
// @Failure 400 {object} echo.HTTPError
// @Failure 403 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /api/ongoing/{id}/is_released [put]
// @Security _csrf
func (h *OngoingContentEndpoints) SetReleasedState(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, "невалидный ID", err)
	}

	secretKey := ctx.QueryParam("secret_key")
	isReleased, err := strconv.ParseBool(ctx.QueryParam("is_released"))
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, "невалидное значение is_released", err)
	}

	err = h.contentUC.SetReleasedState(secretKey, int(id), isReleased)
	switch {
	case errors.Is(err, usecase.ErrContentNotFound):
		return utils.NewError(ctx, http.StatusNotFound, "контент календаря релизов не найден", err)
	case errors.Is(err, usecase.ErrContentInvalidSecretKey):
		return utils.NewError(ctx, http.StatusForbidden, "неверный секретный ключ", err)
	case err != nil:
		return utils.NewError(ctx, http.StatusInternalServerError, "ошибка при установке состояния релиза", err)
	default:
		return ctx.NoContent(http.StatusOK)
	}
}

func subscribeResponse(ctx echo.Context, err error) error {
	switch {
	case errors.Is(err, usecase.ErrContentNotFound):
		return utils.NewError(ctx, http.StatusNotFound, "контент не найден", err)
	case errors.Is(err, usecase.ErrUserNotFound):
		return utils.NewError(ctx, http.StatusUnauthorized, "не авторизован", err)
	case err != nil:
		return utils.NewError(ctx, http.StatusInternalServerError, "ошибка при подписке на контент", err)
	default:
		return ctx.NoContent(http.StatusOK)
	}
}

// SubscribeOnContent godoc
// @Summary Подписаться на выход контента
// @Tags ongoing_content
// @Accept json
// @Produce json
// @Param id path int true "ID контента"
// @Success 200 {object} string
// @Failure 400 {object} echo.HTTPError
// @Failure 401 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /api/ongoing/{id}/subscribe [post]
// @Security _csrf
func (h *OngoingContentEndpoints) SubscribeOnContent(ctx echo.Context) error {
	contentID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, "невалидный ID", err)
	}
	userID, err := utils.GetUserIDFromSession(ctx, h.authUC)
	if err != nil {
		return utils.NewError(ctx, http.StatusUnauthorized, "Не авторизован", err)
	}

	err = h.contentUC.SubscribeOnContent(userID, int(contentID))
	return subscribeResponse(ctx, err)
}

// UnsubscribeFromContent godoc
// @Summary Отписаться от выхода контента
// @Tags ongoing_content
// @Accept json
// @Produce json
// @Param id path int true "ID контента"
// @Success 200 {object} string
// @Failure 400 {object} echo.HTTPError
// @Failure 401 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /api/ongoing/{id}/subscribe [delete]
// @Security _csrf
func (h *OngoingContentEndpoints) UnsubscribeFromContent(ctx echo.Context) error {
	contentID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, "невалидный ID", err)
	}
	userID, err := utils.GetUserIDFromSession(ctx, h.authUC)
	if err != nil {
		return utils.NewError(ctx, http.StatusUnauthorized, "Не авторизован", err)
	}

	err = h.contentUC.UnsubscribeFromContent(userID, int(contentID))
	return subscribeResponse(ctx, err)
}

// GetSubscribedContentIDs godoc
// @Summary Получить ID контентов, на которые подписан пользователь
// @Tags ongoing_content
// @Produce json
// @Success 200 {object} dto.SubscriptionsResponse
// @Failure 401 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /api/ongoing/subscriptions [get]
func (h *OngoingContentEndpoints) GetSubscribedContentIDs(ctx echo.Context) error {
	userID, err := utils.GetUserIDFromSession(ctx, h.authUC)
	if err != nil {
		return utils.NewError(ctx, http.StatusUnauthorized, "Не авторизован", err)
	}

	subscriptions, err := h.contentUC.GetSubscribedContentIDs(userID)
	if err != nil {
		return utils.NewError(ctx, http.StatusInternalServerError, "ошибка при получении подписок", err)
	}

	return utils.WriteJSON(ctx, subscriptions)
}
