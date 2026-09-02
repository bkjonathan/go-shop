package services

import (
	"crypto/rand"
	"encoding/hex"
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

	name, err := uniqueFileName(ext)
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("products/%d/%s", productID, name)
	return s.UploadFile(file, path)
}

// uniqueFileName keeps the caller's filename out of the storage key: it can
// collide with an existing image, or carry "../" and separators that escape the
// product folder on the local provider.
func uniqueFileName(ext string) (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate file name: %w", err)
	}
	return hex.EncodeToString(buf) + ext, nil
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
