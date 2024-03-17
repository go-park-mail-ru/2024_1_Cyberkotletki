package entity

// BoxOffice представляет кассовые сборы.
type BoxOffice struct {
	Id      int     `json:"id"`      // Уникальный идентификатор
	Country Country `json:"country"` // Страна, в которой были сборы
	Revenue int     `json:"revenue"` // Сумма сборов
}
