package dto

type CreateFavouriteRequest struct {
	ContentID int `json:"contentID" example:"1" format:"int" description:"идентификатор контента"`
	// nolint:lll
	Category string `json:"category" example:"favourite" format:"string" description:"может быть favourite/watching/watched/planned/rewatching/abandoned"`
}

type Favourite struct {
	Content PreviewContent `json:"content" description:"контент"`
	// nolint:lll
	Category string `json:"category" example:"favourite" format:"string" description:"может быть favourite/watching/watched/planned/rewatching/abandoned"`
}

type FavouritesResponse struct {
	Favourites []Favourite `json:"favourites" description:"список избранного"`
}

type FavouriteStatusResponse struct {
	// nolint:lll
	Status string `json:"status" example:"favourite" format:"string" description:"может быть favourite/watching/watched/planned/rewatching/abandoned"`
}
