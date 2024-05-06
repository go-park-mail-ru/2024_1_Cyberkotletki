package http

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	"github.com/labstack/echo/v4"
	"net/http"
)

type SearchEndpoints struct {
	searchUC usecase.Search
}

func NewSearchEndpoints(searchUC usecase.Search) SearchEndpoints {
	return SearchEndpoints{searchUC: searchUC}
}

func (h *SearchEndpoints) Configure(server *echo.Group) {
	server.GET("", h.Search)
}

// Search
// @Tags Search
// @Description Поиск фильмов, сериалов и персон
// @Accept json
// @Param	query	query	string	true	"Поисковый запрос"
// @Success     200
// @Failure		400	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router /search [get]
func (h *SearchEndpoints) Search(ctx echo.Context) error {
	searchQuery := ctx.QueryParam("query")
	if searchQuery == "" {
		return utils.NewError(ctx, http.StatusBadRequest, "Пустой запрос", nil)
	}
	if len(searchQuery) > 100 {
		return utils.NewError(ctx, http.StatusBadRequest, "Слишком длинный запрос", nil)
	}
	searchResult, err := h.searchUC.Search(searchQuery)
	if err != nil {
		return utils.NewError(ctx, http.StatusInternalServerError, "Внутренняя ошибка сервера", err)
	}
	return ctx.JSON(http.StatusOK, searchResult)
}
