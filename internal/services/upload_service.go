package services

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/bkjonathan/go-shop/internal/interfaces"
)

type UploadService struct {
	provider interfaces.UploadProvider
}

func NewUploadService(provider interfaces.UploadProvider) *UploadService {
	return &UploadService{provider: provider}
}

func (s *UploadService) UploadFile(file *multipart.FileHeader, path string) (string, error) {
	return s.provider.UploadFile(file, path)
}

func (s *UploadService) DeleteFile(filePath string) error {
	return s.provider.DeleteFile(filePath)
}

func (s *UploadService) UploadProductImage(productID uint, file *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !isValidImageExtension(ext) {
		return "", fmt.Errorf("invalid file type: %s", ext)
	}

	path := fmt.Sprintf("products/%d/%s", productID, file.Filename)
	return s.UploadFile(file, path)
}

func isValidImageExtension(ext string) bool {
	validExtensions := []string{".jpg", ".jpeg", ".png", ".gif"}
	for _, validExt := range validExtensions {
		if ext == validExt {
			return true
		}
	}
	return false
}
