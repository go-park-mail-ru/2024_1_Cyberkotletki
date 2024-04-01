package repository

type Static interface {
	GetStatic(staticID int) (string, error)
	UploadStatic(path, filename string, data []byte) (int, error)
}
