package usecase

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_static.go
type Static interface {
	GetAvatar(staticID int) (string, error)
	UploadAvatar(data []byte) (int, error)
	GetStaticUrl(id int) (string, error)
}