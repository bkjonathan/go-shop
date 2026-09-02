package interfaces

import "mime/multipart"

type UploadProvider interface {
	UploadFile(file *multipart.FileHeader, path string) (string, error)
	DeleteFile(filePath string) error
}
