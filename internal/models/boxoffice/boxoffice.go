package boxoffice

import "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/country"

// BoxOffice представляет кассовые сборы.
type BoxOffice struct {
	Id      int             `json:"id"`      // Уникальный идентификатор
	Country country.Country `json:"country"` // Страна, в которой были сборы
	Revenue int             `json:"revenue"` // Сумма сборов
}
