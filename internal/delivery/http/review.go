package http

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type ReviewEndpoints struct {
	reviewUC usecase.Review
	authUC   usecase.Auth
}

func NewReviewEndpoints(reviewUC usecase.Review, authUC usecase.Auth) ReviewEndpoints {
	return ReviewEndpoints{reviewUC: reviewUC, authUC: authUC}
}

func (h *ReviewEndpoints) Configure(server *echo.Group) {
	server.GET("/:id", h.GetReview)
	server.GET("/myReview", h.GetMyContentReview)
	server.POST("", h.CreateReview)
	server.PUT("", h.UpdateReview)
	server.DELETE("/:id", h.DeleteReview)
	server.GET("/recent", h.GetRecentReviews)
	server.GET("/user/:id/recent", h.GetUserLatestReviews)
	server.GET("/user/:id/:page", h.GetUserReviews)
	server.GET("/content/:id/:page", h.GetContentReviews)
	server.PUT("/:id/like", h.LikeReview)
	server.PUT("/:id/dislike", h.DislikeReview)
	server.DELETE("/:id/like", h.UnlikeReview)
}

// GetReview
// @Summary Получить рецензию
// @Tags review
// @Description Получить рецензию по id
// @Accept json
// @Produce json
// @Param id path int true "ID рецензии"
// @Success 200 {object} dto.ReviewResponse
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /review/{id} [get]
func (h *ReviewEndpoints) GetReview(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, entity.NewClientError("невалидный id рецензии"))
	}
	review, err := h.reviewUC.GetReview(int(id))
	if err != nil {
		switch {
		case entity.Contains(err, entity.ErrNotFound):
			return utils.NewError(ctx, http.StatusNotFound, err)
		default:
			return utils.NewError(ctx, http.StatusInternalServerError, err)
		}
	}
	return utils.WriteJSON(ctx, review)
}

// GetMyContentReview
// @Summary Получить рецензию пользователя к контенту
// @Tags review
// @Description Получить рецензию пользователя к контенту
// @Produce json
// @Param content_id query int true "ID контента"
// @Success 200 {object} dto.ReviewResponse
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /review/myReview [get]
func (h *ReviewEndpoints) GetMyContentReview(ctx echo.Context) error {
	contentID, err := strconv.ParseInt(ctx.QueryParam("content_id"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, entity.NewClientError("невалидный id контента"))
	}
	userID, err := utils.GetUserIDFromSession(ctx, h.authUC)
	if err != nil {
		return err
	}
	reviews, err := h.reviewUC.GetContentReviewByAuthor(userID, int(contentID))
	if err != nil {
		switch {
		case entity.Contains(err, entity.ErrNotFound):
			return utils.NewError(ctx, http.StatusNotFound, err)
		default:
			return utils.NewError(ctx, http.StatusInternalServerError, err)
		}
	}
	return utils.WriteJSON(ctx, reviews)
}

// CreateReview
// @Summary Создать рецензию
// @Tags review
// @Description Создать рецензию
// @Accept json
// @Produce json
// @Param reviewCreate body dto.ReviewCreateRequest true "Данные для создания рецензии"
// @Success 200 {object} dto.ReviewResponse
// @Failure 400 {object} echo.HTTPError
// @Failure 401 {object} echo.HTTPError
// @Failure 409 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /review [post]
func (h *ReviewEndpoints) CreateReview(ctx echo.Context) error {
	var reviewCreate dto.ReviewCreateRequest
	err := ctx.Bind(&reviewCreate)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, entity.NewClientError("невалидный запрос"))
	}
	userID, err := utils.GetUserIDFromSession(ctx, h.authUC)
	if err != nil {
		return err
	}
	review, err := h.reviewUC.CreateReview(dto.ReviewCreate{
		ReviewCreateRequest: reviewCreate,
		UserID:              userID,
	})
	if err != nil {
		switch {
		case entity.Contains(err, entity.ErrBadRequest):
			return utils.NewError(ctx, http.StatusBadRequest, err)
		case entity.Contains(err, entity.ErrAlreadyExists):
			return utils.NewError(ctx, http.StatusConflict, err)
		default:
			return utils.NewError(ctx, http.StatusInternalServerError, err)
		}
	}
	return utils.WriteJSON(ctx, review)
}

// UpdateReview
// @Summary Обновить рецензию
// @Tags review
// @Description Обновить рецензию
// @Accept json
// @Produce json
// @Param reviewUpdate body dto.ReviewUpdateRequest true "Данные для обновления рецензии"
// @Success 200 {object} dto.ReviewResponse
// @Failure 400 {object} echo.HTTPError
// @Failure 401 {object} echo.HTTPError
// @Failure 403 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /review [put]
func (h *ReviewEndpoints) UpdateReview(ctx echo.Context) error {
	var reviewUpdate dto.ReviewUpdateRequest
	err := ctx.Bind(&reviewUpdate)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, entity.NewClientError("невалидный запрос"))
	}
	userID, err := utils.GetUserIDFromSession(ctx, h.authUC)
	if err != nil {
		return err
	}
	review, err := h.reviewUC.EditReview(dto.ReviewUpdate{
		ReviewUpdateRequest: reviewUpdate,
		UserID:              userID,
	})
	if err != nil {
		switch {
		case entity.Contains(err, entity.ErrNotFound):
			return utils.NewError(ctx, http.StatusNotFound, err)
		case entity.Contains(err, entity.ErrForbidden):
			return utils.NewError(ctx, http.StatusForbidden, err)
		case entity.Contains(err, entity.ErrBadRequest):
			return utils.NewError(ctx, http.StatusBadRequest, err)
		default:
			return utils.NewError(ctx, http.StatusInternalServerError, err)
		}
	}
	return utils.WriteJSON(ctx, review)
}

// DeleteReview
// @Summary Удалить рецензию
// @Tags review
// @Description Удалить рецензию
// @Accept json
// @Param id path int true "ID рецензии"
// @Success 200
// @Failure 400 {object} echo.HTTPError
// @Failure 401 {object} echo.HTTPError
// @Failure 403 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /review/{id} [delete]
func (h *ReviewEndpoints) DeleteReview(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, entity.NewClientError("невалидный id рецензии"))
	}
	userID, err := utils.GetUserIDFromSession(ctx, h.authUC)
	if err != nil {
		return err
	}
	err = h.reviewUC.DeleteReview(int(id), userID)
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
	return ctx.NoContent(http.StatusOK)
}

// GetRecentReviews
// @Summary Получить последние рецензии
// @Tags review
// @Description Получить последние рецензии
// @Produce json
// @Success 200 {object} dto.ReviewResponseList
// @Failure 500 {object} echo.HTTPError
// @Router /review/recent [get]
func (h *ReviewEndpoints) GetRecentReviews(ctx echo.Context) error {
	reviews, err := h.reviewUC.GetLatestReviews(3)
	if err != nil {
		return utils.NewError(ctx, http.StatusInternalServerError, err)
	}
	return utils.WriteJSON(ctx, reviews)
}

// GetUserLatestReviews
// @Summary Получить последние рецензии пользователя
// @Tags review
// @Description Получить последние рецензии пользователя
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} dto.ReviewResponseList
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /review/user/{id}/recent [get]
func (h *ReviewEndpoints) GetUserLatestReviews(ctx echo.Context) error {
	userID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, entity.NewClientError("невалидный id пользователя"))
	}
	reviews, err := h.reviewUC.GetUserReviews(int(userID), 3, 1)
	if err != nil {
		switch {
		case entity.Contains(err, entity.ErrNotFound):
			return utils.NewError(ctx, http.StatusNotFound, err)
		default:
			return utils.NewError(ctx, http.StatusInternalServerError, err)
		}
	}
	return utils.WriteJSON(ctx, reviews)
}

// GetUserReviews
// @Summary Получить рецензии пользователя
// @Tags review
// @Description Получить рецензии пользователя
// @Produce json
// @Param id path int true "ID пользователя"
// @Param page path int true "Номер страницы"
// @Success 200 {object} dto.UserReviewResponseList
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /review/user/{id}/{page} [get]
func (h *ReviewEndpoints) GetUserReviews(ctx echo.Context) error {
	// если не удалось получить id пользователя из сессии, то это не ошибка, просто неавторизованный пользователь
	// no-lint
	clientUserID, _ := utils.GetUserIDFromSession(ctx, h.authUC)
	userID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, entity.NewClientError("невалидный id пользователя"))
	}
	page, err := strconv.ParseInt(ctx.Param("page"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, entity.NewClientError("невалидный номер страницы"))
	}
	reviews, err := h.reviewUC.GetUserReviews(int(userID), 10, int(page))
	if err != nil {
		switch {
		case entity.Contains(err, entity.ErrNotFound):
			return utils.NewError(ctx, http.StatusNotFound, err)
		default:
			return utils.NewError(ctx, http.StatusInternalServerError, err)
		}
	}
	return utils.WriteJSON(ctx, dto.UserReviewResponseList{
		ReviewResponseList: *reviews,
		Me:                 clientUserID == int(userID),
	})
}

// GetContentReviews
// @Summary Получить рецензии контента
// @Tags review
// @Description Получить рецензии контента
// @Produce json
// @Param id path int true "ID контента"
// @Param page path int true "Номер страницы"
// @Success 200 {object} dto.ReviewResponseList
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /review/content/{id}/{page} [get]
func (h *ReviewEndpoints) GetContentReviews(ctx echo.Context) error {
	contentID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, entity.NewClientError("невалидный id контента"))
	}
	page, err := strconv.ParseInt(ctx.Param("page"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, entity.NewClientError("невалидный номер страницы"))
	}
	reviews, err := h.reviewUC.GetContentReviews(int(contentID), 10, int(page))
	if err != nil {
		switch {
		case entity.Contains(err, entity.ErrNotFound):
			return utils.NewError(ctx, http.StatusNotFound, err)
		default:
			return utils.NewError(ctx, http.StatusInternalServerError, err)
		}
	}
	return utils.WriteJSON(ctx, reviews)
}

// LikeReview
// @Summary Поставить лайк рецензии
// @Tags review
// @Description Поставить лайк рецензии
// @Accept json
// @Param id path int true "ID рецензии"
// @Success 200
// @Failure 400 {object} echo.HTTPError
// @Failure 401 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /review/{id}/like [put]
func (h *ReviewEndpoints) LikeReview(ctx echo.Context) error {
	reviewID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, entity.NewClientError("невалидный id рецензии"))
	}
	userID, err := utils.GetUserIDFromSession(ctx, h.authUC)
	if err != nil {
		return err
	}
	err = h.reviewUC.LikeReview(userID, int(reviewID))
	if err != nil {
		switch {
		case entity.Contains(err, entity.ErrNotFound):
			return utils.NewError(ctx, http.StatusNotFound, err)
		default:
			return utils.NewError(ctx, http.StatusInternalServerError, err)
		}
	}
	return ctx.NoContent(http.StatusOK)
}

// DislikeReview
// @Summary Поставить дизлайк рецензии
// @Tags review
// @Description Поставить дизлайк рецензии
// @Accept json
// @Param id path int true "ID рецензии"
// @Success 200
// @Failure 400 {object} echo.HTTPError
// @Failure 401 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /review/{id}/dislike [put]
func (h *ReviewEndpoints) DislikeReview(ctx echo.Context) error {
	reviewID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, entity.NewClientError("невалидный id рецензии"))
	}
	userID, err := utils.GetUserIDFromSession(ctx, h.authUC)
	if err != nil {
		return err
	}
	err = h.reviewUC.DislikeReview(userID, int(reviewID))
	if err != nil {
		switch {
		case entity.Contains(err, entity.ErrNotFound):
			return utils.NewError(ctx, http.StatusNotFound, err)
		default:
			return utils.NewError(ctx, http.StatusInternalServerError, err)
		}
	}
	return ctx.NoContent(http.StatusOK)
}

// UnlikeReview
// @Summary Убрать лайк с рецензии
// @Tags review
// @Description Убрать лайк с рецензии
// @Accept json
// @Param id path int true "ID рецензии"
// @Success 200
// @Failure 400 {object} echo.HTTPError
// @Failure 401 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /review/{id}/like [delete]
func (h *ReviewEndpoints) UnlikeReview(ctx echo.Context) error {
	reviewID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, entity.NewClientError("невалидный id рецензии"))
	}
	userID, err := utils.GetUserIDFromSession(ctx, h.authUC)
	if err != nil {
		return err
	}
	err = h.reviewUC.UnlikeReview(userID, int(reviewID))
	if err != nil {
		switch {
		case entity.Contains(err, entity.ErrNotFound):
			return utils.NewError(ctx, http.StatusNotFound, err)
		default:
			return utils.NewError(ctx, http.StatusInternalServerError, err)
		}
	}
	return ctx.NoContent(http.StatusOK)
}
