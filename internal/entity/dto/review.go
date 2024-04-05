package dto

type Review struct {
	ID        int    `json:"id"        example:"1"                    format:"int"`
	AuthorID  int    `json:"authorID"  example:"1"                    format:"int"`
	ContentID int    `json:"contentID" example:"1"                    format:"int"`
	Rating    int    `json:"rating"    example:"5"                    format:"int"`
	Title     string `json:"title"     example:"Title"                format:"string"`
	Text      string `json:"text"      example:"i like it"            format:"string"`
	CreatedAt string `json:"createdAt" example:"2022-01-02T15:04:05Z" format:"int"`
	Likes     int    `json:"likes"     example:"5"                    format:"int"`
	Dislikes  int    `json:"dislikes"  example:"5"                    format:"int"`
}

// ReviewResponse - структура для ответа на запросы
type ReviewResponse struct {
	Review
	AuthorName   string `json:"authorName"   example:"Author"             format:"string"`
	AuthorAvatar string `json:"authorAvatar" example:"avatars/avatar.jpg" format:"string"`
	ContentName  string `json:"contentName"  example:"Content"            format:"string"`
}

type ReviewCreate struct {
	UserID    int    `json:"userID"    example:"1"         format:"int"`
	ContentID int    `json:"contentID" example:"1"         format:"int"`
	Rating    int    `json:"rating"    example:"5"         format:"int"`
	Title     string `json:"title"     example:"Title"     format:"string"`
	Text      string `json:"text"      example:"i like it" format:"string"`
}

type ReviewUpdate struct {
	UserID   int    `json:"userID"   example:"1"         format:"int"`
	ReviewID int    `json:"reviewID" example:"1"         format:"int"`
	Rating   int    `json:"rating"   example:"5"         format:"int"`
	Title    string `json:"title"    example:"Title"     format:"string"`
	Text     string `json:"text"     example:"i like it" format:"string"`
}
