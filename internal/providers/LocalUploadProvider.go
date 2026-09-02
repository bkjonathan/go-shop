package providers

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
)

type LocalUploadProvider struct {
	basePath string
}

func NewLocalUploadProvider(basePath string) *LocalUploadProvider {
	return &LocalUploadProvider{basePath: basePath}
}

func (p *LocalUploadProvider) UploadFile(file *multipart.FileHeader, path string) (string, error) {
	// Implementation for uploading a file to local storage
	fullPath := filepath.Join(p.basePath, path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return "", err
	}

	// Open Source file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Create destination file making sure the path is writeable.
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// read from source to destination
	if _, err := dst.ReadFrom(src); err != nil {
		return "", err
	}
	return fmt.Sprintf("/upload/%s", path), nil
}

func (p *LocalUploadProvider) DeleteFile(filePath string) error {
	fullPath := filepath.Join(p.basePath, filePath)
	return os.Remove(fullPath)
}
