package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrInvalidUploadFile    = errors.New("invalid upload file")
	ErrUploadFileTooLarge   = errors.New("upload file too large")
	ErrUnsupportedImageType = errors.New("unsupported image type")
)

const maxMenuImageSize = 5 * 1024 * 1024

var allowedImageMIMETypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp":  ".webp",
}

func SaveMenuImage(fileHeader *multipart.FileHeader) (string, error) {
	if fileHeader == nil {
		return "", ErrInvalidUploadFile
	}

	if fileHeader.Size <= 0 {
		return "", ErrInvalidUploadFile
	}

	if fileHeader.Size > maxMenuImageSize {
		return "", ErrUploadFileTooLarge
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	if len(content) == 0 || int64(len(content)) > maxMenuImageSize {
		return "", ErrInvalidUploadFile
	}

	mimeType := http.DetectContentType(content)
	extension, ok := allowedImageMIMETypes[mimeType]
	if !ok {
		return "", ErrUnsupportedImageType
	}

	originalExt := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if originalExt != "" {
		switch originalExt {
		case ".jpg", ".jpeg":
			if extension != ".jpg" {
				return "", ErrUnsupportedImageType
			}
			extension = ".jpg"
		case ".png":
			if extension != ".png" {
				return "", ErrUnsupportedImageType
			}
		case ".webp":
			if extension != ".webp" {
				return "", ErrUnsupportedImageType
			}
		default:
			return "", ErrUnsupportedImageType
		}
	}

	if err := os.MkdirAll(filepath.Join("uploads", "menus"), 0o755); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%s%s", uuid.NewString(), extension)
	relativePath := filepath.ToSlash(filepath.Join("uploads", "menus", filename))
	if err := os.WriteFile(relativePath, content, 0o644); err != nil {
		return "", err
	}

	return "/" + relativePath, nil
}

func DeleteUploadedFile(relativePath string) error {
	trimmed := strings.TrimSpace(relativePath)
	if trimmed == "" {
		return nil
	}

	cleanPath := strings.TrimPrefix(trimmed, "/")
	return os.Remove(filepath.FromSlash(cleanPath))
}

func ResolvePublicURL(baseURL, relativePath string) string {
	trimmedPath := strings.TrimSpace(relativePath)
	if trimmedPath == "" {
		return ""
	}

	if strings.HasPrefix(trimmedPath, "http://") || strings.HasPrefix(trimmedPath, "https://") {
		return trimmedPath
	}

	trimmedBase := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if trimmedBase == "" {
		return trimmedPath
	}

	return trimmedBase + "/" + strings.TrimLeft(trimmedPath, "/")
}