package imageupload

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const MaxFileSize = 5 * 1024 * 1024 // 5MB

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".svg":  true,
}

var allowedMimeTypes = map[string]bool{
	"image/jpeg":    true,
	"image/png":     true,
	"image/svg+xml": true,
}

func SaveImageFile(subfolder string, files []*multipart.FileHeader) (string, error) {
	if len(files) == 0 {
		return "", nil // Dosya zorunlu değilse boş dön
	}

	file := files[0]

	// Dosya boyutu kontrolü
	if file.Size > MaxFileSize {
		return "", fmt.Errorf("dosya boyutu 5MB'den büyük olamaz")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		return "", fmt.Errorf("uzantı '%s' desteklenmiyor", ext)
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// MIME kontrolü için ilk 512 byte'ı oku
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		return "", err
	}
	contentType := http.DetectContentType(buffer)
	if !allowedMimeTypes[contentType] {
		return "", fmt.Errorf("MIME türü '%s' desteklenmiyor", contentType)
	}
	_, err = src.Seek(0, 0) // Başa sar
	if err != nil {
		return "", err
	}

	// Dosya adını rastgele üret
	filename := uuid.New().String() + ext
	uploadPath := filepath.Join("uploads", subfolder)
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return "", err
	}

	dstPath := filepath.Join(uploadPath, filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Gerçek içeriği kopyala
	_, err = dst.ReadFrom(src)
	if err != nil {
		return "", err
	}

	return filename, nil
}
