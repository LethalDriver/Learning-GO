package service

type StorageService interface {
	DownloadFile(string) ([]byte, error)
	UploadFile([]byte) (string, error)
	DeleteFile(string) error
}
