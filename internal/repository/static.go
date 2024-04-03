package repository

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_static.go
type Static interface {
	GetStatic(staticID int) (string, error)
	UploadStatic(path, filename string, data []byte) (int, error)
}
