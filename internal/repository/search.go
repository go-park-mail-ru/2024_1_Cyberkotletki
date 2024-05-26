package repository

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_search.go
type Search interface {
	SearchContent(query string) ([]int, error)
	SearchPerson(query string) ([]int, error)
}
