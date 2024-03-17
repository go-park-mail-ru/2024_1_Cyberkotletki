package entity

/*
TODO: тесты
*/

type Film struct {
	Content
	Duration int `json:"duration"` // Продолжительность
}
