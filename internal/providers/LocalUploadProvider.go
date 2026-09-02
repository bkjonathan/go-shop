package providers

import (
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// publicPrefix must match the route the files are served from in
// SetupRoutes: router.Static("/uploads", "./uploads").
const publicPrefix = "/uploads"

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
	return publicPrefix + "/" + strings.TrimPrefix(path, "/"), nil
}

func (p *LocalUploadProvider) DeleteFile(filePath string) error {
	// Accept either a stored public URL or a plain relative path.
	rel := strings.TrimPrefix(strings.TrimPrefix(filePath, publicPrefix), "/")
	fullPath := filepath.Join(p.basePath, rel)
	return os.Remove(fullPath)
}
