package dto

type Review struct {
	ID        int    `json:"id"        example:"1"          format:"int"`
	AuthorID  int    `json:"authorID"  example:"1"          format:"int"`
	ContentID int    `json:"contentID" example:"1"          format:"int"`
	Rating    int    `json:"rating"    example:"5"          format:"int"`
	Title     string `json:"title"     example:"Title"      format:"string"`
	Text      string `json:"text"      example:"i like it"  format:"string"`
	CreatedAt string `json:"createdAt" example:"1711821349" format:"int"`
	Likes     int    `json:"likes"     example:"5"          format:"int"`
	Dislikes  int    `json:"dislikes"  example:"5"          format:"int"`
}
