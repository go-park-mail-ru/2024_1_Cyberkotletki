package usecase

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_profanity.go
type Profanity interface {
	// FilterMessage фильтрует сообщение от ненормативной лексики
	FilterMessage(text string) (string, error)
}
